// Package constants provides a centralized place for application-wide constants.
// This version uses idiomatic Go naming conventions (MixedCaps for exported identifiers).
package constants

import (
	"os"
	"regexp"
	"strconv"
	"sync/atomic"
)

// --- Simple Compile-time Constants ---
const ( // 注意：不能合并到第二个const()块，否则常量值可能会变
	// Data Types (using iota for sequential values)
	TimeSequenceType    int16 = iota + 1 // 时序类型 (1)
	RelationType                         // 关系类型 (2)
	CalculationRealType                  // 实时计算 (3)
	CalculationHistType                  // 历史值计算 (4)
	AlarmRuleType                        // 报警规则类型 (5)
	MergeType                            // 聚合类型 (6)
	CitingType                           // 引用类型 (7)
	JsonbType                            //整个json当做一个字段存储
)
const (
	// System Fields and Flags
	SystemFieldPrev = "_"
	SystemSeqTag    = "tag"
	SystemSeqValue  = "value"
	SysSaveTime     = SystemFieldPrev + "st"
	SysFieldID      = SystemFieldPrev + "id"
	MergeFlag       = "#mg#" // 按时间戳合并消息的标志
	FirstMsgFlag    = "#1#"  // 启动时首条消息的标志

	// Path Types
	PathTypeDir      = int16(0) // 目录
	PathTypeFile     = int16(2) // 文件
	PathTypeTemplate = int16(1) // 模板
	PathTypeLabel    = int16(7) // 标签

	// Topic and Message Data Keys
	ResultTopicPrev = "_rs/"       // 处理结果 topic前缀
	MsgRawDataKey   = "_source_"   // 原始数据的 json key
	MsgResDataKey   = "_resource_" // 处理过的 json key

	// Pagination Defaults
	DefaultPageSize = 20
	DefaultPageNum  = 1
	MaxPageSize     = 1000
	SQLBatchSize    = 200

	// homepage
	DefaultHomepage = "/uns"

	// Paths (composed at compile time)
	RootPath              = "/data"
	PluginPath            = RootPath + "/plugins"
	PluginTempPath        = PluginPath + "/temp"
	PluginInstalledPath   = PluginPath + "/installed"
	PluginUpgradePath     = PluginPath + "/upgrade"
	PluginUpgradeTempPath = PluginPath + "/upgrade-temp"
	PluginFrontendPath    = RootPath + "/plugins-frontend"
	LogPath               = "/logs"
	ImportErr             = "/import_err/"
	UploadRoot            = "/upload/"
	ExportRoot            = "/export/"

	GlobalImport        = "/global-import/"
	GlobalExport        = "/global-export/"
	GlobalImportError   = "/global-import-error/"
	ExampleRoot         = "/example/"
	UnsRoot             = "/uns/"
	SystemRoot          = "/system/"
	ExcelTemplatePath   = "/templates/all-namespace.xlsx"
	ExcelTemplateZHPath = "/templates/all-namespace-zh-CN.xlsx"
	JSONTemplatePath    = "/templates/all-namespace.json"
	ExcelOutPath        = "/export/all-namespace.xlsx"
	JSONOutPath         = "/export/all-namespace.json"
	ExcelOutTempPath    = "/export/temp-namespace.xlsx"
	I18nExcelOutPath    = "/export/i18n_languageCode.xlsx"
	BlobPath            = "/data/uns"
	AccessTokenKey      = "supos_community_token"
	CookieMaxAge        = 60 * 60 * 24 * 365

	// Regex Strings
	AliasReg = `[a-zA-Z_$][a-zA-Z0-9_$]*$`
	TopicReg = `^[\x{4e00}-\x{9fa5}a-zA-Z0-9/_-]+$`
	NameReg  = `^[\x{4e00}-\x{9fa5}a-zA-Z0-9_-]+$`

	VarPrev               = "a"
	DefaultRoleID         = "d12d7ca2-34e1-4f26-9a03-6b4f7f411567"
	AttachmentMaxSize     = 10 * 1024 * 1024
	TDJDBCURL             = "tdengine:6041"
	PGJDBCURL             = "postgresql:5432"
	TSDBJDBCURL           = "tsdb:5432"
	AuthCheckKongPluginID = "1845ee75-d704-40e1-a8b0-aa2baaf9d71b"
	EmqxAPIKey            = "b441dbabd9bd5c26"
	EmqxSecretKey         = "59CdRlRvDaygamiil6789A2JvbXfO9ADRcLEcgxB9CYVv5Y"
	UnknownUser           = "Unknown User"
	ExampleFuxaFile       = "fuxa-"
	ExampleGrafana        = "grafana-"
	ExampleMetadata       = "metadata.json"
	ExampleProtocol       = "protocol.json"
	FuxaAPIURL            = "http://fuxa:1881"

	// Command Types
	CmdSub     = 1 // 实时数据订阅
	CmdSubRes  = 2 // 订阅响应
	CmdValPush = 3 // 实时值推送

	ValuePushBatchSize = 500

	PrideTemplatePrefix      = "system."
	PrideCollectorFilePrefix = "C"

	// Access Levels
	AccessLevelReadOnly  = "READ_ONLY"
	AccessLevelReadWrite = "READ_WRITE"

	// Config files and export names
	NotifyConfigFile       = "notify_config.yml"
	GlobalExportYAML       = "export.yml"
	GlobalExportSourceFlow = "source_flow.json"
	GlobalExportEventFlow  = "event_flow.json"
	GlobalExportDashboard  = "dashboard.json"

	PlugPrefix = "PLUGIN"
	AppPrefix  = "APP"

	// nodered
	FlowTypeNODERED   = "node-red"
	FlowTypeEVENTFLOW = "event-flow"
)

const (
	_                   = iota // 0 is unused
	JDBCTypePostgres           // 1
	JDBCTypeMariaDB            // 2
	JDBCTypeTimeScaleDB        // 3
)

// UnsFlags are bitmask constants for UNS features.
const (
	UnsFlagWithFlow                   = 1 << iota // 是否添加数据采集流程
	UnsFlagWithDashboard                          // 是否添加数据看板
	UnsFlagWithSave2DB                            // 是否持久化到数据库
	UnsFlagRetainTableWhenDelInstance             // 刪除实例时是否保留数据表
	UnsFlagAlarmAcceptPerson                      // 报警规则接收方式 16人员
	UnsFlagAlarmAcceptWorkflow                    // 报警规则接收方式 32工作流
	UnsFlagAccessLevelReadOnly                    // 北向访问级别:READ_ONLY-只读
	UnsFlagAccessLevelReadWrite                   // 北向访问级别:READ_WRITE-读写
	UnsFlagWithSubscribeEnable                    // Skip 1 << 8
	UnsFlagWithAttachment                         // UNS带附件的标志
	UnsFlagHasData                                // UNS有存过数据的标志
)

// --- Runtime-initialized Variables and Compiled Regex ---
var (
	readOnlyMode atomic.Bool // Equivalent to AtomicBoolean

	UseAliasAsTopic    bool // 是否使用别名alias作为 mqtt topic
	MqttPlugin         string
	SysFieldCreateTime string
	QosField           string
	UnsAddBatchSize    int
	WSSessionLimit     int // ws会话限制
	UnsOverdueDelete   int
	OSVersion          string
	TokenMaxAge        int // token失效时间（秒）

	// SystemFields is a set of system field names for quick lookups.
	SystemFields map[string]struct{}

	// Compiled Regex Patterns
	AliasPattern = regexp.MustCompile(AliasReg)
	NamePattern  = regexp.MustCompile(NameReg)
)

func init() {
	readOnlyMode.Store(true)

	UseAliasAsTopic = getEnvAsBool("SYS_OS_USE_ALIAS_PATH_AS_TOPIC", false)
	OSVersion = getEnv("SYS_OS_VERSION", "1.0")
	SysFieldCreateTime = getEnv("SYS_OS_TIMESTAMP_NAME", "timeStamp")
	QosField = getEnv("SYS_OS_QUALITY_NAME", "status")
	UnsAddBatchSize = getEnvAsInt("UNS_ADD_BATCH_SIZE", 1000)
	MqttPlugin = getEnv("MQTT_PLUGIN", "emqx")
	WSSessionLimit = getEnvAsInt("WS_SESSION_LIMIT", 50)
	UnsOverdueDelete = getEnvAsInt("UNS_HISTORY_OVER_DUE", 7)
	TokenMaxAge = getEnvAsInt("TOKEN_MAX_AGE", 3600)

	// Initialize the set of system fields
	SystemFields = map[string]struct{}{
		SystemSeqTag:       {},
		SysFieldID:         {},
		SysFieldCreateTime: {}, // Note: Using the runtime value here
		QosField:           {}, // Note: Using the runtime value here
		SysSaveTime:        {},
		"_ct":              {},
	}
}

// --- Helper functions for initialization ---

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default.
func getEnvAsInt(key string, defaultValue int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultValue
}

// getEnvAsBool retrieves an environment variable as a boolean or returns a default.
func getEnvAsBool(key string, defaultValue bool) bool {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultValue
}

// --- Helper Functions (equivalent to public static methods) ---

func WithFlow(unsFlag int32) bool {
	return (unsFlag & UnsFlagWithFlow) == UnsFlagWithFlow
}

func WithDashBoard(unsFlag int32) bool {
	return (unsFlag & UnsFlagWithDashboard) == UnsFlagWithDashboard
}

func WithSave2db(unsFlag int32) bool {
	return (unsFlag & UnsFlagWithSave2DB) == UnsFlagWithSave2DB
}

func WithRetainTableWhenDeleteInstance(unsFlag int32) bool {
	return (unsFlag & UnsFlagRetainTableWhenDelInstance) == UnsFlagRetainTableWhenDelInstance
}

func WithReadOnly(unsFlag int32) string {
	if (unsFlag & UnsFlagAccessLevelReadOnly) == UnsFlagAccessLevelReadOnly {
		return AccessLevelReadOnly
	}
	if (unsFlag & UnsFlagAccessLevelReadWrite) == UnsFlagAccessLevelReadWrite {
		return AccessLevelReadWrite
	}
	return "" // Go uses "" instead of null for strings
}

func WithAttachment(flags *int32) bool {
	if flags == nil {
		return false
	}
	return (*flags & UnsFlagWithAttachment) == UnsFlagWithAttachment
}

func WithHasData(flags *int32) bool {
	if flags == nil {
		return false
	}
	return (*flags & UnsFlagHasData) == UnsFlagHasData
}
func WithSubscribeEnable(flags int32) bool {
	return (flags & UnsFlagWithSubscribeEnable) == UnsFlagWithSubscribeEnable
}

func IsValidDataType(dataType int16) bool {
	return dataType >= TimeSequenceType && dataType <= JsonbType
}

func GenerateCodeByPrefix(prefix, code string) string {
	return prefix + "#" + code
}
