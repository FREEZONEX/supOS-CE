package postgresql

import (
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeromicro/go-zero/core/logx"
)

type tableProcessInfo struct {
	data []map[string]string
	def  *types.CreateTopicDto
}

func persistence(dbPool *pgxpool.Pool, defaultSchema string, batchSize int, unsData []serviceApi.UnsData) error {
	if len(unsData) == 0 {
		return nil
	}
	// 准备表处理信息
	tableInfoMap := make(map[string]*tableProcessInfo, len(unsData))
	for _, data := range unsData {
		if len(data.Data) == 0 {
			continue
		}
		uns := data.Uns
		tagField := uns.GetTbFieldName()
		if tagField != "" {
			for _, da := range data.Data {
				da[tagField] = strconv.FormatInt(uns.Id, 10)
			}
		}
		tableName := uns.GetTable()
		tableInfo, ok := tableInfoMap[tableName]
		if !ok {
			tableInfo = &tableProcessInfo{
				data: data.Data,
				def:  uns,
			}
			tableInfoMap[tableName] = tableInfo
		} else {
			tableInfo.data = append(tableInfo.data, data.Data...)
		}
	}

	if len(tableInfoMap) == 0 {
		return nil
	}
	tbs := base.MapValues(tableInfoMap)
	var allErrors []string
	conn, er := dbPool.Acquire(context.Background())
	if er != nil {
		logPoolError("pg_persistence", time.Time{}, dbPool, "getConn", er)
		logx.Error("pg_persistence fail", er)
		return er
	}
	defer conn.Release()
	// 分批处理大数据量
	for _, segment := range base.Partition(tbs, batchSize) {
		var batch = &pgx.Batch{}
		for _, table := range segment {
			sql, params := getInsertStatement(table.def, table.data)
			batch.Queue(sql, params...)
		}
		// 执行批次
		err := execBatch(conn, batchTask{batch: batch, uns: segment}, defaultSchema, 0)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
	}
	if len(allErrors) > 0 {
		return fmt.Errorf("处理完成，但有错误: %s", strings.Join(allErrors, "; "))
	}

	return nil
}

type batchTask struct {
	batch *pgx.Batch
	uns   []*tableProcessInfo
}

func execBatch(conn *pgxpool.Conn, task batchTask, defaultSchema string, retry int) error {
	var retryTask = batchTask{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	br := conn.SendBatch(ctx, task.batch)

	defer func() {
		_ = br.Close()
		cancel()
	}()
	for i, seg := range task.uns {
		_, err := br.Exec()
		if err != nil {
			if retry > 0 {
				return fmt.Errorf("批次操作 %d 失败: %v, %s", i, err, seg.def.Alias)
			} else {
				if retryTask.batch == nil {
					retryTask.batch = &pgx.Batch{}
					retryTask.uns = make([]*tableProcessInfo, 0, len(task.uns))
				}
				q := task.batch.QueuedQueries[i]
				retryTask.batch.Queue(q.SQL, q.Arguments)
				retryTask.uns = append(retryTask.uns, seg)
			}
		}
	}
	if retryTask.batch != nil {
		uns := base.Map[*tableProcessInfo, *types.CreateTopicDto](retryTask.uns, func(e *tableProcessInfo) *types.CreateTopicDto {
			return e.def
		})
		tableInfoMap, er := ListTableInfos(conn, uns)
		if er != nil {
			return er
		}
		BatchCreateTables(conn, defaultSchema, uns, tableInfoMap)
		return execBatch(conn, retryTask, defaultSchema, retry+1)
	}
	return nil
}
