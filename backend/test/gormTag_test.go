package test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestGormColumn(t *testing.T) {
	example := ExampleStruct{}

	fmt.Println("=== 获取所有有效列名 ===")
	columns := GetValidColumns(example)
	for columnName, fieldName := range columns {
		fmt.Printf("字段名: %s -> 列名: %s\n", fieldName, columnName)

		//fmt.Println("\n=== 根据字段名获取列名 ===")
		//fmt.Printf("MountSource 列名: %s\n", GetColumnNameByField(example, "MountSource"))
		//fmt.Printf("PathName 列名: %s\n", GetColumnNameByField(example, "PathName"))
		//fmt.Printf("Labels 列名: %s\n", GetColumnNameByField(example, "Labels"))
		//fmt.Printf("Name 列名: %s\n", GetColumnNameByField(example, "Name"))
		//fmt.Printf("Age 列名: %s\n", GetColumnNameByField(example, "Age"))
	}
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
