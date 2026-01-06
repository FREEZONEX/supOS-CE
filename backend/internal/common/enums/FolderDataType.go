package enums

import (
	"backend/internal/common/constants"
	"strings"
)

// FolderDataType 文件夹数据类型
type FolderDataType int16

const (
	NORMAL FolderDataType = iota
	STATE
	ACTION
	METRIC //时序 计算 聚合 引用
)

// FolderDataType 映射关系
var (
	folderDataTypeNames = []string{
		"NORMAL", "uns.folder.type.normal", // Name,ShowName(i18n)
		"STATE", "uns.folder.type.state",
		"ACTION", "uns.folder.type.action",
		"METRIC", "uns.folder.type.metrics",
	}
)

// TypeIndex 返回枚举的类型索引
func (fdt FolderDataType) TypeIndex() int16 {
	return int16(fdt)
}
func (fdt FolderDataType) Name() string {
	i := int(fdt)
	if i < 0 || i >= len(folderDataTypeNames)/2 {
		i = 0
	}
	return folderDataTypeNames[i<<1]
}

// String 返回枚举的i18n名称
func (fdt FolderDataType) String() string {
	i := int(fdt)
	if i < 0 || i >= len(folderDataTypeNames)/2 {
		i = 0
	}
	return folderDataTypeNames[i*2+1]
}
func GetFolderDataTypeByName(name string) (FolderDataType, bool) {
	name = strings.ToUpper(name)
	for i := 0; i < len(folderDataTypeNames); i += 2 {
		if folderDataTypeNames[i] == name {
			return FolderDataType(i / 2), true
		}
	}
	return NORMAL, false
}

// GetFolderDataType 根据索引获取FolderDataType
func GetFolderDataType(index int16) FolderDataType {
	i := int(index)
	if i < 0 || i >= len(folderDataTypeNames)/2 {
		i = 0
	}
	return FolderDataType(i)
}

// IsTypeMatched 检查父数据类型和文件数据类型是否匹配
func IsTypeMatched(parentDataType, fileDataType *int16) bool {
	if parentDataType == nil || fileDataType == nil {
		return false
	}
	dt := *fileDataType
	switch *parentDataType {
	case STATE.TypeIndex():
		return dt == constants.RelationType ||
			dt == constants.JsonbType ||
			dt == constants.MergeType ||
			dt == constants.CitingType
	case ACTION.TypeIndex():
		return dt == constants.JsonbType
	case METRIC.TypeIndex():
		return dt == constants.TimeSequenceType ||
			dt == constants.CalculationRealType ||
			dt == constants.CalculationHistType ||
			dt == constants.MergeType ||
			dt == constants.CitingType
	default:
		return false
	}
}

// GetDefaultParentType 根据文件数据类型获取默认的父类型
func GetDefaultParentType(fileDataType *int16) int16 {
	if fileDataType == nil {
		return STATE.TypeIndex()
	}

	switch *fileDataType {
	case constants.RelationType:
		return STATE.TypeIndex()
	case constants.JsonbType:
		return ACTION.TypeIndex()
	case constants.TimeSequenceType, constants.CalculationRealType, constants.CalculationHistType:
		return METRIC.TypeIndex()
	default:
		return STATE.TypeIndex()
	}
}
