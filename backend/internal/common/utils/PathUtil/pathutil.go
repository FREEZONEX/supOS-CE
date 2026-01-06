package PathUtil

import (
	"backend/internal/common/constants"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mozillazg/go-pinyin"
)

var (
	topicPattern     = regexp.MustCompile(constants.TopicReg)
	aliasPattern     = regexp.MustCompile(constants.AliasReg)
	fieldNamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
)

// ValidTopicFormat 校验topic格式
func ValidTopicFormat(topic string, dataType *int16) bool {
	// 报警规则类型没有格式限制
	if dataType != nil && *dataType == constants.AlarmRuleType {
		return true
	}

	// 不能以 / 开头或包含 //
	if strings.HasPrefix(topic, "/") || strings.Contains(topic, "//") {
		return false
	}

	// 如果没有 /，检查整个字符串
	if !strings.Contains(topic, "/") {
		return topicPattern.MatchString(topic)
	}

	// 检查每一部分
	parts := strings.Split(topic, "/")
	for _, part := range parts {
		if !topicPattern.MatchString(part) {
			return false
		}
	}

	return true
}

// IsAliasFormatOK 校验alias格式
func IsAliasFormatOK(alias string) bool {
	return aliasPattern.MatchString(alias)
}

// IsFieldNameFormatOK 校验字段名格式
func IsFieldNameFormatOK(name string) bool {
	return fieldNamePattern.MatchString(name)
}

// GetName 从路径中提取名称
func GetName(path string) string {
	if path == "" {
		return path
	}

	ed := len(path) - 1
	if path[ed] == '/' {
		ed--
		if ed < 0 {
			return ""
		}
	}

	x := strings.LastIndex(path[:ed+1], "/")
	return path[x+1 : ed+1]
}

// CleanPath 移除路径末尾的斜杠
func CleanPath(path string) string {
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		return path[:len(path)-1]
	}
	return path
}

// IsRootPath 检查是否为根路径
func IsRootPath(path string) bool {
	st, ed := 0, len(path)

	if st < ed && path[0] == '/' {
		st++
	}
	if st < ed && path[len(path)-1] == '/' {
		ed--
	}

	for i := st; i < ed; i++ {
		if path[i] == '/' {
			return false
		}
	}
	return true
}

// generatePinyinAlias 是一个辅助函数，用于生成拼音别名
func generatePinyinAlias(s string) string {
	a := pinyin.NewArgs()
	a.Style = pinyin.Normal
	a.Heteronym = false
	pinyinResult := pinyin.Pinyin(s, a)
	var builder strings.Builder
	for _, item := range pinyinResult {
		builder.WriteString(item[0])
	}
	return builder.String()
}

// GenerateFileAlias 为文件生成别名
func GenerateFileAlias(path string) string {
	aliasPath := path
	pathLen := len(path)

	if pathLen > 20 {
		startPo := lastNthIndex(path, "/", 2)
		if startPo >= 0 {
			aliasPath = path[startPo+1:]
		} else if pathLen > 20 {
			aliasPath = path[pathLen-20:]
		}
	}

	aliasPath = strings.ReplaceAll(aliasPath, "/", "_")
	aliasPath = strings.ReplaceAll(aliasPath, "-", "_")
	piny := generatePinyinAlias(aliasPath)
	if len(piny) > 0 {
		aliasPath = piny
	}
	if len(aliasPath) > 0 && !isLetter(aliasPath[0]) {
		aliasPath = "_" + aliasPath
	}

	if pathLen < 20 {
		return aliasPath
	}

	hash := md5.Sum([]byte(path))
	hashStr := hex.EncodeToString(hash[:])

	prefixLen := 4
	if len(aliasPath) < 4 {
		prefixLen = len(aliasPath)
	}

	return aliasPath[:prefixLen] + "_" + hashStr
}
func GenerateAliasWithRandom(name string) string {
	return randomByName(name)
}

// GenerateAlias 根据路径和类型生成别名
func GenerateAlias(path string, pathType int16) string {
	if pathType == constants.PathTypeFile { // 2 is file type
		return GenerateFileAlias(path)
	}

	var aliasPath string
	if pathType == constants.PathTypeDir { // 0 is dir type
		// folder:folder1/、folder1/folder2/
		if strings.Count(path, "/") > 1 {
			// folder:folder1/folder2/
			aliasPath = path[lastNthIndex(path, "/", 2):]
		} else {
			// folder:folder1/
			aliasPath = path
		}
	} else {
		aliasPath = path
	}

	return randomByName(aliasPath)
}

func randomByName(aliasPath string) string {
	aliasPath = strings.ReplaceAll(aliasPath, "/", "_")
	aliasPath = strings.ReplaceAll(aliasPath, "-", "_")
	aliasPath = generatePinyinAlias(aliasPath)

	if len(aliasPath) > 20 {
		aliasPath = aliasPath[:20]
	}

	uuidStr := uuid.NewString()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")[:20]
	aliasPath = aliasPath + "_" + uuidStr

	if len(aliasPath) > 0 && !isLetter(aliasPath[0]) {
		aliasPath = "_" + aliasPath
	}
	return aliasPath
}

// SubParentPath 提取父路径
func SubParentPath(path string) string {
	if path == "" {
		return ""
	}

	ed := len(path) - 1
	if path[ed] == '/' {
		ed--
		if ed < 0 {
			return ""
		}
	}

	x := strings.LastIndex(path[:ed+1], "/")
	if x > 0 {
		return path[:x]
	}
	return ""
}

// GenIDForPath 为路径生成MD5 ID
func GenIDForPath(path string) string {
	hash := md5.Sum([]byte(path))
	return hex.EncodeToString(hash[:])
}

// GetNextNodeAfterBasePath 获取基础路径后的下一个节点
func GetNextNodeAfterBasePath(basePath, path string) string {
	if basePath == "" {
		parts := strings.Split(path, "/")
		if len(parts) > 0 {
			return parts[0]
		}
		return ""
	}

	if strings.HasPrefix(path, basePath) {
		remaining := strings.TrimPrefix(path, basePath)
		remaining = strings.TrimPrefix(remaining, "/")

		idx := strings.Index(remaining, "/")
		if idx != -1 {
			return remaining[:idx]
		}
		if remaining != "" {
			return remaining
		}
	}
	return ""
}

// IsNextLevelPath 检查是否是直接子路径
func IsNextLevelPath(basePath, path string) bool {
	if basePath == "" {
		return !strings.Contains(path, "/")
	}
	baseDepth := len(strings.Split(basePath, "/"))
	return strings.HasPrefix(path, basePath+"/") && len(strings.Split(path, "/")) == baseDepth+1
}

// EscapeName 使名称成为有效的标识符
func EscapeName(name string) string {
	var builder strings.Builder
	changed := false
	for _, c := range name {
		if !isJavaIdentifierPart(c) || c == '$' {
			builder.WriteRune('_')
			changed = true
		} else {
			builder.WriteRune(c)
		}
	}
	if changed {
		return builder.String()
	}
	return name
}

// RemoveLastLevel 从路径中移除最后一个级别
func RemoveLastLevel(path string) string {
	if path == "" {
		return ""
	}
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return ""
	}
	return path[:lastSlash]
}

// GetLastLevel 返回路径的最后一个级别
func GetLastLevel(path string) string {
	if path == "" {
		return ""
	}
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return path
	}
	return path[lastSlash+1:]
}

// GenerateUniqueName 生成唯一的名称
func GenerateUniqueName(baseName string, all []string) string {
	maxSuffix := -1
	exists := false
	for _, name := range all {
		if name == baseName {
			exists = true
			break
		}
	}
	if !exists {
		return baseName
	}

	for _, name := range all {
		if name == baseName {
			maxSuffix = max(maxSuffix, 0)
		} else if strings.HasPrefix(name, baseName+"-") {
			suffixStr := strings.TrimPrefix(name, baseName+"-")
			if suffix, err := strconv.Atoi(suffixStr); err == nil {
				maxSuffix = max(maxSuffix, suffix)
			}
		}
	}

	return fmt.Sprintf("%s-%d", baseName, maxSuffix+1)
}

// --- 辅助函数 ---

// lastNthIndex 查找第n个子字符串的最后一个索引
func lastNthIndex(s, substr string, n int) int {
	if n <= 0 {
		return -1
	}
	idx := len(s)
	for i := 0; i < n; i++ {
		idx = strings.LastIndex(s[:idx], substr)
		if idx == -1 {
			return -1
		}
	}
	return idx
}

// isLetter 检查字节是否为字母
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// isJavaIdentifierPart 检查符文是否是有效的Java标识符部分
func isJavaIdentifierPart(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// max 返回两个整数中的较大者
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
