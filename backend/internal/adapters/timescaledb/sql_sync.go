package timescaledb

import (
	"backend/internal/types"
	"fmt"
)

// GenerateSyncSQLs 生成所有同步 SQL
func (g *SQLGenerator) GenerateSyncSQLs(
	physicsTableFields []*types.FieldDefine,
	unsList []UnsViewInfo,
) *SyncSQLs {

	result := NewSyncSQLs()

	// 1. 分析需要的字段
	required := g.AnalyzeRequiredFields(unsList, physicsTableFields)

	// 2. 根据物理表是否存在生成不同的 SQL
	if len(physicsTableFields) == 0 {
		// 物理表不存在，创建表
		createSQLs := g.GenerateCreateTableSQL(required)
		result.CreateTableSQL = append(result.CreateTableSQL, createSQLs...)
	} else if len(required.FieldsToAdd) > 0 {
		// 物理表存在，但有新字段需要添加
		alterSQLs := g.GenerateAlterTableSQL(required)
		result.AlterTableSQL = append(result.AlterTableSQL, alterSQLs...)
	}

	// 3. 处理每个 Uns
	for _, unsView := range unsList {
		unsInfo := unsView.Uns
		viewInfo := unsView.View
		// 获取字段映射
		mappings, removedFields, _ := g.getFieldMappings(unsInfo, viewInfo.Columns)

		// 保存字段映射信息
		var mappingInfo []string
		for viewField, sourceCol := range mappings {
			mappingInfo = append(mappingInfo, fmt.Sprintf("%s -> %s", viewField, sourceCol))
		}
		result.FieldMappings[unsInfo.GetAlias()] = mappingInfo

		// 生成更新数据的 SQL
		updateSQL := g.GenerateDataUpdateSQL(unsInfo, removedFields)
		if updateSQL != "" {
			result.UpdateDataSQL = append(result.UpdateDataSQL, updateSQL)
		}

		// 生成创建视图的 SQL
		viewSQL := g.GenerateViewSQL(unsInfo, mappings)
		result.CreateViewSQL = append(result.CreateViewSQL, viewSQL)
	}

	return result
}

// GetSQLExecutionOrder 获取 SQL 执行顺序
func (s *SyncSQLs) GetSQLExecutionOrder() []string {
	// 定义执行顺序：
	// 1. 先创建表
	// 2. 再修改表（添加字段）
	// 3. 然后更新数据（将删除的字段设为 NULL）
	// 4. 最后创建视图

	var orderedSQLs []string

	// 1. 创建表
	orderedSQLs = append(orderedSQLs, s.CreateTableSQL...)

	// 2. 修改表
	orderedSQLs = append(orderedSQLs, s.AlterTableSQL...)

	// 3. 更新数据
	orderedSQLs = append(orderedSQLs, s.UpdateDataSQL...)

	// 4. 创建视图
	orderedSQLs = append(orderedSQLs, s.CreateViewSQL...)

	return orderedSQLs
}

// HasErrors 检查是否有错误
func (s *SyncSQLs) HasErrors() bool {
	return len(s.Errors) > 0
}

// AddError 添加错误
func (s *SyncSQLs) AddError(err error) {
	s.Errors = append(s.Errors, err)
}
