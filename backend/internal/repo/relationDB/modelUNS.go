package relationDB

import (
	"backend/internal/common/constants"
	"backend/internal/common/utils/JsonUtil"
	"backend/internal/types"
	"backend/share/base"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type Example struct {
	ID    int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`    // id编号
	alias int64 `gorm:"column:alias;type:bigint;primary_key;AUTO_INCREMENT"` // 别名
}

func (m *Example) TableName() string {
	return "example"
}

const TableNameUnsNamespace = "uns_namespace"

// UnsNamespace mapped from table <uns_namespace>
type UnsNamespace struct {
	Id               int64            `gorm:"column:id;primaryKey" json:"id"`
	LayRec           string           `gorm:"column:lay_rec;not null" json:"lay_rec"`
	Alias            string           `gorm:"column:alias;not null" json:"alias"`
	ParentAlias      *string          `gorm:"column:parent_alias" json:"parent_alias"`
	Name             string           `gorm:"column:name;not null" json:"name"`
	Path             string           `gorm:"column:path;not null" json:"path"`
	PathType         int16            `gorm:"column:path_type;not null" json:"path_type"`
	DataType         *int16           `gorm:"column:data_type" json:"data_type"`
	ParentDataType   *int16           `gorm:"column:parent_data_type" json:"parent_data_type"`
	Fields           Fields           `gorm:"column:fields;type:json;" json:"fields"`
	CreateAt         time.Time        `gorm:"column:create_at;default:now()" json:"create_at"`
	Status           *int16           `gorm:"column:status;default:1" json:"status"`
	Description      *string          `gorm:"column:description" json:"description"`
	UpdateAt         time.Time        `gorm:"column:update_at" json:"update_at"`
	Protocol         *string          `gorm:"column:protocol" json:"protocol"`
	DataPath         *string          `gorm:"column:data_path" json:"data_path"`
	WithFlags        *int32           `gorm:"column:with_flags" json:"with_flags"`
	DataSrcId        int16            `gorm:"column:data_src_id" json:"data_src_id"`
	RefUns           RefUns           `gorm:"column:ref_uns;default:{};type:jsonb;" json:"ref_uns"`
	Refers           Refers           `gorm:"column:refers;type:json;" json:"refers"`
	Expression       *string          `gorm:"column:expression" json:"expression"`
	TableName_       *string          `gorm:"column:table_name" json:"table_name"`
	NumberFields     *int16           `gorm:"column:number_fields" json:"number_fields"`
	ParentId         *int64           `gorm:"column:parent_id" json:"parent_id"`
	ModelId          *int64           `gorm:"column:model_id" json:"model_id"`
	ProtocolType     *string          `gorm:"column:protocol_type" json:"protocol_type"`
	Extend           map[string]any   `gorm:"column:extend;type:jsonb;serializer:json;" json:"extend"`
	DisplayName      *string          `gorm:"column:display_name" json:"display_name"`
	LabelIds         map[int64]string `gorm:"column:label_ids;type:jsonb;serializer:json;" json:"label_ids"`
	ExtendFieldFlags *int32           `gorm:"column:extend_field_flags" json:"extend_field_flags"`
	MountType        *int16           `gorm:"column:mount_type" json:"mount_type"`
	MountSource      *string          `gorm:"column:mount_source" json:"mount_source"`
	SubscribeAt      *time.Time       `gorm:"column:subscribe_at" json:"subscribe_at"`

	ModelAlias          string `gorm:"-" json:"modelAlias"`
	PathName            string `gorm:"-" json:"pathName"`
	OldPath             string `gorm:"-" json:"oldPath"`
	CountExistsSiblings int64  `gorm:"-" json:"countExistsSiblings"`
}

func (u *UnsNamespace) String() string {
	if u.Id != 0 {
		return strconv.FormatInt(u.Id, 10)
	}
	return u.Alias
}
func (u *UnsNamespace) GetTbFieldName() string {
	if len(u.Fields) == 0 || base.P2v(u.TableName_) == "" {
		return ""
	}
	for _, f := range u.Fields {
		if f.TbValueName != nil {
			return f.Name
		}
	}
	return ""
}

func (u *UnsNamespace) GetLayRec() string {
	return u.LayRec
}

func (u *UnsNamespace) GetDataPath() string {
	return base.P2v(u.DataPath)
}
func (u *UnsNamespace) GetCalculationType() *int32 {
	return nil //TODO
}
func (u *UnsNamespace) GetProtocolMap() (pmap map[string]interface{}) {
	if protocol := u.Protocol; protocol != nil && strings.HasPrefix(*protocol, "{") {
		JsonUtil.FromJson(*protocol, &pmap)
	}
	return pmap
}

func (u *UnsNamespace) GetDescription() string {
	return base.P2v(u.Description)
}

func (u *UnsNamespace) GetExpression() string {
	return base.P2v(u.Expression)
}

func (u *UnsNamespace) GetExtend() map[string]interface{} {
	return u.Extend
}

func (u *UnsNamespace) GetLabelIds() map[int64]string {
	return u.LabelIds
}

func (u *UnsNamespace) GetModelId() *int64 {
	return u.ModelId
}

func (u *UnsNamespace) GetRefers() []*types.InstanceField {
	return u.Refers
}

func (u *UnsNamespace) GetFields() []*types.FieldDefine {
	return u.Fields
}

func (u *UnsNamespace) GetRefUns() map[int64]int {
	return u.RefUns
}

func (u *UnsNamespace) GetCreateAt() int64 {
	return u.CreateAt.UnixMilli()
}

func (u *UnsNamespace) GetUpdateAt() int64 {
	return u.UpdateAt.UnixMilli()
}

func (u *UnsNamespace) GetFlags() *int32 {
	return u.WithFlags
}

func (u *UnsNamespace) GetExtendFieldFlags() *int32 {
	return u.ExtendFieldFlags
}

func (u *UnsNamespace) GetSrcJdbcType() types.SrcJdbcType {
	return types.SrcJdbcType(u.DataSrcId)
}

func (u *UnsNamespace) GetStatus() *int16 {
	return u.Status
}

func (t *UnsNamespace) GetId() int64 {
	return t.Id
}

func (t *UnsNamespace) GetParentId() *int64 {
	return t.ParentId
}

func (t *UnsNamespace) GetAlias() string {
	return t.Alias
}

func (t *UnsNamespace) GetParentAlias() *string {
	return t.ParentAlias
}

func (t *UnsNamespace) GetName() string {
	return t.Name
}

func (t *UnsNamespace) GetDisplayName() string {
	return base.P2v(t.DisplayName)
}

func (t *UnsNamespace) GetPath() string {
	return t.Path
}

func (t *UnsNamespace) GetDataType() *int16 {
	return t.DataType
}
func (t *UnsNamespace) GetParentDataType() *int16 {
	return t.ParentDataType
}
func (t *UnsNamespace) GetPathType() int16 {
	return t.PathType
}

func (t *UnsNamespace) GetMountType() *int16 {
	return t.MountType
}

func (t *UnsNamespace) GetMountSource() string {
	return base.P2v(t.MountSource)
}

type UnsPo struct {
	UnsNamespace
	// 以下都是数据库表不存在的字段，对应 java 的注解： @TableField(exist = false)
	PathName            string `gorm:"->;<-:false;column:path_name" json:"pathName"`
	ModelAlias          string `gorm:"->;<-:false;column:model_alias" json:"modelAlias"`
	TemplateName        string `gorm:"->;<-:false;column:template_name" json:"templateName"`
	TemplateAlias       string `gorm:"->;<-:false;column:template_alias" json:"TemplateAlias"`
	CountChildren       int    `gorm:"->;<-:false;column:count_children" json:"countChildren"`
	CountDirectChildren int    `gorm:"->;<-:false;column:count_direct_children" json:"countDirectChildren"`
	Labels              string `gorm:"->;<-:false;column:labels" json:"labels"`
}
type Fields []*types.FieldDefine

func (f *Fields) Scan(value interface{}) error {
	if value == nil {
		*f = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Fields")
	}
	return json.Unmarshal(bytes, f)
}

func (f Fields) Value() (rs driver.Value, er error) {
	if f == nil {
		return nil, nil
	}
	var bs []byte
	bs, er = json.Marshal(f)
	rs = string(bs)
	return
}

type RefUns map[int64]int

func (f *RefUns) Scan(value interface{}) error {
	if value == nil {
		*f = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan RefUns")
	}
	return json.Unmarshal(bytes, f)
}

func (f RefUns) Value() (rs driver.Value, er error) {
	if f == nil {
		return nil, nil
	}
	var bs []byte
	bs, er = json.Marshal(f)
	rs = string(bs)
	return
}

type Refers []*types.InstanceField

func (f *Refers) Scan(value interface{}) error {
	if value == nil {
		*f = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Refers")
	}
	return json.Unmarshal(bytes, f)
}

func (f Refers) Value() (rs driver.Value, er error) {
	if f == nil {
		return nil, nil
	}
	var bs []byte
	bs, er = json.Marshal(f)
	rs = string(bs)
	return
}

type LabelIds map[int64]string

func (f *LabelIds) Scan(value interface{}) error {
	if value == nil {
		*f = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan LabelIds")
	}
	return json.Unmarshal(bytes, f)
}

func (f LabelIds) Value() (driver.Value, error) {
	if f == nil {
		return nil, nil
	}
	return json.Marshal(f)
}

// TableName UnsNamespace's table name
func (u *UnsNamespace) TableName() string {
	return TableNameUnsNamespace
}

func (u *UnsNamespace) GetTable() string {
	if tb := base.P2v(u.TableName_); tb != "" {
		return tb
	}
	return u.Alias
}
func (u *UnsNamespace) GetID() int64 {
	return u.Id
}
func (u *UnsNamespace) GetParentID() *int64 {
	return u.ParentId
}

const TableNameUnsLabel = "uns_label"

// UnsLabel mapped from table <uns_label>
type UnsLabel struct {
	ID                 int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	LabelName          string     `gorm:"column:label_name" json:"label_name"`
	CreateAt           time.Time  `gorm:"column:create_at;default:now()" json:"create_at"`
	WithFlags          *int32     `gorm:"column:with_flags" json:"with_flags"`
	SubscribeFrequency string     `gorm:"column:subscribe_frequency" json:"subscribe_frequency"`
	SubscribeAt        *time.Time `gorm:"column:subscribe_at" json:"subscribe_at"`
	UpdateAt           time.Time  `gorm:"column:update_at;default:now()" json:"update_at"`
}

// TableName UnsLabel's table name
func (l *UnsLabel) TableName() string {
	return TableNameUnsLabel
}
func (l *UnsLabel) GetId() int64 {
	return l.ID
}

func (l *UnsLabel) GetParentId() *int64 {
	return nil
}

func (l *UnsLabel) GetAlias() string {
	return ""
}

func (l *UnsLabel) GetParentAlias() *string {
	return nil
}

func (l *UnsLabel) GetName() string {
	return l.LabelName
}

func (l *UnsLabel) GetDisplayName() string {
	return l.LabelName
}

func (l *UnsLabel) GetPath() string {
	return "label/" + l.LabelName
}

func (l *UnsLabel) GetDataType() *int16 {
	return nil
}
func (l *UnsLabel) GetParentDataType() *int16 {
	return nil
}
func (l *UnsLabel) GetPathType() int16 {
	return constants.PathTypeLabel
}

func (l *UnsLabel) GetMountType() *int16 {
	return nil
}

func (l *UnsLabel) GetMountSource() string {
	return ""
}

const TableNameUnsLabelRef = "uns_label_ref"

// UnsLabelRef mapped from table <uns_label_ref>
type UnsLabelRef struct {
	LabelID  int64     `gorm:"column:label_id;primaryKey" json:"label_id"`
	UnsID    int64     `gorm:"column:uns_id;primaryKey" json:"uns_id"`
	CreateAt time.Time `gorm:"column:create_at;default:now()" json:"create_at"`
}

// TableName UnsLabelRef's table name
func (*UnsLabelRef) TableName() string {
	return TableNameUnsLabelRef
}

const TableNameUnsAttachment = "uns_attachment"

// UnsAttachment mapped from table <uns_attachment>
type UnsAttachment struct {
	ID             int64     `gorm:"column:id;primaryKey" json:"id"`
	UnsAlias       string    `gorm:"column:uns_alias;not null" json:"uns_alias"`
	OriginalName   string    `gorm:"column:original_name;not null" json:"original_name"`
	AttachmentName string    `gorm:"column:attachment_name;not null" json:"attachment_name"`
	AttachmentPath string    `gorm:"column:attachment_path" json:"attachment_path"`
	ExtensionName  string    `gorm:"column:extension_name" json:"extension_name"`
	CreateAt       time.Time `gorm:"column:create_at;default:now()" json:"create_at"`
}

// TableName UnsAttachment's table name
func (*UnsAttachment) TableName() string {
	return TableNameUnsAttachment
}

const TableNameUnsAlarmsDatum = "uns_alarms_data"

// UnsAlarmsDatum mapped from table <uns_alarms_data>
type UnsAlarmsDatum struct {
	Ct           time.Time `gorm:"column:_ct;default:now()" json:"_ct"`
	ID           int64     `gorm:"column:_id;primaryKey;autoIncrement:true" json:"_id"`
	CurrentValue float32   `gorm:"column:current_value" json:"current_value"`
	IsAlarm      bool      `gorm:"column:is_alarm;default:true" json:"is_alarm"`
	LimitValue   float32   `gorm:"column:limit_value" json:"limit_value"`
	ReadStatus   bool      `gorm:"column:read_status" json:"read_status"`
	Uns          int64     `gorm:"column:uns" json:"uns"`
	UnsPath      string    `gorm:"column:uns_path" json:"uns_path"`
}

// TableName UnsAlarmsDatum's table name
func (*UnsAlarmsDatum) TableName() string {
	return TableNameUnsAlarmsDatum
}

const TableNameUnsAlarmsHandler = "uns_alarms_handler"

// UnsAlarmsHandler mapped from table <uns_alarms_handler>
type UnsAlarmsHandler struct {
	ID       int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UnsID    int64     `gorm:"column:uns_id" json:"uns_id"`
	UserID   string    `gorm:"column:user_id" json:"user_id"`
	Username string    `gorm:"column:username" json:"username"`
	CreateAt time.Time `gorm:"column:create_at;default:now()" json:"create_at"`
}

// TableName UnsAlarmsHandler's table name
func (*UnsAlarmsHandler) TableName() string {
	return TableNameUnsAlarmsHandler
}

const TableNameUnsExportRecord = "uns_export_record"

// UnsExportRecord mapped from table <uns_export_record>
type UnsExportRecord struct {
	ID         int64     `gorm:"column:id;primaryKey" json:"id"`
	UserID     string    `gorm:"column:user_id" json:"user_id"`
	FilePath   string    `gorm:"column:file_path" json:"file_path"`
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	Confirm    bool      `gorm:"column:confirm" json:"confirm"`
}

// TableName UnsExportRecord's table name
func (*UnsExportRecord) TableName() string {
	return TableNameUnsExportRecord
}

const TableNameUnsHistoryDeleteJob = "uns_history_delete_job"

// UnsHistoryDeleteJob mapped from table <uns_history_delete_job>
type UnsHistoryDeleteJob struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Alias_     string    `gorm:"column:alias;not null" json:"alias"`
	Name       string    `gorm:"column:name;not null" json:"name"`
	TableName_ string    `gorm:"column:table_name" json:"table_name"`
	Path       string    `gorm:"column:path;not null" json:"path"`
	PathType   int16     `gorm:"column:path_type;not null" json:"path_type"`
	DataType   int16     `gorm:"column:data_type" json:"data_type"`
	Fields     string    `gorm:"column:fields" json:"fields"`
	Status     int16     `gorm:"column:status;default:1" json:"status"`
	CreateAt   time.Time `gorm:"column:create_at;default:now()" json:"create_at"`
}

// TableName UnsHistoryDeleteJob's table name
func (*UnsHistoryDeleteJob) TableName() string {
	return TableNameUnsHistoryDeleteJob
}

const TableNameUnsMountExtend = "uns_mount_extend"

// UnsMountExtend mapped from table <uns_mount_extend>
type UnsMountExtend struct {
	ID                int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	SourceSubType     string `gorm:"column:source_sub_type" json:"source_sub_type"`
	MountSeq          string `gorm:"column:mount_seq" json:"mount_seq"`
	TargetAlias       string `gorm:"column:target_alias" json:"target_alias"`
	FirstSourceAlias  string `gorm:"column:first_source_alias" json:"first_source_alias"`
	SecondSourceAlias string `gorm:"column:second_source_alias" json:"second_source_alias"`
	SourceName        string `gorm:"column:source_name" json:"source_name"`
	Extend            string `gorm:"column:extend" json:"extend"`
}

// TableName UnsMountExtend's table name
func (*UnsMountExtend) TableName() string {
	return TableNameUnsMountExtend
}

const TableNameUnsMount = "uns_mount"

// UnsMount mapped from table <uns_mount>
type UnsMount struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	MountSeq    string `gorm:"column:mount_seq" json:"mount_seq"`
	TargetType  string `gorm:"column:target_type" json:"target_type"`
	TargetAlias string `gorm:"column:target_alias" json:"target_alias"`
	MountModel  string `gorm:"column:mount_model" json:"mount_model"`
	SourceAlias string `gorm:"column:source_alias" json:"source_alias"`
	MountStatus int16  `gorm:"column:mount_status" json:"mount_status"`
	Status      string `gorm:"column:status" json:"status"`
	DataType    int16  `gorm:"column:data_type" json:"data_type"`
	WithFlags   int32  `gorm:"column:with_flags" json:"with_flags"`
	Version     string `gorm:"column:version" json:"version"`
	NextVersion string `gorm:"column:next_version" json:"next_version"`
}

// TableName UnsMount's table name
func (*UnsMount) TableName() string {
	return TableNameUnsMount
}

const TableNameUnsPersonConfig = "uns_person_config"

// UnsPersonConfig mapped from table <uns_person_config>
type UnsPersonConfig struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID       string    `gorm:"column:user_id;not null" json:"user_id"`
	MainLanguage string    `gorm:"column:main_language" json:"main_language"`
	CreateAt     time.Time `gorm:"column:create_at;default:now()" json:"create_at"`
	UpdateAt     time.Time `gorm:"column:update_at;default:now()" json:"update_at"`
}

// TableName UnsPersonConfig's table name
func (*UnsPersonConfig) TableName() string {
	return TableNameUnsPersonConfig
}

const TableNameUnsSysCode = "uns_sys_code"

// UnsSysCode mapped from table <uns_sys_code>
type UnsSysCode struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ModuleCode  string    `gorm:"column:module_code;not null" json:"module_code"`
	EntityCode  string    `gorm:"column:entity_code;not null" json:"entity_code"`
	Code        string    `gorm:"column:code;not null" json:"code"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	Sort        float64   `gorm:"column:sort;default:1000" json:"sort"`
	DesA        string    `gorm:"column:des_a" json:"des_a"`
	DesB        string    `gorm:"column:des_b" json:"des_b"`
	DesC        string    `gorm:"column:des_c" json:"des_c"`
	CreateTime  time.Time `gorm:"column:create_time;default:now()" json:"create_time"`
	Description string    `gorm:"column:description" json:"description"`
}

// TableName UnsSysCode's table name
func (*UnsSysCode) TableName() string {
	return TableNameUnsSysCode
}

const TableNameUnsSysEntity = "uns_sys_entity"

// UnsSysEntity mapped from table <uns_sys_entity>
type UnsSysEntity struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ModuleCode  string    `gorm:"column:module_code;not null" json:"module_code"`
	Code        string    `gorm:"column:code;not null" json:"code"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	CreateTime  time.Time `gorm:"column:create_time;default:now()" json:"create_time"`
	SysDefault  bool      `gorm:"column:sys_default;default:true" json:"sys_default"`
}

// TableName UnsSysEntity's table name
func (*UnsSysEntity) TableName() string {
	return TableNameUnsSysEntity
}

const TableNameUnsSysModule = "uns_sys_module"

// UnsSysModule mapped from table <uns_sys_module>
type UnsSysModule struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Code       string    `gorm:"column:code;not null" json:"code"`
	Name       string    `gorm:"column:name;not null" json:"name"`
	Type       string    `gorm:"column:type;not null" json:"type"`
	CreateTime time.Time `gorm:"column:create_time;default:now()" json:"create_time"`
}

// TableName UnsSysModule's table name
func (*UnsSysModule) TableName() string {
	return TableNameUnsSysModule
}

const TableNameUnsTag = "uns_tag"

// UnsTag mapped from table <uns_tag>
type UnsTag struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	Topic     string    `gorm:"column:topic" json:"topic"`
	TagName   string    `gorm:"column:tag_name" json:"tag_name"`
	IsDeleted bool      `gorm:"column:is_deleted" json:"is_deleted"`
	CreateAt  time.Time `gorm:"column:create_at;default:now()" json:"create_at"`
}

// TableName UnsTag's table name
func (*UnsTag) TableName() string {
	return TableNameUnsTag
}

const TableNameUnsWebhookAction = "uns_webhook_action"

// UnsWebhookAction mapped from table <uns_webhook_action>
type UnsWebhookAction struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name       string    `gorm:"column:name;not null" json:"name"`
	Service    string    `gorm:"column:service;not null" json:"service"`
	EventMeta  string    `gorm:"column:event_meta;not null" json:"event_meta"`
	Action     string    `gorm:"column:action;not null" json:"action"`
	MethodType string    `gorm:"column:method_type;not null" json:"method_type"`
	Timeout    int32     `gorm:"column:timeout;not null" json:"timeout"`
	URL        string    `gorm:"column:url;not null" json:"url"`
	Headers    string    `gorm:"column:headers" json:"headers"`
	Params     string    `gorm:"column:params" json:"params"`
	Payload    string    `gorm:"column:payload" json:"payload"`
	Status     string    `gorm:"column:status" json:"status"`
	Message    string    `gorm:"column:message" json:"message"`
	CreateTime time.Time `gorm:"column:create_time;default:now()" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	Body       string    `gorm:"column:body" json:"body"`
}

// TableName UnsWebhookAction's table name
func (*UnsWebhookAction) TableName() string {
	return TableNameUnsWebhookAction
}

const TableNameUnsWebhook = "uns_webhook"

// UnsWebhook mapped from table <uns_webhook>
type UnsWebhook struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	Service     string    `gorm:"column:service;not null" json:"service"`
	EventMeta   string    `gorm:"column:event_meta;not null" json:"event_meta"`
	Actions     string    `gorm:"column:actions;not null" json:"actions"`
	MethodType  string    `gorm:"column:method_type;not null" json:"method_type"`
	URL         string    `gorm:"column:url;not null" json:"url"`
	Headers     string    `gorm:"column:headers" json:"headers"`
	Params      string    `gorm:"column:params" json:"params"`
	Description string    `gorm:"column:description" json:"description"`
	CreateTime  time.Time `gorm:"column:create_time;default:now()" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`
	Timeout     int32     `gorm:"column:timeout;not null;default:1000" json:"timeout"`
	Body        string    `gorm:"column:body" json:"body"`
}

// TableName UnsWebhook's table name
func (*UnsWebhook) TableName() string {
	return TableNameUnsWebhook
}
