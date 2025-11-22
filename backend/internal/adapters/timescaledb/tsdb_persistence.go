package timescaledb

import (
	"backend/internal/adapters/postgresql"
	"backend/internal/common"
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
	tempTableName string
	data          []map[string]string
	def           *types.CreateTopicDto
}

func (t *tableProcessInfo) GetTableName() string {
	return t.def.GetTable()
}

func persistence(dbPool *pgxpool.Pool, defaultSchema string, batchSize int, unsData []serviceApi.UnsData) error {
	if len(unsData) == 0 {
		return nil
	}

	// 获取单个连接
	ctx := context.Background()
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}
	defer conn.Release()

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
				tempTableName: fmt.Sprintf("tm_%s_%d", strings.ToLower(tableName), common.NextId()),
				data:          data.Data,
				def:           uns,
			}
			tableInfoMap[tableName] = tableInfo
		} else {
			tableInfo.data = append(tableInfo.data, data.Data...)
		}
	}

	if len(tableInfoMap) == 0 {
		return nil
	}
	tableInfos := base.MapValues(tableInfoMap)
	var allErrors []string

	// 步骤1: 在一个SendBatch中创建所有临时表
	if err := createAllTempTables(dbPool, conn, defaultSchema, tableInfos, 0); err != nil {
		allErrors = append(allErrors, fmt.Sprintf("创建临时表失败: %v", err))
		return fmt.Errorf("处理失败: %s", strings.Join(allErrors, "; "))
	}

	// 步骤2: 循环每个表，使用CopyFrom导入数据到对应的临时表
	for _, tableInfo := range tableInfos {
		if err := copyDataToTempTable(conn, batchSize, tableInfo); err != nil {
			allErrors = append(allErrors, fmt.Sprintf("表 %s 数据导入失败: %v", tableInfo.GetTableName(), err))
			// 继续处理其他表
		}
	}

	// 如果有表的数据导入失败，我们可能不想继续合并操作
	if len(allErrors) > 0 {
		return fmt.Errorf("部分表数据导入失败: %s", strings.Join(allErrors, "; "))
	}

	// 步骤3: 在一个SendBatch中执行所有表的合并操作
	if err := mergeAllTables(conn, tableInfos); err != nil {
		allErrors = append(allErrors, fmt.Sprintf("合并数据失败: %v", err))
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("处理完成，但有错误: %s", strings.Join(allErrors, "; "))
	}

	return nil
}

func createAllTempTables(dbPool *pgxpool.Pool, conn *pgxpool.Conn, defaultSchema string, tableInfos []*tableProcessInfo, retry int) error {
	batch := &pgx.Batch{}

	// 为每个表添加创建临时表的操作
	for _, tableInfo := range tableInfos {
		// 创建临时表
		createSQL := fmt.Sprintf(`CREATE TEMP TABLE %s (LIKE "%s" EXCLUDING INDEXES)`, tableInfo.tempTableName, tableInfo.GetTableName())
		logx.Debug("创建临时表:", retry, createSQL)
		batch.Queue(createSQL)
	}

	// 执行批次
	br := conn.SendBatch(context.Background(), batch)
	defer br.Close()

	// 检查所有操作结果
	var retryTables []*tableProcessInfo
	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			if retry > 0 {
				return fmt.Errorf("【%d】批次操作 %d 失败: %v", retry, i, err)
			} else {
				retryTables = append(retryTables, tableInfos[i])
			}
		}
	}
	if len(retryTables) > 0 {
		br.Close()
		uns := base.Map[*tableProcessInfo, *types.CreateTopicDto](retryTables, func(e *tableProcessInfo) *types.CreateTopicDto {
			return e.def
		})
		tableInfoMap, _ := postgresql.ListTableInfos(conn, uns)
		postgresql.BatchCreateTables(dbPool, defaultSchema, uns, tableInfoMap)
		return createAllTempTables(dbPool, conn, defaultSchema, retryTables, retry+1)
	}
	return br.Close()
}

func copyDataToTempTable(conn *pgxpool.Conn, batchSize int, tableInfo *tableProcessInfo) error {
	if len(tableInfo.data) == 0 {
		return nil
	}

	// 构建列名（排除自动生成的字段）
	var columns = base.Map(tableInfo.def.Fields, func(e *types.FieldDefine) string {
		return e.Name
	})

	// 分批处理大数据量
	for i := 0; i < len(tableInfo.data); i += batchSize {
		end := i + batchSize
		if end > len(tableInfo.data) {
			end = len(tableInfo.data)
		}

		batch := tableInfo.data[i:end]
		batch = postgresql.DeduplicationById(tableInfo.def, batch)
		// 准备数据行
		rows := make([][]interface{}, len(batch))
		for j, record := range batch {
			row := make([]interface{}, len(columns))
			for k, f := range tableInfo.def.Fields {
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
		logx.Debugf("%s: rows: %+v", tableInfo.def.Alias, rows)

		// 执行COPY
		_, err := conn.CopyFrom(
			context.Background(),
			pgx.Identifier{tableInfo.tempTableName},
			columns,
			pgx.CopyFromRows(rows),
		)

		if err != nil {
			return fmt.Errorf("COPY数据到临时表 %s 失败: %v", tableInfo.tempTableName, err)
		}

		//p.log.Debugf("表 %s 批次 %d-%d 数据导入成功，数据量: %d",
		//	tableInfo.GetTableName(), i, end, len(batch))
	}

	return nil
}

func mergeAllTables(conn *pgxpool.Conn, tableInfos []*tableProcessInfo) error {
	batch := &pgx.Batch{}

	// 为每个表添加合并操作
	for _, tableInfo := range tableInfos {
		primaryFields := tableInfo.def.GetPrimaryField()
		// 合并数据SQL
		mergeSQL := &base.StringBuilder{}
		mergeSQL.Grow(128 + len(primaryFields)*10)
		mergeSQL.Append(`INSERT INTO "`).Append(tableInfo.GetTableName()).
			Append(`" SELECT *  FROM `).Append(tableInfo.tempTableName)

		if len(primaryFields) > 0 {
			mergeSQL.Append(` ON CONFLICT (`)
			for i, f := range primaryFields {
				if i > 0 {
					mergeSQL.Append(`, `)
				}
				mergeSQL.Append(`"`).Append(f).Append(`"`)
			}
			mergeSQL.Append(`)`)
			if len(tableInfo.def.Fields) > len(primaryFields) {
				mergeSQL.Append(" DO UPDATE SET ")
				postgresql.GetUpdateColumns(tableInfo.def, mergeSQL)
			} else {
				mergeSQL.Append(" DO NOTHING ")
			}
		}
		batch.Queue(mergeSQL.String())
	}

	// 执行批次
	br := conn.SendBatch(context.Background(), batch)

	// 检查所有操作结果
	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("合并操作 %d (表 %s) 失败: %v, SQL: %v",
				i, tableInfos[i].GetTableName(), err, batch.QueuedQueries[i].SQL)
		}
	}

	return br.Close()
}
