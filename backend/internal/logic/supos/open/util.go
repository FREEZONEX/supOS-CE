// Package open provides common utility functions for the open module
package open

import "backend/internal/types"

// ConvertInstanceDetailToDefinition converts InstanceDetail with Fields to Definition
func ConvertInstanceDetailToDefinition(detail *types.InstanceDetail) interface{} {
	// 创建一个 map 来存储转换后的结果
	result := make(map[string]interface{})

	// 转换基础字段
	result["id"] = detail.Id
	result["alias"] = detail.ParentAlias
	result["parentAlias"] = detail.ParentAlias
	result["path"] = detail.Path
	result["dataType"] = detail.DataType
	result["parentDataType"] = detail.ParentDataType
	result["pathType"] = detail.PathType
	result["createTime"] = detail.CreateTime
	result["updateTime"] = detail.UpdateTime
	result["protocol"] = detail.Protocol
	result["modelDescription"] = detail.ModelDescription
	result["description"] = detail.Description
	result["persistence"] = detail.WithSave2db // withSave2db 转换为 persistence
	result["refers"] = detail.Refers
	result["name"] = detail.Name
	result["displayName"] = detail.DisplayName
	result["pathName"] = detail.PathName
	result["extendProperties"] = detail.Extend // extend 转换为 extendProperties
	result["payload"] = detail.Payload
	result["templateName"] = detail.TemplateName
	result["templateAlias"] = detail.TemplateAlias

	if detail.JsonFields != nil {
		result["definition"] = detail.JsonFields
	}

	return result
}

// ConvertModelDetailToDefinition converts ModelDetail with Fields to Definition
func ConvertModelDetailToDefinition(detail *types.ModelDetail) interface{} {
	// 创建一个 map 来存储转换后的结果
	result := make(map[string]interface{})

	// 转换基础字段
	result["id"] = detail.Id
	result["topic"] = detail.Topic
	result["alias"] = detail.Alias
	result["parentAlias"] = detail.ParentAlias
	result["path"] = detail.Path
	result["pathType"] = detail.PathType
	result["dataType"] = detail.DataType
	result["createTime"] = detail.CreateTime
	result["updateTime"] = detail.UpdateTime
	result["description"] = detail.Description
	result["name"] = detail.Name
	result["displayName"] = detail.DisplayName
	result["pathName"] = detail.PathName
	result["modelId"] = detail.ModelId
	result["modelName"] = detail.ModelName
	result["extend"] = detail.Extend
	result["templateAlias"] = detail.TemplateAlias
	result["mount"] = detail.Mount

	// 将 Fields 转换为 Definition
	if detail.Fields != nil {
		result["definition"] = detail.Fields
	}

	return result
}
