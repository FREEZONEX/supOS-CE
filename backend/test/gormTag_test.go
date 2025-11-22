package test

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"strings"
	"testing"
)

func TestGormColumn(t *testing.T) {
	jsonStr := `{
  "name" : "state-wait",
  "parentId" : "",
  "alias" : "stateWait2",
  "dataType" : 2,
  "save2db" : false,
  "pathType" : 2,
  "labelNames" : [ ],
  "fields" : [ {
    "name" : "value",
    "type" : "FLOAT"
  } ],
  "parentDataType" : 1
}`
	t.Log(jsonStr)
}

// MD5Digest16 生成16字符的MD5摘要（等效于Java的digestHex16）
func MD5Digest16(input string) string {
	// 创建MD5哈希对象
	hasher := md5.New()

	// 写入输入字符串
	hasher.Write([]byte(input))

	// 计算MD5哈希值
	hashBytes := hasher.Sum(nil)

	// 转换为32字符的十六进制字符串
	fullHash := hex.EncodeToString(hashBytes)

	// 取前16字符作为结果
	return fullHash[8:24]
}

// 定义示例结构体
type ExampleStruct struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"column:user_name" json:"name"`
	MountSource string `gorm:"column:mount_source" json:"mount_source"`
	PathName    string `gorm:"-" json:"pathName"`
	Labels      string `gorm:"->;<-:false;column:labels" json:"labels"`
	Age         int    `json:"age"`
	CreatedAt   int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   int64  `gorm:"autoUpdateTime" json:"updated_at"`
	Ignored     string `gorm:"-" json:"ignored"`
}

// GetValidColumns 获取结构体有效的数据库列名
func GetValidColumns(model interface{}) map[string]string {
	result := make(map[string]string)
	t := reflect.TypeOf(model)

	// 如果传入的是指针，获取其指向的类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 确保是结构体类型
	if t.Kind() != reflect.Struct {
		return result
	}

	// 遍历结构体字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 获取gorm标签
		gormTag := field.Tag.Get("gorm")
		if gormTag == "" {
			continue
		}

		// 检查是否被忽略的字段
		if strings.Contains(gormTag, "-") {
			continue
		}

		// 解析column名称
		columnName := parseColumnName(gormTag)
		if columnName != "" {
			result[columnName] = field.Name
		}
	}

	return result
}

// parseColumnName 解析gorm标签中的column名称
func parseColumnName(gormTag string) string {
	// 如果标签包含"-"，表示忽略该字段
	if strings.Contains(gormTag, "-") {
		return ""
	}

	// 查找column:xxx的模式
	tagParts := strings.Split(gormTag, ";")
	for _, part := range tagParts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}

	// 如果没有明确指定column，使用默认的蛇形命名
	return "" //toSnakeCase(fieldName)
}

// toSnakeCase 将驼峰命名转换为蛇形命名
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, char := range s {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteByte(byte(char))
	}
	return strings.ToLower(result.String())
}

// GetColumnNameByField 根据字段名获取对应的数据库列名
func GetColumnNameByField(model interface{}, fieldName string) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, found := t.FieldByName(fieldName)
	if !found {
		return ""
	}

	gormTag := field.Tag.Get("gorm")
	if gormTag == "" || strings.Contains(gormTag, "-") {
		return ""
	}

	return parseColumnName(gormTag)
}
