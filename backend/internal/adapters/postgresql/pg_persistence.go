package postgresql

import (
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
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
	// 准备表处理信息
	tableInfoMap := GetTableDataMap(unsData)
	if len(tableInfoMap) == 0 {
		return nil
	}
	tbs := base.MapValues(tableInfoMap)
	allErrors := SaveBatch(dbPool, defaultSchema, batchSize, tbs)()
	if len(allErrors) > 0 {
		return fmt.Errorf("处理完成，但有错误: %s", strings.Join(allErrors, "; "))
	}
	return nil
}

func GetTableDataMap(unsData []serviceApi.UnsData) map[string]*serviceApi.UnsData {
	tableInfoMap := make(map[string]*serviceApi.UnsData, len(unsData))
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
			tableInfo = &serviceApi.UnsData{
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

func SaveBatch(dbPool *pgxpool.Pool, defaultSchema string, batchSize int, unsData []*serviceApi.UnsData) func() (allErrors []string) {
	// 分批处理大数据量
	var wg sync.WaitGroup
	var parts = base.Partition(unsData, batchSize)
	var lock sync.Mutex
	errorMsgs := make([]string, 0, len(parts))
	for _, segment := range parts {
		wg.Add(1)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			conn, err := dbPool.Acquire(ctx)
			defer func() {
				if err != nil {
					lock.Lock()
					errorMsgs = append(errorMsgs, err.Error())
					lock.Unlock()
				}
				if conn != nil {
					conn.Release()
				}
				wg.Done()
			}()
			cancel()

			if err != nil {
				return
			} else if conn == nil {
				err = fmt.Errorf("conn is nil")
				return
			}
			var batch = &pgx.Batch{}
			for _, table := range segment {
				sql, params := getInsertStatement(table.Uns, table.Data)
				logx.Debugf("insert sql: %s, values: %+v", sql, params)
				batch.Queue(sql, params...)
			}
			// 执行批次
			err = execBatch(conn, batchTask{batch: batch, uns: segment}, defaultSchema, 0)
		}()
	}
	return func() []string {
		wg.Wait()
		return errorMsgs
	}
}

type batchTask struct {
	batch *pgx.Batch
	uns   []*serviceApi.UnsData
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
			if strings.Contains(err.Error(), "mismatched") {
				sql := task.batch.QueuedQueries[i].SQL
				args := task.batch.QueuedQueries[i].Arguments
				argsBs, _ := json.Marshal(args)
				countArgsInSql := strings.Count(sql, "$")
				logx.Errorf("参数不匹配[%s]: sql[%d]=%s, Args[%d]= %s", seg.Uns.Alias, countArgsInSql, sql, len(args), string(argsBs))
				return fmt.Errorf("参数不匹配[%s]: %v", seg.Uns.Alias, err)
			}
			var pgEr *pgconn.PgError
			if retry > 0 {
				return fmt.Errorf("批次操作失败~[%s]: %v ", seg.Uns.Alias, err)
			} else if errors.As(err, &pgEr) && (pgEr.Code == "42P01" || pgEr.Code == "42703") {
				//错误码 	名称 	               描述 	                            应对建议
				//42P01 	undefined_table 	查询的表不存在 	                检查表名拼写、schema 是否正确，确认表已创建
				//42703 	undefined_column 	查询中引用了不存在的字段 	        核对字段名，使用 DESCRIBE table 或 \d table 检查结构
				//23503 	foreign_key_violation 	外键引用的记录不存在 	        确保被引用的主表中存在对应记录
				//23505 	unique_violation 	违反唯一约束（如重复主键或唯一索引） 	检查插入/更新数据是否已存在，使用 ON CONFLICT 处理
				//22001 	string_data_right_truncation 	字符串长度超出字段定义 	检查字段长度（如 VARCHAR(50)），截断或扩大字段
				//42601 	syntax_error 	SQL 语法错误 	                    检查关键字、括号、引号、分号等语法结构
				//55000 	object_in_use 	对象正在被使用（如删除正在使用的表） 	    等待事务结束，或使用 DROP ... CASCADE 强制清理依赖
				if retryTask.batch == nil {
					retryTask.batch = &pgx.Batch{}
					retryTask.uns = make([]*serviceApi.UnsData, 0, len(task.uns))
				}
				q := task.batch.QueuedQueries[i]
				retryTask.batch.Queue(q.SQL, q.Arguments)
				retryTask.uns = append(retryTask.uns, seg)
			} else if pgEr != nil {
				if pgEr.Code == "42601" || strings.HasPrefix(pgEr.Code, "2") {
					sql := task.batch.QueuedQueries[i].SQL
					args := task.batch.QueuedQueries[i].Arguments
					argsBs, _ := json.Marshal(args)
					countArgsInSql := strings.Count(sql, "$")
					if pgEr.Code == "42601" {
						logx.Errorf("语法错误[%s]: sql[%d]=%s, Args[%d]= %s", seg.Uns.Alias, countArgsInSql, sql, len(args), string(argsBs))
						return fmt.Errorf("语法错误[%s]: %v ", seg.Uns.Alias, err)
					} else {
						logx.Errorf("约束错误[%s]: sql[%d]=%s, Args[%d]= %s", seg.Uns.Alias, countArgsInSql, sql, len(args), string(argsBs))
						return fmt.Errorf("约束错误[%s]: %v ", seg.Uns.Alias, err)
					}
				} else {
					return fmt.Errorf("pg操作失败[%s]: %v ", seg.Uns.Alias, err)
				}
			} else {
				return fmt.Errorf("批次操作失败![%s]: %v ", seg.Uns.Alias, err)
			}
		}
	}
	if retryTask.batch != nil {
		uns := base.Map[*serviceApi.UnsData, *types.CreateTopicDto](retryTask.uns, func(e *serviceApi.UnsData) *types.CreateTopicDto {
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
