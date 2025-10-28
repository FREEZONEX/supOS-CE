package dto

import (
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/types"
	"strings"
)

// UpdateUnsDto represents UNS update DTO
type UpdateUnsDto struct {
	ID          int64  `json:"id,omitzero"`
	Name        string `json:"name,omitzero"`
	DisplayName string `json:"displayName,omitzero"`
	PathType    int    `json:"pathType,omitzero" validate:"omitzero,min=0,max=2"`
	Path        string `json:"path,omitzero"`
	Alias       string `json:"alias" validate:"required"`
	Description string `json:"description,omitzero"`

	// Model/Template fields
	ModelID    int64  `json:"modelId,omitzero"`
	ModelAlias string `json:"modelAlias,omitzero"`
	Template   string `json:"-"`

	// Parent fields
	ParentAlias string `json:"parentAlias,omitzero"`
	ParentID    int64  `json:"parentId,omitzero"`

	// Data type and fields
	DataType  int                  `json:"dataType,omitzero" validate:"omitzero,min=1,max=7"`
	DataSrcID any                  `json:"-"`
	Fields    []*types.FieldDefine `json:"fields,omitzero"`

	// Table fields
	TableName string `json:"-"`
	DataPath  string `json:"-"`

	// Reference fields
	ReferUns     string               `json:"referUns,omitzero"`
	ReferIDs     []int64              `json:"referIds,omitzero"`
	ReferTable   string               `json:"referTable,omitzero"`
	RefFields    []*types.FieldDefine `json:"refFields,omitzero"`
	ReferModelID string               `json:"referModelId,omitzero"`
	Cited        map[int64]bool       `json:"-"`

	// Calculation fields
	Refers            []*types.InstanceField `json:"refers,omitzero"`
	Expression        string                 `json:"expression,omitzero" validate:"max=255"`
	CompileExpression any                    `json:"-"`
	StreamOptions     *StreamOptions         `json:"streamOptions,omitzero"`

	// Protocol fields
	Protocol     map[string]any `json:"-"`
	ProtocolType string         `json:"-"`
	ProtocolBean any            `json:"-"`

	// Flags and options
	Flags                         int  `json:"-"`
	AddFlow                       bool `json:"addFlow,omitzero"`
	AddDashBoard                  int  `json:"addDashBoard,omitzero"`
	Save2DB                       int  `json:"save2db,omitzero"`
	RetainTableWhenDeleteInstance int  `json:"-"`

	// Frequency for merge type
	Frequency        string `json:"frequency,omitzero"`
	FrequencySeconds int64  `json:"-"`

	// Alarm rule
	AlarmRuleDefine any `json:"-"`

	// Extended fields
	Extend     map[string]any `json:"extend,omitzero"`
	LabelNames []string       `json:"labelNames,omitzero"`
	Order      int            `json:"-"`

	// Pride specific fields
	RefSource string `json:"refSource,omitzero"`
	ValueType string `json:"valueType,omitzero"`
	InitValue any    `json:"initValue,omitzero"`
	StrMaxLen int    `json:"strMaxLen,omitzero"`

	// Access level
	AccessLevel string `json:"accessLevel,omitzero"`

	// Internal fields
	RefTopicFields map[int64]map[string]bool `json:"-"`
}

// GetTopic returns the topic for UpdateUnsDto
func (u *UpdateUnsDto) GetTopic() string {
	if constants.UseAliasAsTopic {
		return u.Alias
	}
	return u.Path
}

// GetTable returns the table name for UpdateUnsDto
func (u *UpdateUnsDto) GetTable() string {
	if u.TableName != "" {
		return u.TableName
	}
	if u.Alias != "" {
		return u.Alias
	}
	return u.Path
}

// SetPath sets and trims the path for UpdateUnsDto
func (u *UpdateUnsDto) SetPath(path string) {
	u.Path = strings.TrimSpace(path)
}

// SetAlias sets and trims the alias for UpdateUnsDto
func (u *UpdateUnsDto) SetAlias(alias string) {
	u.Alias = strings.TrimSpace(alias)
}

// SetFrequency sets frequency and calculates frequency in seconds for UpdateUnsDto
func (u *UpdateUnsDto) SetFrequency(frequency string) {
	u.Frequency = frequency
	if frequency != "" {
		frequency = strings.TrimSpace(frequency)
		nano, ok := enums.TimeUnitsParseToNanoSecond(frequency)
		if ok {
			seconds := nano / enums.TimeUnitSecond.Multiple
			u.FrequencySeconds = seconds
		}
	}
}

// SetDataPath sets data path for UpdateUnsDto
func (u *UpdateUnsDto) SetDataPath(dataPath string) *UpdateUnsDto {
	u.DataPath = dataPath
	return u
}

// SetCalculation sets calculation parameters for UpdateUnsDto
func (u *UpdateUnsDto) SetCalculation(refers []*types.InstanceField, expression string) *UpdateUnsDto {
	u.Refers = refers
	u.Expression = expression
	return u
}

// SetStreamCalculation sets stream calculation parameters for UpdateUnsDto
func (u *UpdateUnsDto) SetStreamCalculation(referTopic string, streamOptions *StreamOptions) *UpdateUnsDto {
	u.ReferUns = referTopic
	u.StreamOptions = streamOptions
	return u
}

// CountNumberFields counts number of numeric fields for UpdateUnsDto
func (u *UpdateUnsDto) CountNumberFields() int {
	if u.Fields == nil {
		return 0
	}
	count := 0
	for _, f := range u.Fields {
		if types.FieldType(f.Type).IsNumber() && !f.IsSystemField() {
			count++
		}
	}
	return count
}

// FilterBlobField filters BLOB and LBLOB field names for UpdateUnsDto
func (u *UpdateUnsDto) FilterBlobField() []string {
	if len(u.Fields) == 0 {
		return []string{}
	}
	result := make([]string, 0)
	for _, f := range u.Fields {
		if f.Type == types.FieldTypeBlob || f.Type == types.FieldTypeLBlob {
			result = append(result, f.Name)
		}
	}
	return result
}

// SimpleUnsInfo interface for basic UNS information
type SimpleUnsInfo interface {
	GetId() int64
	GetAlias() string
	GetName() string
	GetTableName() string
	GetPath() string
	GetDataType() *int16
	GetFields() []*types.FieldDefine
}

// SimpleUnsInstance represents a simple UNS instance
type SimpleUnsInstance struct {
	ID                            int64                `json:"id,omitzero"`
	Name                          string               `json:"name,omitzero"`
	Path                          string               `json:"path,omitzero"`
	Alias                         string               `json:"alias,omitzero"`
	TableName                     string               `json:"tableName,omitzero"`
	DataType                      int16                `json:"dataType,omitzero"`
	ParentID                      int64                `json:"parentId,omitzero"`
	Fields                        []*types.FieldDefine `json:"fields,omitzero"`
	RemoveDashboard               bool                 `json:"removeDashboard"`
	RemoveTableWhenDeleteInstance bool                 `json:"removeTableWhenDeleteInstance"`
	LabelIDs                      map[int64]bool       `json:"labelIds,omitzero"`
	Flags                         int                  `json:"flags,omitzero"`
}

// GetTopic returns the topic based on configuration
func (s *SimpleUnsInstance) GetTopic() string {
	if constants.UseAliasAsTopic {
		return s.Alias
	}
	return s.Path
}

// GetTableNameOnly returns the table name field value
func (s *SimpleUnsInstance) GetTableNameOnly() string {
	return s.TableName
}

// GetTableName returns the effective table name
func (s *SimpleUnsInstance) GetTableName() string {
	if s.TableName != "" {
		return s.TableName
	}
	if s.Alias != "" {
		return s.Alias
	}
	return s.Path
}

// Implement SimpleUnsInfo interface
func (s *SimpleUnsInstance) GetId() int64                    { return s.ID }
func (s *SimpleUnsInstance) GetAlias() string                { return s.Alias }
func (s *SimpleUnsInstance) GetName() string                 { return s.Name }
func (s *SimpleUnsInstance) GetPath() string                 { return s.Path }
func (s *SimpleUnsInstance) GetDataType() *int16             { return &s.DataType }
func (s *SimpleUnsInstance) GetFields() []*types.FieldDefine { return s.Fields }

// UnsSearchCondition represents UNS search conditions

// UnsTreeCondition represents UNS tree query conditions
type UnsTreeCondition struct {
	PaginationDTO

	SearchType int    `json:"searchType" form:"searchType"` // 查询类型：1-UNS（名称+别名） 2-含标签 3-含模板
	Keyword    string `json:"keyword,omitzero" form:"keyword"`
	ParentID   *int64 `json:"parentId,omitzero" form:"parentId"` // 父级ID  可为空，传0查询顶级节点，空值时查询所有
	DataType   *int   `json:"dataType,omitzero" form:"dataType"` // 数据类型：1--时序，2--关系，3--计算型, 5--告警 6--聚合 7--引用
	PathType   *int   `json:"pathType,omitzero" form:"pathType"` // 路径类型: 0--文件夹，2--文件
}

// SaveDataDto represents data saving DTO
type SaveDataDto struct {
	ID             int64                 `json:"id,omitzero"`
	Table          string                `json:"table,omitzero"`
	FieldDefines   *types.FieldDefines   `json:"-"`
	CreateTopicDto *types.CreateTopicDto `json:"-"`
	List           []map[string]any      `json:"list" validate:"required"`
	Tables         map[string]bool       `json:"-"`
}

// Clone creates a deep copy of SaveDataDto
func (s *SaveDataDto) Clone() *SaveDataDto {
	clone := &SaveDataDto{
		ID:             s.ID,
		Table:          s.Table,
		FieldDefines:   s.FieldDefines,
		CreateTopicDto: s.CreateTopicDto,
	}

	if s.List != nil {
		clone.List = make([]map[string]any, len(s.List))
		for i, item := range s.List {
			clonedItem := make(map[string]any)
			for k, v := range item {
				clonedItem[k] = v
			}
			clone.List[i] = clonedItem
		}
	}

	if s.Tables != nil {
		clone.Tables = make(map[string]bool)
		for k, v := range s.Tables {
			clone.Tables[k] = v
		}
	}

	return clone
}

// UpdateFileDTO represents file update DTO
type UpdateFileDTO struct {
	Alias string         `json:"alias,omitzero"`
	Data  map[string]any `json:"data,omitzero"`
}

type UnsResultDto struct {
	AllDefinitions []types.CreateTopicDto `json:"allDefinitions"`
	MatchResults   []types.CreateTopicDto `json:"matchResults"`
}

type UnsCountDTO struct {
	CountChildren       int  `json:"countChildren"`
	CountDirectChildren int  `json:"countDirectChildren"`
	HasChildren         bool `json:"hasChildren"`
}

// BatchRemoveUnsDto represents batch removal DTO
type BatchRemoveUnsDto struct {
	AliasList       []string `json:"aliasList" validate:"required,min=1"` // 别名集合
	WithFlow        bool     `json:"-"`                                   // 是否删除关联的Flow
	WithDashboard   bool     `json:"-"`                                   // 是否删除关联的Dashboard
	RemoveRefer     bool     `json:"-"`                                   // 是否删除引用
	CheckMount      bool     `json:"-"`                                   // 是否检查挂载
	OnlyRemoveChild bool     `json:"-"`                                   // 是否仅删除子节点
}
