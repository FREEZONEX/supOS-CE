package audit

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"backend/internal/contextx"
	"backend/internal/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type ctxKey string

const stateKey ctxKey = "audit.state"
const auditResultMsgKey = "_resultMsg"

const (
	ScopeTypePlatform = 1
	ScopeTypeProject  = 2
)

const (
	ResTypeAuth           = "Auth"
	ResTypeUserManagement = "UserManagement"
	ResTypeRole           = "Role"
	ResTypeUNS            = "UNS"
	ResTypeResource       = "Resource"
	ResTypeTemplate       = "Template"
	ResTypeLabel          = "Label"
	ResTypeGroup          = "Group"
	ResTypeDashboard      = "Dashboard"
	ResTypeAppKey         = "AppKey"
	ResTypeSourceFlow     = "SourceFlow"
	ResTypeEventFlow      = "EventFlow"
	ResTypeAsset          = "Asset"
	ResTypePlatform       = "Platform"
)

const (
	BizCreate        = "Create"
	BizUpdate        = "Update"
	BizDelete        = "Delete"
	BizEnable        = "Enable"
	BizDisable       = "Disable"
	BizImport        = "Import"
	BizExport        = "Export"
	BizStart         = "Start"
	BizStop          = "Stop"
	BizLogin         = "Login"
	BizLogout        = "Logout"
	BizResetPassword = "ResetPassword"
	BizAssignRole    = "AssignRole"
	BizReplace       = "Replace"
	BizActivate      = "Activate"
	BizInstall       = "Install"
	BizUninstall     = "Uninstall"
	BizSave          = "Save"
	BizCopy          = "Copy"
	BizDeploy        = "Deploy"
	BizMark          = "Mark"
	BizUnmark        = "Unmark"
	BizBindUNS       = "BindUNS"
	BizConfirm       = "Confirm"
	BizRestore       = "Restore"
	BizView          = "View"
)

const (
	ShowInRecentYes int64 = 1
	ShowInRecentNo  int64 = 2
)

type RecordInput struct {
	ScopeType      int
	ScopeID        int64
	ScopeName      string
	ResType        string
	ResID          string
	ResName        string
	BusinessType   string
	IsShowInRecent int64
	Detail         map[string]any
}

type Operator struct {
	ID    string
	Name  string
	Email string
}

type State struct {
	Record         *RecordInput
	Operator       *Operator
	RequestPayload map[string]any
	Suppressed     bool
}

type Service struct {
	repo  *repo.AuditRepo
	queue chan repo.AuditLog
	once  sync.Once
}

type PageQuery struct {
	PageNo          int
	PageSize        int
	ScopeType       int
	ScopeID         int64
	ResType         string
	BusinessType    string
	OperatorKeyword string
	Code            *int
	StartTime       string
	EndTime         string
}

type ResponseMeta struct {
	StatusCode int
	Status     string
}

func New(ctx context.Context) *Service {
	auditRepo := repo.NewAuditRepo(ctx)
	svc := &Service{
		repo:  auditRepo,
		queue: make(chan repo.AuditLog, 256),
	}
	svc.start()
	return svc
}

func (s *Service) Page(ctx context.Context, query PageQuery) (map[string]any, error) {
	filter := repo.AuditLogFilter{
		PageNo:           query.PageNo,
		PageSize:         query.PageSize,
		ScopeType:        query.ScopeType,
		ScopeID:          query.ScopeID,
		ResType:          query.ResType,
		BusinessTypeCode: repo.AuditBusinessTypeCode(query.BusinessType),
		OperatorKeyword:  query.OperatorKeyword,
		Code:             query.Code,
		StartTime:        parseTime(query.StartTime),
		EndTime:          parseTime(query.EndTime),
	}
	items, total, err := s.repo.ListAuditLogs(ctx, filter)
	if err != nil {
		return nil, err
	}
	pageNo := query.PageNo
	if pageNo <= 0 {
		pageNo = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	return map[string]any{
		"pageNo":   pageNo,
		"pageSize": pageSize,
		"total":    total,
		"data":     items,
	}, nil
}

func (s *Service) Detail(ctx context.Context, id int64) (repo.AuditLog, error) {
	item, err := s.repo.GetAuditLog(ctx, id)
	if errors.Is(err, repo.ErrNotFound) {
		return repo.AuditLog{}, ErrNotFound
	}
	return item, err
}

func (s *Service) WithState(ctx context.Context) context.Context {
	if StateFromContext(ctx) != nil {
		return ctx
	}
	return context.WithValue(ctx, stateKey, &State{})
}

func StateFromContext(ctx context.Context) *State {
	if ctx == nil {
		return nil
	}
	if state, ok := ctx.Value(stateKey).(*State); ok {
		return state
	}
	return nil
}

func (s *Service) Record(ctx context.Context, input RecordInput) {
	state := StateFromContext(ctx)
	if state == nil {
		return
	}
	copied := normalizeRecord(input)
	state.Record = &copied
}

func (s *Service) Suppress(ctx context.Context) {
	state := StateFromContext(ctx)
	if state != nil {
		state.Suppressed = true
	}
}

func (s *Service) BindSubject(ctx context.Context, subject contextx.Subject) {
	state := StateFromContext(ctx)
	if state == nil {
		return
	}
	name := strings.TrimSpace(subject.UserName)
	if subject.AuthType == "apiKey" && strings.TrimSpace(subject.APIKeyName) != "" {
		name = strings.TrimSpace(subject.APIKeyName)
	}
	state.Operator = &Operator{
		ID:    strconv.FormatInt(subject.UserID, 10),
		Name:  truncateString(name, 200),
		Email: truncateString(subject.Email, 255),
	}
}

func (s *Service) BindRequestPayload(ctx context.Context, payload map[string]any) {
	state := StateFromContext(ctx)
	if state == nil || len(payload) == 0 {
		return
	}
	state.RequestPayload = payload
}

func (s *Service) FlushRequest(r *http.Request, meta ResponseMeta) {
	item := s.BuildAuditLog(r, meta)
	if item == nil {
		return
	}
	s.dispatch(*item)
}

func (s *Service) BuildAuditLog(r *http.Request, meta ResponseMeta) *repo.AuditLog {
	if r == nil {
		return nil
	}
	state := StateFromContext(r.Context())
	var record *RecordInput
	var operator *Operator
	if state != nil {
		if state.Suppressed {
			return nil
		}
		record = state.Record
		operator = state.Operator
	}
	if record == nil {
		record = autoRecord(r)
	}
	if record != nil {
		normalized := normalizeRecord(*record)
		record = &normalized
	}
	if !shouldRecordAudit(r.Method, record) {
		return nil
	}
	requestPayload := requestPayloadFromContext(r.Context())
	if strings.TrimSpace(record.ResName) == "" {
		record.ResName = resourceNameFromPayload(requestPayload)
	}
	if strings.TrimSpace(record.ResName) == "" {
		record.ResName = record.ResID
	}

	code := meta.StatusCode
	if code <= 0 {
		code = http.StatusOK
	}
	detail := sanitizeDetail(record.Detail)
	if len(detail) == 0 {
		detail = map[string]any{}
	}
	if _, ok := detail["path"]; !ok {
		detail["path"] = r.URL.Path
	}
	if _, ok := detail["method"]; !ok {
		detail["method"] = r.Method
	}
	if r.URL.RawQuery != "" {
		detail["query"] = truncateString(r.URL.RawQuery, 1024)
	}
	if status := strings.TrimSpace(meta.Status); status != "" {
		detail[auditResultMsgKey] = status
	} else {
		detail[auditResultMsgKey] = http.StatusText(code)
	}

	item := &repo.AuditLog{
		UserID:           operatorUserID(operator),
		ScopeType:        record.ScopeType,
		ScopeID:          record.ScopeID,
		ScopeName:        truncateString(record.ScopeName, 200),
		ResType:          truncateString(record.ResType, 50),
		ResID:            truncateString(record.ResID, 50),
		ResName:          truncateString(record.ResName, 256),
		BusinessType:     truncateString(record.BusinessType, 64),
		BusinessTypeCode: repo.AuditBusinessTypeCode(record.BusinessType),
		IsShowInRecent:   showInRecentForRecord(record),
		Code:             int64(code),
		DetailJSON:       marshalDetail(detail),
		CreatedTime:      time.Now().UTC(),
	}
	if operator != nil {
		item.OperatorName = truncateString(operator.Name, 200)
		item.OperatorEmail = truncateString(operator.Email, 255)
	}
	return item
}

func (s *Service) start() {
	if s == nil {
		return
	}
	s.once.Do(func() {
		go func() {
			for item := range s.queue {
				if s.repo == nil {
					continue
				}
				if err := s.repo.InsertAuditLog(context.Background(), item); err != nil {
					logx.Errorf("insert audit log failed: %v", err)
				}
			}
		}()
	})
}

func (s *Service) dispatch(item repo.AuditLog) {
	if s == nil {
		return
	}
	s.start()
	select {
	case s.queue <- item:
	default:
		logx.Errorf("audit log queue full, dropping item resType=%s businessType=%s resID=%s", item.ResType, item.BusinessType, item.ResID)
	}
}

func operatorUserID(operator *Operator) int64 {
	if operator == nil {
		return 0
	}
	userID, err := strconv.ParseInt(strings.TrimSpace(operator.ID), 10, 64)
	if err != nil || userID < 0 {
		return 0
	}
	return userID
}

func autoRecord(r *http.Request) *RecordInput {
	if r == nil || !shouldAutoAudit(r) {
		return nil
	}
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	core := stripAPIPrefix(parts)
	if len(core) == 0 {
		return nil
	}
	resType := resourceTypeFromPath(core)
	bizType := businessTypeFromPath(r.Method, core)
	if resType == "" || bizType == "" {
		return nil
	}
	return &RecordInput{
		ScopeType:    ScopeTypePlatform,
		ResType:      resType,
		ResID:        resourceIDFromPath(core),
		ResName:      resourceNameFromPayload(requestPayloadFromContext(r.Context())),
		BusinessType: bizType,
		Detail:       map[string]any{"auto": true},
	}
}

func requestPayloadFromContext(ctx context.Context) map[string]any {
	state := StateFromContext(ctx)
	if state == nil {
		return nil
	}
	return state.RequestPayload
}

func shouldAutoAudit(r *http.Request) bool {
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	}
	path := strings.ToLower(strings.TrimSpace(r.URL.Path))
	if path == "" {
		return false
	}
	skips := []string{
		"/api/core/uns/newmsg",
		"/api/core/system/",
		"/api/core/assets",
		"/api/core/asset-bindings",
		"/api/openapi/",
		"/openapi/",
		"/healthz",
		"/readyz",
		"/metrics",
	}
	for _, skip := range skips {
		if strings.HasPrefix(path, skip) {
			return false
		}
	}
	// Vision worker 面的机器调用(注册/结果上报/心跳)按秒级频率触发,
	// 不是用户操作,不进审计;heartbeat 带路径参数,只能按后缀识别。
	if strings.HasPrefix(path, "/api/core/vision/worker/") ||
		strings.HasPrefix(path, "/api/core/vision/ingest/") ||
		(strings.HasPrefix(path, "/api/core/vision/tasks/") && strings.HasSuffix(path, "/heartbeat")) {
		return false
	}
	return strings.HasPrefix(path, "/api/core/")
}

func stripAPIPrefix(parts []string) []string {
	if len(parts) >= 3 && parts[0] == "api" && parts[1] == "core" {
		return parts[2:]
	}
	return parts
}

func resourceTypeFromPath(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	first := strings.ToLower(parts[0])
	second := ""
	if len(parts) > 1 {
		second = strings.ToLower(parts[1])
	}
	switch first {
	case "auth":
		return ResTypeAuth
	case "users", "user":
		return ResTypeUserManagement
	case "roles", "role":
		return ResTypeRole
	case "resources", "resource", "menus", "menu":
		return ResTypeResource
	case "uns":
		switch second {
		case "templates":
			return ResTypeTemplate
		case "labels":
			return ResTypeLabel
		default:
			return ResTypeUNS
		}
	case "flows", "flow":
		if len(parts) > 1 && (parts[1] == "event" || parts[1] == "event-flow") {
			return ResTypeEventFlow
		}
		return ResTypeSourceFlow
	case "dashboards", "dashboard":
		return ResTypeDashboard
	case "apikeys", "api-keys", "apikey":
		return ResTypeAppKey
	case "assets", "asset":
		return ResTypeAsset
	default:
		return toPascal(first)
	}
}

func businessTypeFromPath(method string, parts []string) string {
	joined := strings.ToLower(strings.Join(parts, "/"))
	switch {
	case strings.Contains(joined, "login"):
		return BizLogin
	case strings.Contains(joined, "logout"):
		return BizLogout
	case strings.Contains(joined, "reset"):
		return BizResetPassword
	case strings.Contains(joined, "enable"):
		return BizEnable
	case strings.Contains(joined, "disable"):
		return BizDisable
	case strings.Contains(joined, "import"):
		return BizImport
	case strings.Contains(joined, "export"):
		return BizExport
	case strings.Contains(joined, "deploy"):
		return BizDeploy
	case strings.Contains(joined, "restore") || strings.Contains(joined, "revert"):
		return BizRestore
	case strings.Contains(joined, "copy") || strings.Contains(joined, "clone"):
		return BizCopy
	case strings.Contains(joined, "mark"):
		return BizMark
	case strings.Contains(joined, "bind"):
		return BizBindUNS
	case strings.Contains(joined, "save"):
		return BizSave
	case strings.Contains(joined, "install"):
		return BizInstall
	case strings.Contains(joined, "uninstall"):
		return BizUninstall
	}
	switch method {
	case http.MethodPost:
		return BizCreate
	case http.MethodPut, http.MethodPatch:
		return BizUpdate
	case http.MethodDelete:
		return BizDelete
	default:
		return ""
	}
}

func resourceIDFromPath(parts []string) string {
	for i := len(parts) - 1; i >= 0; i-- {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}
		lower := strings.ToLower(part)
		if lower == "create" || lower == "update" || lower == "delete" || lower == "enable" || lower == "disable" ||
			lower == "import" || lower == "export" || lower == "deploy" || lower == "save" || lower == "reset" ||
			lower == "list" || lower == "page" || lower == "version" || lower == "data" || lower == "status" ||
			lower == "content" || lower == "snapshots" {
			continue
		}
		return truncateString(part, 128)
	}
	return ""
}

func resourceNameFromPath(parts []string) string {
	id := resourceIDFromPath(parts)
	if id == "" {
		return strings.Join(parts, "/")
	}
	return id
}

func resourceNameFromPayload(payload map[string]any) string {
	if len(payload) == 0 {
		return ""
	}
	keys := []string{
		"name",
		"flowName",
		"userName",
		"username",
		"nickName",
		"clientName",
		"apiKeyName",
		"tokenName",
		"nodeName",
		"userId",
		"resourceKey",
		"routeKey",
		"fileName",
		"originalName",
		"attachmentName",
		"versionName",
		"newName",
		"keyName",
		"title",
		"label",
		"namespace",
		"alias",
		"path",
	}
	for _, key := range keys {
		if value := payloadStringValue(payload, key); value != "" {
			return truncateString(value, 256)
		}
	}
	return ""
}

func payloadStringValue(payload map[string]any, key string) string {
	for currentKey, value := range payload {
		if !strings.EqualFold(strings.TrimSpace(currentKey), key) {
			continue
		}
		switch typed := value.(type) {
		case string:
			return strings.TrimSpace(typed)
		case json.Number:
			return strings.TrimSpace(typed.String())
		case float64:
			if typed == float64(int64(typed)) {
				return strconv.FormatInt(int64(typed), 10)
			}
			return strconv.FormatFloat(typed, 'f', -1, 64)
		case int64:
			return strconv.FormatInt(typed, 10)
		case int:
			return strconv.Itoa(typed)
		}
	}
	return ""
}

func shouldRecordAudit(method string, record *RecordInput) bool {
	if record == nil {
		return false
	}
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return strings.TrimSpace(record.ResType) == ResTypeAuth && strings.TrimSpace(record.BusinessType) == BizLogin
	}
	return strings.TrimSpace(record.ResType) != "" && strings.TrimSpace(record.BusinessType) != ""
}

func normalizeRecord(input RecordInput) RecordInput {
	if input.ScopeType == 0 {
		input.ScopeType = ScopeTypePlatform
	}
	input.ResType = strings.TrimSpace(input.ResType)
	input.ResID = strings.TrimSpace(input.ResID)
	input.ResName = strings.TrimSpace(input.ResName)
	input.BusinessType = strings.TrimSpace(input.BusinessType)
	if input.Detail == nil {
		input.Detail = map[string]any{}
	}
	return input
}

func showInRecentForRecord(record *RecordInput) int64 {
	if record == nil {
		return ShowInRecentNo
	}
	switch record.IsShowInRecent {
	case ShowInRecentYes, ShowInRecentNo:
		return record.IsShowInRecent
	}
	if hiddenRecentResourceType(record.ResType) {
		return ShowInRecentNo
	}
	return ShowInRecentYes
}

func hiddenRecentResourceType(resType string) bool {
	switch strings.TrimSpace(resType) {
	case ResTypeUNS,
		ResTypeTemplate,
		ResTypeLabel,
		ResTypeGroup,
		ResTypeDashboard,
		ResTypeSourceFlow,
		ResTypeEventFlow:
		return false
	default:
		return true
	}
}

func requestID(r *http.Request) string {
	if r == nil {
		return ""
	}
	for _, key := range []string{"X-Request-Id", "X-Request-ID", "Traceparent"} {
		if value := strings.TrimSpace(r.Header.Get(key)); value != "" {
			return truncateString(value, 64)
		}
	}
	return ""
}

func clientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	for _, key := range []string{"X-Forwarded-For", "X-Real-IP"} {
		value := strings.TrimSpace(r.Header.Get(key))
		if value == "" {
			continue
		}
		if key == "X-Forwarded-For" {
			value = strings.TrimSpace(strings.Split(value, ",")[0])
		}
		return truncateString(value, 64)
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return truncateString(host, 64)
	}
	return truncateString(strings.TrimSpace(r.RemoteAddr), 64)
}

func sanitizeDetail(detail map[string]any) map[string]any {
	if len(detail) == 0 {
		return map[string]any{}
	}
	sanitized := make(map[string]any, len(detail))
	for key, value := range detail {
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "password", "newpassword", "client_secret", "clientsecret", "authorization", "token", "cookie", "license", "licensetoken", "license_token":
			continue
		default:
			sanitized[key] = value
		}
	}
	return sanitized
}

func marshalDetail(detail map[string]any) string {
	raw, err := json.Marshal(detail)
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func truncateString(value string, limit int) string {
	value = strings.TrimSpace(value)
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit]
}

func toPascal(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	seps := func(r rune) bool { return r == '-' || r == '_' || r == ' ' }
	parts := strings.FieldsFunc(value, seps)
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

func parseTime(value string) int64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if n, err := strconv.ParseInt(value, 10, 64); err == nil {
		return n
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t.UnixMilli()
	}
	return 0
}

var ErrNotFound = errors.New("audit log not found")
