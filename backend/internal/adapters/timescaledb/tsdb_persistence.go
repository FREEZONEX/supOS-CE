package timescaledb

import (
	"backend/internal/adapters/postgresql"
	"backend/internal/common/constants"
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
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
		uns := column2uns(rs.conflict.columns)
		err = copyAndMergeFromTempTable(conn, uns, rs.conflict)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
	}
	if len(rs.normal.rows) > 0 {
		//直接COPY数据到目标表
		uns := column2uns(rs.normal.columns)
		err = copyDataToTable(context.Background(), conn, uns.TableName, rs.normal)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
	}
	if len(allErrors) > 0 {
		return fmt.Errorf("处理完成，但有错误: %s", strings.Join(allErrors, "; "))
	}
	return nil
}
func column2uns(cols []string) *types.CreateTopicDto {
	var uns = &types.CreateTopicDto{Fields: make([]*types.FieldDefine, 0, 32), TableName: "uns_timeserial"}
	for _, col := range cols {
		uns.Fields = append(uns.Fields, &types.FieldDefine{
			Name:   col,
			Unique: base.V2p(col == constants.SysFieldCreateTime || col == constants.SystemSeqTag),
		})
	}
	return uns
}
func copyAndMergeFromTempTable(conn *pgxpool.Conn, uns *types.CreateTopicDto, params copyParams) error {
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
		return err
	}

	// 3. 从临时表合并到主表
	err = mergeFromTempTable(ctx, tx, uns, tempTableName)
	if err != nil {
		return err
	}

	// 4. 提交事务，提交时临时表自动销毁
	return tx.Commit(ctx)
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
