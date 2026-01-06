package postgresql

import (
	"fmt"
	"strings"
)

// 获取完整表名
func getFullTableName(tableName string) string {
	// 如果已经是引号包围的，直接返回
	if len(tableName) >= 2 && tableName[0] == '"' && tableName[len(tableName)-1] == '"' {
		return tableName
	}

	// 查找点号分隔符
	dot := strings.Index(tableName, ".")
	if dot > 0 {
		// 分割数据库名和表名
		db := tableName[:dot]
		table := tableName[dot+1:]

		// 返回带引号的完整表名
		return fmt.Sprintf(`"%s"."%s"`, db, table)
	}

	// 返回带引号的表名
	return fmt.Sprintf(`"%s"`, tableName)
}

// 获取不带双引号的表名
func getCleanTableName(tableName string) string {
	st := 0
	// 查找最后一个点之后的位置
	if idx := strings.LastIndex(tableName, "."); idx >= 0 {
		st = idx + 1
	}

	ed := len(tableName)

	// 移除开头的引号
	if st < len(tableName) && tableName[st] == '"' {
		st++
	}

	// 移除结尾的引号
	if ed > 0 && tableName[ed-1] == '"' {
		ed--
	}
	// 移除斜杠部分（如果有）
	cleanName := tableName[st:ed]
	if idx := strings.LastIndex(cleanName, "/"); idx >= 0 {
		cleanName = cleanName[idx+1:]
	}

	return cleanName
}
