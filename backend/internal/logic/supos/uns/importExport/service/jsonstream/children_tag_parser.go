package jsonstream

import (
	"reflect"
	"strings"
)

func getChildrenJsonTagName(ts any) string {
	t := reflect.TypeOf(ts)
	// 如果传入的是指针，获取其指向的类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 确保是结构体类型
	if t.Kind() != reflect.Struct {
		return ""
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		if field.Type.Kind() == reflect.Slice &&
			(field.Type.Elem() == t || (field.Type.Elem().Kind() == reflect.Ptr && field.Type.Elem().Elem() == t)) {
			name := parseColumnName(jsonTag)
			if len(name) > 0 {
				return name
			}
		}
	}
	return ""
}
func parseColumnName(tag string) string {
	// 如果标签包含"-"，表示忽略该字段
	if strings.Contains(tag, "-") {
		return ""
	}
	tagParts := strings.Split(tag, ",")
	if len(tagParts) > 0 {
		return strings.TrimSpace(tagParts[0])
	}
	return ""
}
