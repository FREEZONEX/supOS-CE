package types

import (
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/common/utils/FieldFlags"
	"backend/share/base"
	"fmt"
	"strings"
)

func (c *CreateTopicDto) GetFieldDefines() *FieldDefines {
	if c.FieldDefines == nil && len(c.Fields) > 0 {
		c.SetFields(c.Fields)
	}
	return c.FieldDefines
}

// GetId 获取 Id
func (c *CreateTopicDto) GetId() int64 {
	return c.Id
}

// GetName 获取 Name
func (c *CreateTopicDto) GetName() string {
	return c.Name
}

// GetDisplayName 获取 DisplayName
func (c *CreateTopicDto) GetDisplayName() string {
	if dn := c.DisplayName; dn != nil {
		return *dn
	}
	return ""
}
func (c *CreateTopicDto) GetTopic() string {
	if constants.UseAliasAsTopic {
		return c.Alias
	} else {
		return c.Path
	}
}

// GetPathType 获取 PathType
func (c *CreateTopicDto) GetPathType() int16 {
	return c.PathType
}

// GetReferIds 获取 ReferIds
func (c *CreateTopicDto) GetReferIds() []int64 {
	return c.ReferIds
}

// GetReferTable 获取 ReferTable
func (c *CreateTopicDto) GetReferTable() string {
	return c.ReferTable
}

// GetRefFields 获取 RefFields
func (c *CreateTopicDto) GetRefFields() []*FieldDefine {
	return c.RefFields
}

// GetReferModelId 获取 ReferModelId
func (c *CreateTopicDto) GetReferModelId() string {
	return c.ReferModelID
}

// GetAlias 获取 Alias
func (c *CreateTopicDto) GetAlias() string {
	return c.Alias
}

// GetModelId 获取 ModelId
func (c *CreateTopicDto) GetModelId() *int64 {
	return c.ModelId
}

// GetModelAlias 获取 ModelAlias
func (c *CreateTopicDto) GetModelAlias() *string {
	return c.ModelAlias
}

// GetParentAlias 获取 ParentAlias
func (c *CreateTopicDto) GetParentAlias() *string {
	return c.ParentAlias
}

// GetParentId 获取 ParentId
func (c *CreateTopicDto) GetParentId() *int64 {
	return c.ParentId
}

// GetDataType 获取 DataType
func (c *CreateTopicDto) GetDataType() *int16 {
	return c.DataType
}

// GetParentDataType 获取 ParentDataType
func (c *CreateTopicDto) GetParentDataType() *int16 {
	return c.ParentDataType
}

// GetFields 获取 Fields
func (c *CreateTopicDto) GetFields() []*FieldDefine {
	return c.Fields
}

// GetExtendFieldUsed 获取 ExtendFieldUsed
func (c *CreateTopicDto) GetExtendFieldUsed() []string {
	return c.ExtendFieldUsed
}

// GetDataPath 获取 DataPath
func (c *CreateTopicDto) GetDataPath() string {
	return base.P2v(c.DataPath)
}
func (c *CreateTopicDto) GetPath() string {
	return c.Path
}

// GetDescription 获取 Description
func (c *CreateTopicDto) GetDescription() string {
	return base.P2v(c.Description)
}

// GetProtocolType 获取 ProtocolType
func (c *CreateTopicDto) GetProtocolType() string {
	return base.P2v(c.ProtocolType)
}

// GetProtocol 获取 Protocol
func (c *CreateTopicDto) GetProtocol() map[string]interface{} {
	return c.Protocol
}

// GetRefers 获取 Refers
func (c *CreateTopicDto) GetRefers() []*InstanceField {
	return c.Refers
}

// GetExpression 获取 Expression
func (c *CreateTopicDto) GetExpression() string {
	return base.P2v(c.Expression)
}

// GetStreamOptions 获取 StreamOptions
func (c *CreateTopicDto) GetStreamOptions() *StreamOptions {
	return c.StreamOptions
}

// GetAddFlow 获取 AddFlow
func (c *CreateTopicDto) GetAddFlow() *bool {
	return c.AddFlow
}

// GetAddDashBoard 获取 AddDashBoard
func (c *CreateTopicDto) GetAddDashBoard() *bool {
	return c.AddDashBoard
}

// GetSave2db 获取 Save2db
func (c *CreateTopicDto) GetSave2db() *bool {
	return c.Save2Db
}

// GetRetainTableWhenDeleteInstance 获取 RetainTableWhenDeleteInstance
func (c *CreateTopicDto) GetRetainTableWhenDeleteInstance() *bool {
	return c.RetainTableWhenDeleteInstance
}

// GetCreateTemplate 获取 CreateTemplate
func (c *CreateTopicDto) GetCreateTemplate() *bool {
	return c.CreateTemplate
}

// GetFrequency 获取 Frequency
func (c *CreateTopicDto) GetFrequency() string {
	return c.Frequency
}

// GetExtend 获取 Extend
func (c *CreateTopicDto) GetExtend() map[string]interface{} {
	return c.Extend
}

// GetLabelNames 获取 LabelNames
func (c *CreateTopicDto) GetLabelNames() []string {
	return c.LabelNames
}

// GetRefSource 获取 RefSource
func (c *CreateTopicDto) GetRefSource() string {
	return c.RefSource
}

// GetValueType 获取 ValueType
func (c *CreateTopicDto) GetValueType() string {
	return c.ValueType
}

// GetInitValue 获取 InitValue
func (c *CreateTopicDto) GetInitValue() interface{} {
	return c.InitValue
}

// GetStrMaxLen 获取 StrMaxLen
func (c *CreateTopicDto) GetStrMaxLen() int {
	return c.StrMaxLen
}

// GetAccessLevel 获取 AccessLevel
func (c *CreateTopicDto) GetAccessLevel() string {
	return c.AccessLevel
}

// GetMountType 获取 MountType
func (c *CreateTopicDto) GetMountType() *int16 {
	return c.MountType
}

// GetMountSource 获取 MountSource
func (c *CreateTopicDto) GetMountSource() string {
	return base.P2v(c.MountSource)
}
func (c *CreateTopicDto) GetCreateAt() int64 {
	return c.UpdateAt
}

// GetUpdateAt 获取 UpdateAt
func (c *CreateTopicDto) GetUpdateAt() int64 {
	return c.UpdateAt
}
func (c *CreateTopicDto) GetPrimaryField() []string {
	if c.FieldDefines == nil {
		c.SetFields(c.Fields)
	}
	return c.PrimaryField
}

// GetTimestampField returns the timestamp field name
func (c *CreateTopicDto) GetTimestampField() string {
	if c.TmField != "" {
		return c.TmField
	}

	if len(c.Fields) > 0 {
		// Find timestamp field (implementation depends on FieldUtils)
		for _, f := range c.Fields {
			if f.Name == constants.SysFieldCreateTime || f.Name == "timestamp" {
				c.TmField = f.Name
				return c.TmField
			}
		}
	}

	return ""
}

// GetQualityField returns the quality field name
func (c *CreateTopicDto) GetQualityField() string {
	if len(c.Fields) > 2 && c.DataSrcID > 0 {
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
func (c *CreateTopicDto) GetTbFieldName() string {
	if c.FieldDefines == nil && len(c.Fields) > 0 {
		c.SetFields(c.Fields)
	}
	return c.TbFieldName
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
	c.FieldDefines = NewFieldDefines(fields)
	c.HasBlobField = len(c.FilterAllBlobField()) > 0

	if len(fields) > 0 {
		pkSet := make(map[string]bool)
		for _, f := range fields {
			if f.IsUnique() {
				pkSet[f.Name] = true
			}
			if f.TbValueName != nil {
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
			c.FrequencySeconds = &seconds
		}
	}
}

// SetExpression sets expression and clears compiled expression
func (c *CreateTopicDto) SetExpression(expression string) {
	if expression != "" {
		c.Expression = &expression
	}
	c.CompileExpression = nil
}

// CountNumberFields counts number of numeric fields
func (c *CreateTopicDto) CountNumberFields() int16 {
	if c.Fields == nil {
		return 0
	}

	count := int16(0)
	for _, f := range c.Fields {
		if FieldType(strings.ToUpper(f.Type)).IsNumber() && !f.IsSystemField() {
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
	c.FlagNo = fmt.Sprintf("%d-%d", c.Batch, c.Index)
	return c.FlagNo
}

// FilterAllBlobField filters all BLOB and LBLOB fields
func (c *CreateTopicDto) FilterAllBlobField() []*FieldDefine {
	if len(c.Fields) == 0 {
		return []*FieldDefine{}
	}

	result := make([]*FieldDefine, 0)
	for _, f := range c.Fields {
		if f.Type == FieldTypeBlob || f.Type == FieldTypeLBlob {
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
		if f.Type == FieldTypeBlob {
			result = append(result, f)
		}
	}
	return result
}

// SetCalculation sets calculation parameters
func (c *CreateTopicDto) SetCalculation(refers []*InstanceField, expression string) *CreateTopicDto {
	c.Refers = refers
	if len(expression) > 0 {
		c.Expression = base.V2p(expression)
	}
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
	if len(dataPath) > 0 {
		c.DataPath = &dataPath
	}
	return c
}

func (c *CreateTopicDto) GetLayRec() string {
	return c.LayRec
}

func (c *CreateTopicDto) GetProtocolMap() map[string]interface{} {
	return c.Protocol
}

func (c *CreateTopicDto) GetCalculationType() *int32 {
	return nil //TODO
}

//func (c *CreateTopicDto) GetReadWriteMode() string {
//	return c.readWriteMode
//}

func (c *CreateTopicDto) GetLabelIds() map[int64]string {
	return c.LabelIDs
}

func (c *CreateTopicDto) GetRefUns() map[int64]int {
	return c.RefUns
}

func (c *CreateTopicDto) GetFlags() *int32 {
	return c.WithFlags
}

func (c *CreateTopicDto) GetExtendFieldFlags() *int32 {
	flag := FieldFlags.GenerateFlag(c.ExtendFieldUsed)
	return &flag
}

func (c *CreateTopicDto) GetSrcJdbcType() SrcJdbcType {
	return SrcJdbcType(c.DataSrcID)
}

func (c *CreateTopicDto) GetStatus() *int16 {
	return &c.Status
}
