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

func persistence(dbPool *pgxpool.Pool, defaultSchema string, batchSize int, unsData []serviceApi.UnsData) error {
	if len(unsData) == 0 {
		return nil
	}
	// 准备表处理信息
	tableInfoMap := GetTableDataMap(unsData)
	if len(tableInfoMap) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	conn, err := dbPool.Acquire(ctx)
	cancel()
	if conn != nil {
		defer conn.Release()
	}
	if err != nil {
		return err
	} else if conn == nil {
		return fmt.Errorf("conn is nil")
	}
	tbs := base.MapValues(tableInfoMap)
	allErrors := SaveBatch(conn, defaultSchema, batchSize, tbs)
	if len(allErrors) > 0 {
		return fmt.Errorf("处理完成，但有错误: %s", strings.Join(allErrors, "; "))
	}
	return nil
}

func GetTableDataMap(unsData []serviceApi.UnsData) map[string]serviceApi.UnsData {
	tableInfoMap := make(map[string]serviceApi.UnsData, len(unsData))
	for _, data := range unsData {
		uns, list := data.Uns, data.Data
		if len(list) == 0 || uns == nil {
			continue
		}
		if tagField := uns.GetTbFieldName(); tagField != "" {
			for _, da := range list {
				da[tagField] = strconv.FormatInt(uns.Id, 10)
			}
		}
		tableName := uns.GetTable()
		tableInfo, ok := tableInfoMap[tableName]
		if !ok {
			tableInfo = serviceApi.UnsData{
				Data: list,
				Uns:  uns,
			}
			tableInfoMap[tableName] = tableInfo
		} else {
			tableInfo.Data = append(tableInfo.Data, list...)
		}
	}
	return tableInfoMap
}

func SaveBatch(conn *pgxpool.Conn, defaultSchema string, batchSize int, unsData []serviceApi.UnsData) (allErrors []string) {
	// 分批处理大数据量
	for _, segment := range base.Partition(unsData, batchSize) {
		var batch = &pgx.Batch{}
		for _, table := range segment {
			sql, params := getInsertStatement(table.Uns, table.Data)
			logx.Debugf("insert sql: %s, values: %+v", sql, params)
			batch.Queue(sql, params...)
		}
		// 执行批次
		err := execBatch(conn, batchTask{batch: batch, uns: segment}, defaultSchema, 0)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
	}
	return allErrors
}

type batchTask struct {
	batch *pgx.Batch
	uns   []serviceApi.UnsData
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
				return fmt.Errorf("批次操作 %d 失败: %v, %s", i, err, seg.Uns.Alias)
			} else {
				if retryTask.batch == nil {
					retryTask.batch = &pgx.Batch{}
					retryTask.uns = make([]serviceApi.UnsData, 0, len(task.uns))
				}
				q := task.batch.QueuedQueries[i]
				retryTask.batch.Queue(q.SQL, q.Arguments)
				retryTask.uns = append(retryTask.uns, seg)
			}
		}
	}
	if retryTask.batch != nil {
		uns := base.Map[serviceApi.UnsData, *types.CreateTopicDto](retryTask.uns, func(e serviceApi.UnsData) *types.CreateTopicDto {
			return e.Uns
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
