package timescaledb

import (
	"backend/internal/adapters/postgresql"
	"backend/internal/common/constants"
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeromicro/go-zero/core/logx"
)

func persistence(dbPool *pgxpool.Pool, defaultSchema string, batchSize int, unsData []serviceApi.UnsData) error {
	if len(unsData) == 0 {
		return nil
	}

	rs := preprocess(unsData)

	// 获取单个连接
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	conn, err := dbPool.Acquire(ctx)
	cancel()
	if conn != nil {
		defer conn.Release()
	}
	if err != nil {
		logPoolError("persistence", time.Time{}, dbPool, "getConn", err)
		return fmt.Errorf("获取数据库连接失败: %v", err)
	} else if conn == nil {
		return fmt.Errorf("conn is nil")
	}
	var allErrors []string
	if len(rs.conflict.rows) > 0 {
		uns := column2uns(rs.conflict.columns, unsData[0].Uns)
		err = copyAndMergeFromTempTable(conn, uns, defaultSchema, rs.conflict)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
	}
	if len(rs.normal.rows) > 0 {
		//直接COPY数据到目标表
		uns := column2uns(rs.normal.columns, unsData[0].Uns)
		err = copyDataToTable(context.Background(), conn, uns.TableName, rs.normal)
		if err != nil {
			if retry, reCreateTable := shouldRetry(err); retry {
				if reCreateTable {
					er := fixTable(conn, uns)
					if er != nil {
						logx.Error("更新物理表失败!", er)
						return err
					}
				}
				err = copyDataToTable(context.Background(), conn, uns.TableName, rs.normal) //重新拷贝一次
			}
			if err != nil {
				allErrors = append(allErrors, err.Error())
			}
		}
	}
	if len(allErrors) > 0 {
		return fmt.Errorf("处理完成，但有错误: %s", strings.Join(allErrors, "; "))
	}
	return nil
}
func column2uns(cols []string, reference *types.UnsDefinition) *types.UnsDefinition {
	var uns = types.CreateTopicDto{Fields: make([]*types.FieldDefine, 0, 32), TableName: "uns_timeserial",
		DataSrcID: types.SrcJdbcTypeTimeScaleDB.Id()}
	for _, col := range cols {
		fd := &types.FieldDefine{
			Name:   col,
			Unique: base.V2p(col == constants.SysFieldCreateTime || col == constants.SysFieldID),
		}
		if def, has := reference.GetFieldDefines().FieldsMap[col]; has {
			fd.Type = def.Type
		} else if pi := strings.LastIndex(col, "_"); pi > 0 {
			prev := col[:pi]
			fd.Type = prefixToFieldType[prev]
		}

		uns.Fields = append(uns.Fields, fd)
	}
	return &types.UnsDefinition{CreateTopicDto: uns}
}
func copyAndMergeFromTempTable(conn *pgxpool.Conn, uns *types.UnsDefinition, defaultSchema string, params copyParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. 创建临时表 (此临时表仅在此事务内存在)
	tableName := uns.GetTable()
	tempTableName := fmt.Sprintf("tmp_%s_%d", strings.ToLower(tableName), time.Now().UnixNano())
	createTmpTableSQL := fmt.Sprintf(`CREATE TEMP TABLE %s (LIKE "%s" EXCLUDING INDEXES)  ON COMMIT DROP`, tempTableName, tableName)
	_, err = tx.Exec(ctx, createTmpTableSQL)
	if err != nil {
		return err
	}

	// 2. COPY数据到临时表
	err = copyDataToTable(ctx, tx, tempTableName, params)
	if err != nil {
		_ = tx.Rollback(ctx)
		if retry, reCreateTable := shouldRetry(err); retry {
			if reCreateTable {
				er := fixTable(conn, uns)
				if er != nil {
					logx.Error("创建物理表失败", er)
					return err
				}
			}
			tx, err = conn.Begin(ctx)
			if err != nil {
				return err
			}
			_, err = tx.Exec(ctx, createTmpTableSQL) //重建临时表
			if err != nil {
				return err
			}
			err = copyDataToTable(ctx, tx, tempTableName, params) //重新拷贝一次
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// 3. 从临时表合并到主表
	err = mergeFromTempTable(ctx, tx, uns, tempTableName)
	if err != nil {
		return err
	}

	// 4. 提交事务，提交时临时表自动销毁
	return tx.Commit(ctx)
}
func fixTable(conn *pgxpool.Conn, uns *types.UnsDefinition) error {
	info, er := postgresql.ListTableInfos(conn, []string{uns.TableName})
	if er != nil {
		return er
	}
	// 构建当前字段类型映射
	curFieldTypes := make(map[string]*types.FieldDefine)
	for _, field := range uns.Fields {
		curFieldTypes[field.Name] = field
	}
	renameColSQL := ""
	tableInfo := info[uns.TableName]
	if tableInfo != nil {
		for field, Type := range tableInfo.FieldTypes {
			_, exists := curFieldTypes[field]
			if !exists {
				if Type == types.FieldTypeLong && field != constants.SysFieldID && !strings.Contains(field, "_") {
					renameColSQL = "ALTER TABLE " + uns.TableName + ` RENAME COLUMN "` + field + `" TO "` + constants.QosField + `"`

					delete(curFieldTypes, constants.QosField)
				}
				continue
			}
			// 从映射中移除已存在的字段
			delete(curFieldTypes, field)
		}
	}
	if len(renameColSQL) == 0 && len(curFieldTypes) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if len(renameColSQL) > 0 {
		_, err = tx.Exec(ctx, renameColSQL)
		logx.Debug("renameColSQL: ", er, renameColSQL)
	}
	if len(curFieldTypes) > 0 {
		alterSQL := base.StringBuilder{}
		alterSQL.Grow(256)
		alterSQL.Append("ALTER TABLE ").Append(uns.TableName)
		for _, def := range curFieldTypes {
			field := def.GetName()
			var typeStr string
			if def.Type == types.FieldTypeString {
				typeStr = "text"
			} else {
				typeStr = postgresql.GetTypeDefineWithSerial(def, false)
			}
			alterSQL.Append(fmt.Sprintf(` ADD IF NOT EXISTS "%s" %s,`, field, typeStr))
		}
		sql := alterSQL.SetLast(' ').String()
		_, err = tx.Exec(ctx, sql)
		logx.Debug("addColSQL: ", err, sql)
	}
	if err != nil {
		logx.Error("重试失败:", err.Error())
		return err
	} else {
		return tx.Commit(ctx)
	}
}
func shouldRetry(err error) (retry, reCreateTable bool) {
	var pgEr *pgconn.PgError
	if errors.As(err, &pgEr) {
		switch pgEr.Code {
		case "42P01", "42703":
			retry, reCreateTable = true, true
		default:
			if strings.HasPrefix(pgEr.Code, "08") { //Class 08 — Connection Exception
				retry = true
			}
		}
	}
	return
}

type copyFromer interface {
	CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error)
}

func copyDataToTable(ctx context.Context, conn copyFromer, tableName string, params copyParams) error {
	// 执行COPY
	count, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{tableName},
		params.columns,
		pgx.CopyFromRows(params.rows),
	)
	logx.Debugf("copyRows-> %s [%d]: %d, err: %v, cols: %v", tableName, len(params.rows), count, err, params.columns)
	return err
}

func mergeFromTempTable(ctx context.Context, conn pgx.Tx, uns *types.UnsDefinition, tempTableName string) error {
	primaryFields := uns.GetPrimaryField()
	// 合并数据SQL
	mergeSQL := &base.StringBuilder{}
	mergeSQL.Grow(128 + len(primaryFields)*10)
	mergeSQL.Append(`INSERT INTO "`).Append(uns.GetTable()).
		Append(`"AS t SELECT *  FROM `).Append(tempTableName)

	if len(primaryFields) > 0 {
		mergeSQL.Append(` ON CONFLICT (`)
		for i, f := range primaryFields {
			if i > 0 {
				mergeSQL.Append(`, `)
			}
			mergeSQL.Append(`"`).Append(f).Append(`"`)
		}
		mergeSQL.Append(`)`)
		if len(uns.Fields) > len(primaryFields) {
			mergeSQL.Append(" DO UPDATE SET ")
			postgresql.GetUpdateColumns(uns, mergeSQL)
		} else {
			mergeSQL.Append(" DO NOTHING ")
		}
	}
	_, er := conn.Exec(ctx, mergeSQL.String())
	return er
}
func logPoolError(name string, start time.Time, pool *pgxpool.Pool, sql string, err error) {
	if !start.IsZero() {
		duration := time.Since(start)
		stats := pool.Stat()
		logx.Errorf("[%s FAILED] sql:%s, err:%v, duration:%v, poolStats:(Total:%d, Idle:%d, Acquired:%d)", name,
			sql, err, duration, stats.TotalConns(), stats.IdleConns(), stats.AcquiredConns())
	} else {
		stats := pool.Stat()
		logx.Errorf("[%s FAILED] sql:%s, err:%v, poolStats:(Total:%d, Idle:%d, Acquired:%d)", name,
			sql, err, stats.TotalConns(), stats.IdleConns(), stats.AcquiredConns())
	}
}
