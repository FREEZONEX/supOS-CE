package grafanautil

import (
	"backend/internal/common/serviceApi"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"bytes"
	"context"
	"crypto/md5"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"backend/internal/common/constants"
	grafanadto "backend/internal/common/dto/grafana"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/google/uuid"
)

//go:embed templates/*
var templates embed.FS
var _grafanaUrl string

// GrafanaUtils provides utility functions for Grafana operations.
//
// GetGrafanaURL returns the Grafana URL based on the runtime environment.
func GetGrafanaURL() string {
	if _grafanaUrl == "" {
		_grafanaUrl = spring.GetBean[*svc.ServiceContext]().Config.GrafanaUrl
	}
	return _grafanaUrl
}

// GetDashboardUUIDByAlias generates a dashboard UUID from an alias using MD5.
func GetDashboardUUIDByAlias(alias string) string {
	// 创建MD5哈希对象
	hasher := md5.New()
	// 写入输入字符串
	hasher.Write([]byte(alias))
	// 计算MD5哈希值
	hashBytes := hasher.Sum(nil)
	// 转换为32字符的十六进制字符串
	fullHash := hex.EncodeToString(hashBytes)
	// 取前16字符作为结果
	return fullHash[8:24]
}

// DeleteDashboard deletes a Grafana dashboard by UID.
func DeleteDashboard(ctx context.Context, uid string) error {
	url := GetGrafanaURL() + "/api/dashboards/uid/" + uid
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.WithContext(ctx).Debugf("Delete Dashboard: %s, response: %s", uid, string(body))
	return nil
}

// DeleteDatasource deletes a Grafana datasource by UID.
func DeleteDatasource(ctx context.Context, uid string) error {
	url := GetGrafanaURL() + "/api/datasources/uid/" + uid
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.WithContext(ctx).Debugf("Delete DataSource: %s, response: %s", uid, string(body))
	return nil
}

// GetDataSourceByName retrieves a Grafana datasource by name.
func GetDataSourceByName(ctx context.Context, name string) int {
	logger := logx.WithContext(ctx)
	url := GetGrafanaURL() + "/api/datasources/name/" + name
	logger.Debugf("查询 datasource 请求: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode
	}
	body, _ := io.ReadAll(resp.Body)
	logger.Debugf("查询 datasource 返回结果: %s", string(body))

	return resp.StatusCode
}

// CreateDatasource creates a Grafana datasource.
func CreateDatasource(ctx context.Context, jdbcType types.SrcJdbcType, dsProps serviceApi.DataSourceProperties, reCreate bool) (bool, error) {
	title := jdbcType.Alias()
	datasource := &grafanadto.GrafanaDataSourceDto{
		User:     dsProps.UserName,
		Password: dsProps.Password,
		UID:      GetDashboardUUIDByAlias(title),
		Name:     title,
	}

	if reCreate {
		// Delete first, then create
		_ = DeleteDatasource(ctx, datasource.UID)
	}
	logger := logx.WithContext(ctx)
	var dsTemplate string
	switch jdbcType.Id() {
	case types.SrcJdbcTypePostgresql.Id():
		dsTemplate = LoadTemplate(ctx, "templates/pg-datasource.json")
		datasource.URL = constants.PGJDBCURL
	case types.SrcJdbcTypeTimeScaleDB.Id():
		dsTemplate = LoadTemplate(ctx, "templates/pg-datasource.json")
		datasource.URL = constants.TSDBJDBCURL
	case types.SrcJdbcTypeTdEngine.Id():
		datasource.URL = constants.TDJDBCURL
		datasource.CreateBasicAuth()
		dsTemplate = LoadTemplate(ctx, "templates/td-datasource.json")
	default:
		return false, fmt.Errorf("unsupported JDBC type: %d", jdbcType.Id())
	}
	datasource.URL = dsProps.HostPort
	dsJSON := formatTemplate(dsTemplate, datasource)
	logger.Debug("创建 datasource 请求: ", dsJSON)

	host := GetGrafanaURL()
	{
		url := host + "/api/datasources/name/" + datasource.Name
		resp, err := http.Get(url)
		if err != nil {
			logger.Debugf("查询 datasource 失败: %v %s", err, url)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				logger.Debugf("查询 datasource 返回结果: %s", string(body))
				if len(body) > 0 && body[0] == '{' {
					var rsMap map[string]interface{}
					err = json.Unmarshal(body, &rsMap)
					if err == nil && len(rsMap) > 0 {
						var oldUrl = rsMap["url"].(string)
						var oldUser = rsMap["user"].(string)
						if oldUrl != datasource.URL || oldUser != datasource.User {
							var uid = rsMap["uid"].(string)
							logger.Infof("准备删除重建数据源[%s],因为: url: %s->%s, user: %s->%s", uid, oldUrl, datasource.URL, oldUser, datasource.User)
							_ = DeleteDatasource(ctx, uid)
						}
					}
				}
			}
		}
	}

	resp, err := http.Post(host+"/api/datasources", "application/json", bytes.NewBufferString(dsJSON))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logger.Infof("创建 datasource 返回结果: %s", string(body))

	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusConflict, nil
}

// CreateDatasourceByBody creates a Grafana datasource from a JSON body.
func CreateDatasourceByBody(ctx context.Context, name, body string, reCreate bool) (string, error) {
	logger := logx.WithContext(ctx)
	var bodyJSON map[string]any
	if err := json.Unmarshal([]byte(body), &bodyJSON); err != nil {
		return "", err
	}
	bodyJSON["name"] = name

	if reCreate {
		// Delete first, then create
		url := GetGrafanaURL() + "/api/datasources/name/" + name
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		client := &http.Client{}
		_, _ = client.Do(req)
	}

	newBody, _ := json.Marshal(bodyJSON)
	logger.Info("创建 datasource 请求: ", string(newBody))

	resp, err := http.Post(GetGrafanaURL()+"/api/datasources", "application/json", bytes.NewBuffer(newBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	logger.Info("创建 datasource 返回结果: ", string(respBody))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		return "", nil
	}
	return string(respBody), nil
}

// CreateDashboard creates a Grafana dashboard.
func CreateDashboard(uid string, ctx context.Context, table, tagNameCondition string, jdbcType types.SrcJdbcType, schema, title, columns, ct string) error {
	var template string
	var dbParams map[string]any

	switch jdbcType {
	case types.SrcJdbcTypePostgresql:
		template = LoadTemplate(ctx, "templates/pg-dashboard.json")
		dbParams = map[string]any{
			"title":          title,
			"uid":            uid,
			"dataSourceType": jdbcType.DataSrcType,
			"dataSourceUid":  GetDatasourceUUIDByJDBC(jdbcType),
			"schema":         schema,
			"tableName":      table,
			"columns":        columns,
		}
	case types.SrcJdbcTypeTdEngine:
		template = LoadTemplate(ctx, "templates/td-dashboard.json")
		dbParams = map[string]any{
			"title":            title,
			"uid":              uid,
			"dataSourceType":   jdbcType.DataSrcType,
			"dataSourceUid":    GetDatasourceUUIDByJDBC(jdbcType),
			"schema":           schema,
			"tableName":        table,
			"tagNameCondition": tagNameCondition,
			"columns":          columns,
		}
	case types.SrcJdbcTypeTimeScaleDB:
		template = LoadTemplate(ctx, "templates/ts-dashboard.json")
		dbParams = map[string]any{
			"title":            title,
			"uid":              uid,
			"dataSourceType":   jdbcType.DataSrcType,
			"dataSourceUid":    GetDatasourceUUIDByJDBC(jdbcType),
			"schema":           schema,
			"tableName":        table,
			"tagNameCondition": tagNameCondition,
			"columns":          columns,
		}
	default:
		return fmt.Errorf("unsupported JDBC type: %d", jdbcType.Id())
	}

	dbParams["sys_field_create_time"] = ct
	dashboardJSON := FormatTemplateMap(template, dbParams)
	logger := logx.WithContext(ctx)
	logger.Debug("创建 dashboardJson 请求: ", dashboardJSON)

	resp, err := http.Post(GetGrafanaURL()+"/api/dashboards/db", "application/json", bytes.NewBufferString(dashboardJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logger.Debug("创建 dashboardJson 返回结果: ", string(body))

	return nil
}

// CreateDashboardByBody creates a Grafana dashboard from a JSON body.
func CreateDashboardByBody(ctx context.Context, uidsTr, datasourceName, body string) (string, error) {
	var bodyJSON map[string]any
	if err := json.Unmarshal([]byte(body), &bodyJSON); err != nil {
		return "", err
	}

	if uid, ok := bodyJSON["uid"]; ok && uid != nil {
		_ = DeleteDashboard(ctx, uid.(string))
	}

	if datasourceName != "" {
		// Update datasource references in panels
		if panels, ok := bodyJSON["panels"].([]any); ok {
			for _, panel := range panels {
				if p, ok := panel.(map[string]any); ok {
					if _, hasDatasource := p["datasource"]; hasDatasource {
						p["datasource"] = datasourceName
					}
				}
			}
		}

		// Update datasource references in templating
		if templating, ok := bodyJSON["templating"].(map[string]any); ok {
			if list, ok := templating["list"].([]any); ok {
				for _, item := range list {
					if l, ok := item.(map[string]any); ok {
						if _, hasDatasource := l["datasource"]; hasDatasource {
							l["datasource"] = datasourceName
						}
					}
				}
			}
		}
	}

	var newBody []byte
	if _, has := bodyJSON["dashboard"]; has {
		newBody = []byte(body)
	} else {
		newBodyMap := map[string]any{
			"dashboard": bodyJSON,
		}
		newBody, _ = json.Marshal(newBodyMap)
	}

	logger := logx.WithContext(ctx)
	logger.Infof("创建 dashboardJson 请求: %s", string(newBody))

	resp, err := http.Post(GetGrafanaURL()+"/api/dashboards/db", "application/json", bytes.NewBuffer(newBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	logger.Infof("创建 dashboardJson 返回结果: %s", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}
	return string(respBody), nil
}

// GetDatasourceUUIDByJDBC generates a datasource UUID from JDBC type.
func GetDatasourceUUIDByJDBC(jdbcType types.SrcJdbcType) string {
	return GetDashboardUUIDByAlias(jdbcType.Alias())
}

// Fields2Columns converts field definitions to column string for Grafana.
func Fields2Columns(jdbcType types.SrcJdbcType, fields []*types.FieldDefine) string {
	// TDengine uses `, PostgreSQL and TimescaleDB use "
	flag := "`"
	if jdbcType.Id() != types.SrcJdbcTypeTdEngine.Id() {
		flag = `\"`
	}

	var fieldNames []string
	for _, field := range fields {
		// Filter out BLOB types and specific system fields
		if field.Type == types.FieldTypeBlob || field.Type == types.FieldTypeLBlob {
			continue
		}
		if field.Name == constants.QosField ||
			field.Name == constants.SysSaveTime ||
			field.Name == "tag" ||
			field.Name == constants.SysFieldID {
			continue
		}
		fieldNames = append(fieldNames, flag+field.Name+flag)
	}

	return strings.Join(fieldNames, ", ")
}

// CreateTimeSeriesListDashboard creates a time series list dashboard with multiple panels.
func CreateTimeSeriesListDashboard(ctx context.Context, srcJdbcType types.SrcJdbcType, topics []*types.CreateTopicDto, dashboardName string) (string, error) {
	logger := logx.WithContext(ctx)
	logger.Infof("调用 创建时序组合Dashboard: %s", dashboardName)

	var panelTemplate string
	if srcJdbcType.Id() == types.SrcJdbcTypeTimeScaleDB.Id() {
		panelTemplate = LoadTemplate(ctx, "templates/ts-panel.json")
	} else {
		panelTemplate = LoadTemplate(ctx, "templates/td-panel.json")
	}

	var panelJSONList []any
	for i, topic := range topics {
		columns := Fields2Columns(srcJdbcType, topic.Fields)
		title := topic.GetTopic()
		schema := "public"
		table := topic.TableName

		if dot := strings.Index(table, "."); dot > 0 {
			schema = table[:dot]
			table = table[dot+1:]
		}

		// Panel's x-axis position
		gridPosX := i * 8
		if gridPosX > 16 {
			gridPosX = (i % 3) * 8
		}

		panelParam := map[string]any{
			"id":            i + 1,
			"title":         title,
			"dataSourceUid": GetDatasourceUUIDByJDBC(srcJdbcType),
			"columns":       columns,
			"schema":        schema,
			"tableName":     table,
			"gridPosX":      gridPosX,
		}

		panelJSON := FormatTemplateMap(panelTemplate, panelParam)
		var panel any
		json.Unmarshal([]byte(panelJSON), &panel)
		panelJSONList = append(panelJSONList, panel)
	}

	template := LoadTemplate(ctx, "templates/td-dashboard-list.json")
	uid := uuid.New().String()[:32] // Fast simple UUID

	dbParams := map[string]any{
		"uid":    uid,
		"title":  dashboardName,
		"panels": panelJSONList,
	}

	dashboardJSON := FormatTemplateMap(template, dbParams)
	logger.Debugf("创建时序组合DashboardDashboard 请求: %s", dashboardJSON)

	resp, err := http.Post(GetGrafanaURL()+"/api/dashboards/db", "application/json", bytes.NewBufferString(dashboardJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logger.Infof("创建时序组合Dashboard 返回结果: %s", string(body))

	return uid, nil
}

// GetDashboardByUUID retrieves a dashboard by UUID.
func GetDashboardByUUID(ctx context.Context, uuid string) (map[string]any, error) {
	logger := logx.WithContext(ctx)
	url := GetGrafanaURL() + "/api/dashboards/uid/" + uuid
	resp, err := http.Get(url)
	if err != nil {
		logger.Errorf("查询 dashboards 失败: %v, %s", err, url)
		return nil, err
	}
	logger.Debugf("查询 dashboards 请求: %s", url)

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logger.Debugf("查询 dashboards 返回结果: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateFolder creates a Grafana folder.
func CreateFolder(uid, title string) (*grafanadto.GrafanaFolderDto, error) {
	// Check if folder exists
	url := GetGrafanaURL() + "/api/folders/" + uid
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var folder grafanadto.GrafanaFolderDto
		err = json.Unmarshal(body, &folder)
		if err != nil {
			return nil, err
		}
		return &folder, nil
	}

	// Create new folder
	reqBody := map[string]any{
		"uid":   uid,
		"title": title,
	}
	reqJSON, _ := json.Marshal(reqBody)

	resp, err = http.Post(GetGrafanaURL()+"/api/folders", "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var folder grafanadto.GrafanaFolderDto
		err = json.Unmarshal(body, &folder)
		if err != nil {
			return nil, err
		}
		return &folder, nil
	}

	return nil, fmt.Errorf("failed to create folder: status %d", resp.StatusCode)
}

// SetLanguage sets Grafana language preference.
func SetLanguage(ctx context.Context, language string) error {
	logger := logx.WithContext(ctx)
	url := GetGrafanaURL() + "/api/org/preferences"
	reqBody := fmt.Sprintf(`{"language":"%s"}`, language)

	logger.Infof("设置grafana 语言 请求: %s", reqBody)

	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logger.Infof("设置grafana 语言 返回结果: %s", string(body))

	return nil
}

// Create creates a Grafana dashboard from JSON.
func Create(ctx context.Context, dashboardJSON string) (bool, error) {
	logger := logx.WithContext(ctx)
	logger.Debugf("grafana 创建 dashboards 请求: %s", dashboardJSON)

	resp, err := http.Post(GetGrafanaURL()+"/api/dashboards/db", "application/json", bytes.NewBufferString(dashboardJSON))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logger.Debugf("grafana 创建 dashboards 返回结果: %s", string(body))

	return resp.StatusCode == http.StatusOK, nil
}

// Get retrieves a Grafana dashboard by UUID.
func Get(ctx context.Context, uuid string) (string, error) {
	logger := logx.WithContext(ctx)
	url := GetGrafanaURL() + "/api/dashboards/uid/" + uuid
	logger.Debugf("grafana 查询 dashboards 请求: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logger.Debugf("grafana 查询 dashboards 返回结果: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}
	return string(body), nil
}

// Helper functions

// LoadTemplate loads a template file.
func LoadTemplate(ctx context.Context, path string) string {
	logger := logx.WithContext(ctx)
	bs, er := templates.ReadFile(path)
	if er != nil {
		logger.Error("FileNotExist:", path)
		return ""
	}
	return string(bs)
}

// formatTemplate formats a template string with struct values.
func formatTemplate(template string, data any) string {
	// Convert struct to map and use FormatTemplateMap
	jsonData, _ := json.Marshal(data)
	var dataMap map[string]any
	err := json.Unmarshal(jsonData, &dataMap)
	if err != nil {
		return ""
	}
	return FormatTemplateMap(template, dataMap)
}

// FormatTemplateMap formats a template string with map values.
func FormatTemplateMap(template string, data map[string]any) string {
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{%s}", key)
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		default:
			valueStr = fmt.Sprint(v)
		}
		result = strings.ReplaceAll(result, placeholder, valueStr)
	}
	return result
}
