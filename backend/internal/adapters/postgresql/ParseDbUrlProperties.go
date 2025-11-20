package postgresql

import (
	"backend/internal/common/serviceApi"
	"net/url"
	"strings"
)

func ParseDbUrlProperties(dbUrl string) serviceApi.DataSourceProperties {
	props := serviceApi.DataSourceProperties{
		Url: dbUrl,
	}

	// 解析URL
	parsedUrl, err := url.Parse(dbUrl)
	if err != nil {
		return props
	}

	// 提取用户名和密码
	if parsedUrl.User != nil {
		props.UserName = parsedUrl.User.Username()
		password, hasPassword := parsedUrl.User.Password()
		if hasPassword {
			props.Password = password
		}
	}

	// 提取数据库名（Path中的最后一个部分）
	if parsedUrl.Path != "" {
		pathParts := strings.Split(strings.Trim(parsedUrl.Path, "/"), "/")
		if len(pathParts) > 0 {
			// 使用路径中的数据库名作为备选Schema
			props.Schema = pathParts[len(pathParts)-1]
		}
	}

	// 优先从查询参数中获取search_path作为Schema
	if parsedUrl.Query().Get("search_path") != "" {
		props.Schema = parsedUrl.Query().Get("search_path")
	}

	return props
}
