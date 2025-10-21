package grafanautil

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"backend/internal/common"
	"backend/internal/common/constants"
	"backend/internal/common/dto"
	grafanadto "backend/internal/common/dto/grafana"
	"backend/internal/common/enums"
	"backend/internal/common/utils/runtimeutil"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/google/uuid"
)

// GrafanaUtils provides utility functions for Grafana operations.
//
// GetGrafanaURL returns the Grafana URL based on the runtime environment.
func GetGrafanaURL() string {
	if runtimeutil.IsLocalEnv() {
		return "http://100.100.100.22:33997/grafana/home"
	}
	return "http://grafana:3000"
}

// GetDashboardUUIDByAlias generates a dashboard UUID from an alias using MD5.
func GetDashboardUUIDByAlias(alias string) string {
	hash := md5.Sum([]byte(alias))
	return hex.EncodeToString(hash[:8]) // 16-character hex (MD5 digestHex16)
}

// DeleteDashboard deletes a Grafana dashboard by UID.
func DeleteDashboard(uid string) error {
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
	logx.Debugf("Delete Dashboard: %s, response: %s", uid, string(body))
	return nil
}

// DeleteDatasource deletes a Grafana datasource by UID.
func DeleteDatasource(uid string) error {
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
	logx.Debugf("Delete DataSource: %s, response: %s", uid, string(body))
	return nil
}

// CreateDatasource creates a Grafana datasource.
func CreateDatasource(jdbcType *common.SrcJdbcType, username, password string, reCreate bool) (bool, error) {
	title := jdbcType.Alias
	datasource := &grafanadto.GrafanaDataSourceDto{
		User:     username,
		Password: password,
		UID:      GetDatasourceUUIDByJDBC(jdbcType),
		Name:     title,
	}

	if reCreate {
		// Delete first, then create
		_ = DeleteDatasource(datasource.UID)
	}

	var dsTemplate string
	switch jdbcType.ID {
	case common.SrcJdbcTypePostgresql.ID:
		dsTemplate = loadTemplate("templates/pg-datasource.json")
		datasource.URL = constants.PGJDBCURL
	case common.SrcJdbcTypeTimeScaleDB.ID:
		dsTemplate = loadTemplate("templates/pg-datasource.json")
		datasource.URL = constants.TSDBJDBCURL
	case common.SrcJdbcTypeTdEngine.ID:
		datasource.URL = constants.TDJDBCURL
		datasource.CreateBasicAuth()
		dsTemplate = loadTemplate("templates/td-datasource.json")
	default:
		return false, fmt.Errorf("unsupported JDBC type: %d", jdbcType.ID)
	}

	dsJSON := formatTemplate(dsTemplate, datasource)
	logx.Infof("创建 datasource 请求: %s", dsJSON)

	resp, err := http.Post(GetGrafanaURL()+"/api/datasources", "application/json", bytes.NewBufferString(dsJSON))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Infof("创建 datasource 返回结果: %s", string(body))

	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusConflict, nil
}

// CreateDatasourceByBody creates a Grafana datasource from a JSON body.
func CreateDatasourceByBody(name, body string, reCreate bool) (string, error) {
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
	logx.Infof("创建 datasource 请求: %s", string(newBody))

	resp, err := http.Post(GetGrafanaURL()+"/api/datasources", "application/json", bytes.NewBuffer(newBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	logx.Infof("创建 datasource 返回结果: %s", string(respBody))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		return "", nil
	}
	return string(respBody), nil
}

// CreateDashboard creates a Grafana dashboard.
func CreateDashboard(table, tagNameCondition string, jdbcType *common.SrcJdbcType, schema, title, columns, ct string) (string, error) {
	uid := GetDashboardUUIDByAlias(title)
	var template string
	var dbParams map[string]any

	switch jdbcType.ID {
	case common.SrcJdbcTypePostgresql.ID:
		template = loadTemplate("templates/pg-dashboard.json")
		dbParams = map[string]any{
			"title":          title,
			"uid":            uid,
			"dataSourceType": jdbcType.DataSrcType,
			"dataSourceUid":  GetDatasourceUUIDByJDBC(jdbcType),
			"schema":         schema,
			"tableName":      table,
			"columns":        columns,
		}
	case common.SrcJdbcTypeTdEngine.ID:
		template = loadTemplate("templates/td-dashboard.json")
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
	case common.SrcJdbcTypeTimeScaleDB.ID:
		template = loadTemplate("templates/ts-dashboard.json")
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
		return "", fmt.Errorf("unsupported JDBC type: %d", jdbcType.ID)
	}

	dbParams["sys_field_create_time"] = ct
	dashboardJSON := formatTemplateMap(template, dbParams)

	logx.Debugf("创建 dashboardJson 请求: %s", dashboardJSON)

	resp, err := http.Post(GetGrafanaURL()+"/api/dashboards/db", "application/json", bytes.NewBufferString(dashboardJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Debugf("创建 dashboardJson 返回结果: %s", string(body))

	return uid, nil
}

// CreateDashboardByBody creates a Grafana dashboard from a JSON body.
func CreateDashboardByBody(uidsTr, datasourceName, body string) (string, error) {
	var bodyJSON map[string]any
	if err := json.Unmarshal([]byte(body), &bodyJSON); err != nil {
		return "", err
	}

	if uid, ok := bodyJSON["uid"]; ok && uid != nil {
		_ = DeleteDashboard(uid.(string))
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

	newBodyMap := map[string]any{
		"dashboard": bodyJSON,
	}
	newBody, _ := json.Marshal(newBodyMap)

	logx.Infof("创建 dashboardJson 请求: %s", string(newBody))

	resp, err := http.Post(GetGrafanaURL()+"/api/dashboards/db", "application/json", bytes.NewBuffer(newBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	logx.Infof("创建 dashboardJson 返回结果: %s", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}
	return string(respBody), nil
}

// GetDataSourceByName retrieves a Grafana datasource by name.
func GetDataSourceByName(name string) (string, error) {
	url := GetGrafanaURL() + "/api/datasources/name/" + name
	logx.Debugf("查询 datasource 请求: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Debugf("查询 datasource 返回结果: %s", string(body))

	return string(body), nil
}

// GetDatasourceUUIDByJDBC generates a datasource UUID from JDBC type.
func GetDatasourceUUIDByJDBC(jdbcType *common.SrcJdbcType) string {
	hash := md5.Sum([]byte(jdbcType.Alias))
	return hex.EncodeToString(hash[:])
}

// Fields2Columns converts field definitions to column string for Grafana.
func Fields2Columns(jdbcType *common.SrcJdbcType, fields []*dto.FieldDefine) string {
	// TDengine uses `, PostgreSQL and TimescaleDB use "
	flag := "`"
	if jdbcType.ID != common.SrcJdbcTypeTdEngine.ID {
		flag = `\"`
	}

	var fieldNames []string
	for _, field := range fields {
		// Filter out BLOB types and specific system fields
		if field.Type == enums.FieldTypeBlob || field.Type == enums.FieldTypeLBlob {
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
func CreateTimeSeriesListDashboard(srcJdbcType *common.SrcJdbcType, topics []*dto.CreateTopicDto, dashboardName string) (string, error) {
	logx.Infof("调用 创建时序组合Dashboard: %s", dashboardName)

	var panelTemplate string
	if srcJdbcType.ID == common.SrcJdbcTypeTimeScaleDB.ID {
		panelTemplate = loadTemplate("templates/ts-panel.json")
	} else {
		panelTemplate = loadTemplate("templates/td-panel.json")
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

		panelJSON := formatTemplateMap(panelTemplate, panelParam)
		var panel any
		json.Unmarshal([]byte(panelJSON), &panel)
		panelJSONList = append(panelJSONList, panel)
	}

	template := loadTemplate("templates/td-dashboard-list.json")
	uid := uuid.New().String()[:32] // Fast simple UUID

	dbParams := map[string]any{
		"uid":    uid,
		"title":  dashboardName,
		"panels": panelJSONList,
	}

	dashboardJSON := formatTemplateMap(template, dbParams)
	logx.Debugf("创建时序组合DashboardDashboard 请求: %s", dashboardJSON)

	resp, err := http.Post(GetGrafanaURL()+"/api/dashboards/db", "application/json", bytes.NewBufferString(dashboardJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Infof("创建时序组合Dashboard 返回结果: %s", string(body))

	return uid, nil
}

// GetDashboardByUUID retrieves a dashboard by UUID.
func GetDashboardByUUID(uuid string) (map[string]any, error) {
	url := GetGrafanaURL() + "/api/dashboards/uid/" + uuid
	logx.Debugf("查询 dashboards 请求: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Debugf("查询 dashboards 返回结果: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]any
	json.Unmarshal(body, &result)
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
		json.Unmarshal(body, &folder)
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
		json.Unmarshal(body, &folder)
		return &folder, nil
	}

	return nil, fmt.Errorf("failed to create folder: status %d", resp.StatusCode)
}

// SetLanguage sets Grafana language preference.
func SetLanguage(language string) error {
	url := GetGrafanaURL() + "/api/org/preferences"
	reqBody := fmt.Sprintf(`{"language":"%s"}`, language)

	logx.Infof("设置grafana 语言 请求: %s", reqBody)

	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Infof("设置grafana 语言 返回结果: %s", string(body))

	return nil
}

// Create creates a Grafana dashboard from JSON.
func Create(dashboardJSON string) (bool, error) {
	logx.Debugf("grafana 创建 dashboards 请求: %s", dashboardJSON)

	resp, err := http.Post(GetGrafanaURL()+"/api/dashboards/db", "application/json", bytes.NewBufferString(dashboardJSON))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Debugf("grafana 创建 dashboards 返回结果: %s", string(body))

	return resp.StatusCode == http.StatusOK, nil
}

// Get retrieves a Grafana dashboard by UUID.
func Get(uuid string) (string, error) {
	url := GetGrafanaURL() + "/api/dashboards/uid/" + uuid
	logx.Debugf("grafana 查询 dashboards 请求: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Debugf("grafana 查询 dashboards 返回结果: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}
	return string(body), nil
}

// Helper functions

// loadTemplate loads a template file.
// TODO: Implement actual template loading from embedded files or resources directory.
func loadTemplate(path string) string {
	// In Go, we can use embed.FS or read from a resources directory
	// For now, return a placeholder
	return "{}"
}

// formatTemplate formats a template string with struct values.
func formatTemplate(template string, data any) string {
	// Convert struct to map and use formatTemplateMap
	jsonData, _ := json.Marshal(data)
	var dataMap map[string]any
	json.Unmarshal(jsonData, &dataMap)
	return formatTemplateMap(template, dataMap)
}

// formatTemplateMap formats a template string with map values.
func formatTemplateMap(template string, data map[string]any) string {
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{%s}", key)
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		default:
			jsonValue, _ := json.Marshal(v)
			valueStr = string(jsonValue)
		}
		result = strings.ReplaceAll(result, placeholder, valueStr)
	}
	return result
}
