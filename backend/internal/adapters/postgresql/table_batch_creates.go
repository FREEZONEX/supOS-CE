package postgresql

import (
	"backend/internal/common/constants"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeromicro/go-zero/core/logx"
)

type sendBatcher interface {
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}

func BatchCreateTables(conn sendBatcher, defaultSchema string, topics []types.UnsInfo, tableInfoMap map[string]*TableInfo) []error {
	if len(topics) == 0 {
		return nil
	}
	if tableInfoMap == nil {
		tableInfoMap = map[string]*TableInfo{}
	}
	alterTableSQLs := make([]string, 0, len(topics))
	hyperTableSQLs := make([]string, 0, len(topics))

	tables := make(map[string]bool, len(topics))
	for _, uns := range topics {
		tableName := uns.GetTable()
		dbName := defaultSchema
		if dot := strings.Index(tableName, "."); dot > 0 {
			dbName = tableName[:dot]
			tableName = tableName[dot+1:]
		}
		if _, has := tables[tableName]; has {
			continue
		}
		tables[tableName] = true
		tableInfo := tableInfoMap[tableName]
		ch := checkTableModify(uns, dbName, tableName, &alterTableSQLs, tableInfo)
		if ch == MDF_NEW_TABLE || ch == MDF_TYPE_CHANGED {
			createTableSQL := getCreateTableSQLByUns(uns)
			alterTableSQLs = append(alterTableSQLs, createTableSQL)
			if uns.GetSrcJdbcType().TypeCode() == constants.TimeSequenceType {
				ct := uns.GetTimestampField()
				timeScaleDbCreateTableSQL := "SELECT create_hypertable('" + dbName + ".\"" + tableName + "\"', '" + ct + "',chunk_time_interval => INTERVAL '1 day')"
				hyperTableSQLs = append(hyperTableSQLs, timeScaleDbCreateTableSQL)
			}
		}
	}
	errs := make([]error, 0, 8)
	for i, sqlList := range base.Partition(alterTableSQLs, constants.SQLBatchSize) {
		batch := &pgx.Batch{}
		for j, sql := range sqlList {
			logx.Debugf("Batch SQL[%d-%d]: %s", i, j, sql)
			batch.Queue(sql)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		// 执行批次
		br := conn.SendBatch(ctx, batch)
		for i := 0; i < batch.Len(); i++ {
			_, err := br.Exec()
			if err != nil {
				errs = append(errs, err)
			}
		}
		_ = br.Close()
		cancel()
	}
	for i, sqlList := range base.Partition(hyperTableSQLs, constants.SQLBatchSize) {
		batch := &pgx.Batch{}
		for j, sql := range sqlList {
			logx.Debugf("Batch hyperTable SQL[%d-%d]: %s", i, j, sql)
			batch.Queue(sql)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// 执行批次
		br := conn.SendBatch(ctx, batch)
		for i := 0; i < batch.Len(); i++ {
			_, err := br.Exec()
			if err != nil {
				errs = append(errs, err)
			}
		}
		_ = br.Close()
		cancel()
	}
	if len(errs) > 0 {
		if pool, isPool := conn.(*pgxpool.Pool); isPool {
			logPoolError("queryPrimary rows", time.Time{}, pool, alterTableSQLs[0], errs[0])
		}
	}
	return errs
}
