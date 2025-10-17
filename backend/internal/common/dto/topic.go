package dto

import (
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"fmt"
	"strings"
	"time"
)

// CreateTopicDto represents UNS topic creation/update DTO
type CreateTopicDto struct {
	// Basic fields
	ID          int64  `json:"id,omitzero"`
	Batch       int    `json:"-"`
	Index       int    `json:"-"`
	FlagNo      string `json:"-"`
	Name        string `json:"name" validate:"required,max=63"`
	DisplayName string `json:"displayName,omitzero" validate:"max=128"`
	PathType    int    `json:"pathType" validate:"required,min=0,max=2"`
	Path        string `json:"path,omitzero"`
	Alias       string `json:"alias" validate:"required"`
	Description string `json:"description,omitzero" validate:"max=255"`

	// Model/Template fields
	ModelID    int64  `json:"modelId,omitzero"`
	ModelAlias string `json:"modelAlias,omitzero"`
	Template   string `json:"-"`

	// Parent fields
	ParentAlias string `json:"parentAlias,omitzero"`
	ParentID    int64  `json:"parentId,omitzero"`

	// Data type and fields
	DataType  int            `json:"dataType" validate:"required,min=1,max=7"`
	Fields    []*FieldDefine `json:"fields,omitzero"`
	DataSrcID any            `json:"-"` // SrcJdbcType

	// Table fields
	TableName    string   `json:"-"`
	TbFieldName  string   `json:"-"`
	PrimaryField []string `json:"-"`
	HasBlobField bool     `json:"-"`

	// Reference fields
	ReferUns     string         `json:"referUns,omitzero"`
	ReferIDs     []int64        `json:"referIds,omitzero"`
	RefUns       map[int64]int  `json:"-"`
	LayRec       string         `json:"-"`
	ReferTable   string         `json:"referTable,omitzero"`
	RefFields    []*FieldDefine `json:"refFields,omitzero"`
	ReferModelID string         `json:"referModelId,omitzero"`
	Cited        map[int64]bool `json:"-"` // Set of cited IDs

	// Calculation fields
	Refers            []*InstanceField `json:"refers,omitzero"`
	Expression        string           `json:"expression,omitzero" validate:"max=255"`
	CompileExpression any              `json:"-"`
	StreamOptions     *StreamOptions   `json:"streamOptions,omitzero"`

	// Protocol fields
	DataPath     string         `json:"-"`
	Protocol     map[string]any `json:"-"`
	ProtocolType string         `json:"-"`
	ProtocolBean any            `json:"-"`

	// Flags and options
	Flags                         int `json:"-"`
	AddFlow                       int `json:"addFlow,omitzero"`
	AddDashBoard                  int `json:"-"`
	Save2DB                       int `json:"save2db,omitzero"`
	RetainTableWhenDeleteInstance int `json:"-"`
	CreateTemplate                int `json:"-"`

	// Frequency for merge type
	Frequency        string `json:"frequency,omitzero"`
	FrequencySeconds int64  `json:"-"`

	// Alarm rule
	AlarmRuleDefine any `json:"-"` // AlarmRuleDefine type

	// Extended fields
	Extend          map[string]any   `json:"extend,omitzero" validate:"max=3"`
	ExtendFieldUsed []string         `json:"extendFieldUsed,omitzero"`
	LabelNames      []string         `json:"labelNames,omitzero"`
	LabelIDs        map[int64]string `json:"-"`
	Order           int              `json:"-"`

	// Pride specific fields
	RefSource string `json:"refSource,omitzero"`
	ValueType string `json:"valueType,omitzero"`
	InitValue any    `json:"initValue,omitzero"`
	StrMaxLen int    `json:"strMaxLen,omitzero"`

	// Access level
	AccessLevel string `json:"accessLevel,omitzero"`

	// Mount fields
	MountType   int    `json:"mountType,omitzero"`
	MountSource string `json:"mountSource,omitzero"`

	// Update metadata
	UpdateAt      time.Time `json:"updateAt,omitzero"`
	FieldsChanged bool      `json:"-"`

	// Internal fields
	tmField        string                    `json:"-"`
	fieldDefines   *FieldDefines             `json:"-"`
	RefTopicFields map[int64]map[string]bool `json:"-"`
}

// GetTopic returns the topic based on configuration
func (c *CreateTopicDto) GetTopic() string {
	if constants.UseAliasAsTopic {
		return c.Alias
	}
	return c.Path
}

// GetTimestampField returns the timestamp field name
func (c *CreateTopicDto) GetTimestampField() string {
	if c.tmField != "" {
		return c.tmField
	}

	if len(c.Fields) > 0 {
		// Find timestamp field (implementation depends on FieldUtils)
		for _, f := range c.Fields {
			if f.Name == constants.SysFieldCreateTime || f.Name == "timestamp" {
				c.tmField = f.Name
				return c.tmField
			}
		}
	}

	return ""
}

// GetQualityField returns the quality field name
func (c *CreateTopicDto) GetQualityField() string {
	if len(c.Fields) > 2 && c.DataSrcID != nil {
		// Find quality field (implementation depends on FieldUtils and dataSrcId.typeCode)
		for _, f := range c.Fields {
			if f.Name == constants.QosField || f.Name == "quality" {
				return f.Name
			}
		}
	}
	return ""
}

// GetTable returns the table name
func (c *CreateTopicDto) GetTable() string {
	if c.TableName != "" {
		return c.TableName
	}
	if c.Alias != "" {
		return c.Alias
	}
	return c.Path
}

// SetTableName sets the table name and parses field name if present
func (c *CreateTopicDto) SetTableName(table string) {
	if table == "" {
		c.TableName = ""
		return
	}

	c.TableName = table

	// Parse format: database.table.field
	parts := strings.Split(table, ".")
	if len(parts) == 3 {
		c.TableName = strings.TrimSpace(parts[0]) + "." + strings.TrimSpace(parts[1])
		c.TbFieldName = strings.TrimSpace(parts[2])
	}
}

// SetPath sets and trims the path
func (c *CreateTopicDto) SetPath(path string) {
	c.Path = strings.TrimSpace(path)
}

// SetAlias sets and trims the alias
func (c *CreateTopicDto) SetAlias(alias string) {
	c.Alias = strings.TrimSpace(alias)
}

// SetFields sets fields and updates related metadata
func (c *CreateTopicDto) SetFields(fields []*FieldDefine) {
	c.Fields = fields
	c.fieldDefines = NewFieldDefines(fields)
	c.HasBlobField = len(c.FilterAllBlobField()) > 0

	if len(fields) > 0 {
		pkSet := make(map[string]bool)
		for _, f := range fields {
			if f.IsUnique() {
				pkSet[f.Name] = true
			}
			if f.TbValueName != "" {
				c.TbFieldName = f.Name
			}
		}

		if len(pkSet) > 0 {
			c.PrimaryField = make([]string, 0, len(pkSet))
			for k := range pkSet {
				c.PrimaryField = append(c.PrimaryField, k)
			}
		}
	} else {
		c.PrimaryField = nil
	}
}

// SetFrequency sets frequency and calculates frequency in seconds
func (c *CreateTopicDto) SetFrequency(frequency string) {
	c.Frequency = frequency

	if frequency != "" {
		frequency = strings.TrimSpace(frequency)
		nano, ok := enums.TimeUnitsParseToNanoSecond(frequency)
		if ok {
			seconds := nano / enums.TimeUnitSecond.Multiple
			c.FrequencySeconds = seconds
		}
	}
}

// SetExpression sets expression and clears compiled expression
func (c *CreateTopicDto) SetExpression(expression string) {
	c.Expression = expression
	c.CompileExpression = nil
}

// CountNumberFields counts number of numeric fields
func (c *CreateTopicDto) CountNumberFields() int {
	if c.Fields == nil {
		return 0
	}

	count := 0
	for _, f := range c.Fields {
		if f.Type.IsNumber() && !f.IsSystemField() {
			count++
		}
	}
	return count
}

// GainBatchIndex returns batch index string
func (c *CreateTopicDto) GainBatchIndex() string {
	if c.FlagNo != "" {
		return c.FlagNo
	}
	return fmt.Sprintf("%d-%d", c.Batch, c.Index)
}

// FilterAllBlobField filters all BLOB and LBLOB fields
func (c *CreateTopicDto) FilterAllBlobField() []*FieldDefine {
	if len(c.Fields) == 0 {
		return []*FieldDefine{}
	}

	result := make([]*FieldDefine, 0)
	for _, f := range c.Fields {
		if f.Type == enums.FieldTypeBlob || f.Type == enums.FieldTypeLBlob {
			result = append(result, f)
		}
	}
	return result
}

// FilterBlobField filters BLOB fields only
func (c *CreateTopicDto) FilterBlobField() []*FieldDefine {
	if len(c.Fields) == 0 {
		return []*FieldDefine{}
	}

	result := make([]*FieldDefine, 0)
	for _, f := range c.Fields {
		if f.Type == enums.FieldTypeBlob {
			result = append(result, f)
		}
	}
	return result
}

// SetCalculation sets calculation parameters
func (c *CreateTopicDto) SetCalculation(refers []*InstanceField, expression string) *CreateTopicDto {
	c.Refers = refers
	c.Expression = expression
	return c
}

// SetStreamCalculation sets stream calculation parameters
func (c *CreateTopicDto) SetStreamCalculation(referTopic string, streamOptions *StreamOptions) *CreateTopicDto {
	c.ReferUns = referTopic
	c.StreamOptions = streamOptions
	return c
}

// SetDataPath sets data path
func (c *CreateTopicDto) SetDataPath(dataPath string) *CreateTopicDto {
	c.DataPath = dataPath
	return c
}

func (c *CreateTopicDto) GetFieldDefines() *FieldDefines {
	return c.fieldDefines
}

// TimestampField is an alias for GetTimestampField
var TimestampField = (*CreateTopicDto).GetTimestampField

// QualityField is an alias for GetQualityField
var QualityField = (*CreateTopicDto).GetQualityField

// Label represents a label for file/folder tagging
type Label struct {
	ID        string    `json:"id,omitzero"`        // 标签ID：已有标签时必传，新建标签时不需要传
	LabelName string    `json:"labelName,omitzero"` // 标签名称，新建标签时，必传
	CreateAt  time.Time `json:"createAt,omitzero"`  // 创建时间
}

// CreateFileDto represents file creation DTO
type CreateFileDto struct {
	CreateTopicDto         // Embed CreateTopicDto
	LabelList      []Label `json:"labelList,omitzero"` // 标签列表
}

// NewCreateFileDto creates a new CreateFileDto with PATH_TYPE_FILE
func NewCreateFileDto() *CreateFileDto {
	pathType := constants.PathTypeFile
	return &CreateFileDto{
		CreateTopicDto: CreateTopicDto{
			PathType: pathType,
		},
	}
}

// CreateFolderDto represents folder creation DTO
type CreateFolderDto struct {
	CreateTopicDto // Embed CreateTopicDto
}

// NewCreateFolderDto creates a new CreateFolderDto with PATH_TYPE_DIR
func NewCreateFolderDto() *CreateFolderDto {
	pathType := constants.PathTypeDir
	return &CreateFolderDto{
		CreateTopicDto: CreateTopicDto{
			PathType: pathType,
		},
	}
}
