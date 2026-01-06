package timescaledb

import (
	"backend/internal/adapters/postgresql"
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeromicro/go-zero/core/logx"
)

var _smallBatchSize = intEnv("OS_SMALL_SIZE", 500)

func intEnv(key string, def int) int {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
			return parsed
		}
	}
	return def
}

func persistence(dbPool *pgxpool.Pool, defaultSchema string, batchSize int, unsData []serviceApi.UnsData) error {
	if len(unsData) == 0 {
		return nil
	}
	// 准备表处理信息
	tableInfoMap := postgresql.GetTableDataMap(unsData)
	if len(tableInfoMap) == 0 {
		return nil
	}
	smallData := make([]*serviceApi.UnsData, 0, batchSize)
	bigData := make([]*serviceApi.UnsData, 0, batchSize)
	for _, tableInfo := range tableInfoMap {
		if len(tableInfo.Data) < _smallBatchSize {
			smallData = append(smallData, tableInfo)
		} else {
			bigData = append(bigData, tableInfo)
		}
	}

	logx.Debugf("tsdb persistence: %d , %d\n", len(smallData), len(bigData))
	var allErrors []string
	if len(smallData) > 0 {
		postgresql.SaveBatch(dbPool, defaultSchema, 1000, smallData)
	}

	if len(bigData) > 0 {
		// 获取单个连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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
		for _, tableInfo := range bigData {
			er := processSingleTableInTx(conn, tableInfo, batchSize)
			if er != nil {
				allErrors = append(allErrors, er.Error())
			}
		}
	}
	if len(allErrors) > 0 {
		return fmt.Errorf("处理完成，但有错误: %s", strings.Join(allErrors, "; "))
	}
	return nil
}
func processSingleTableInTx(conn *pgxpool.Conn, tableInfo *serviceApi.UnsData, batchSize int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. 创建临时表 (此临时表仅在此事务内存在)
	tableName := tableInfo.Uns.GetTable()
	tempTableName := fmt.Sprintf("tmp_%s_%d", strings.ToLower(tableName), time.Now().UnixNano())
	createTmpTableSQL := fmt.Sprintf(`CREATE TEMP TABLE %s (LIKE "%s" EXCLUDING INDEXES)  ON COMMIT DROP`, tempTableName, tableName)
	_, err = tx.Exec(ctx, createTmpTableSQL)
	if err != nil {
		return err
	}

	// 2. COPY数据到临时表
	err = copyDataToTempTable(ctx, tx, batchSize, tableInfo, tempTableName)
	if err != nil {
		return err
	}

	// 3. 从临时表合并到主表
	err = mergeFromTempTable(ctx, tx, tableInfo.Uns, tempTableName)
	if err != nil {
		return err
	}

	// 4. 提交事务，提交时临时表自动销毁
	return tx.Commit(ctx)
}

func copyDataToTempTable(ctx context.Context, conn pgx.Tx, batchSize int, tableInfo *serviceApi.UnsData, tempTableName string) error {
	if len(tableInfo.Data) == 0 {
		return nil
	}
	uns, data := tableInfo.Uns, tableInfo.Data
	// 构建列名（排除自动生成的字段）
	var columns = base.Map(uns.Fields, func(e *types.FieldDefine) string {
		return e.Name
	})

	// 分批处理大数据量
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		batch := data[i:end]
		batch = postgresql.DeduplicationById(uns, batch)
		// 准备数据行
		rows := make([][]interface{}, len(batch))
		for j, record := range batch {
			row := make([]interface{}, len(columns))
			for k, f := range uns.Fields {
				v, has := record[f.Name]
				if !has {
					row[k] = f.GetType().DefaultValue()
					continue
				}
				if f.Type == types.FieldTypeDatetime {
					mill, _ := strconv.ParseFloat(v, 64)
					if mill > 0 {
						utcTime := time.UnixMilli(int64(mill)).UTC()
						v = utcTime.Format("2006-01-02 15:04:05.000") + "+00"
					}
				}
				row[k] = v
			}
			rows[j] = row
		}
		//logx.Debugf("%s: rows: %+v", uns.Alias, rows)

		// 执行COPY
		_, err := conn.CopyFrom(
			ctx,
			pgx.Identifier{tempTableName},
			columns,
			pgx.CopyFromRows(rows),
		)

		if err != nil {
			return fmt.Errorf("COPY数据到临时表 %s 失败: %v", tempTableName, err)
		}

		//p.log.Debugf("表 %s 批次 %d-%d 数据导入成功，数据量: %d",
		//	tableInfo.GetTableName(), i, end, len(batch))
	}

	return nil
}

func mergeFromTempTable(ctx context.Context, conn pgx.Tx, uns *types.CreateTopicDto, tempTableName string) error {
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
