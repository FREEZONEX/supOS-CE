package excel

import (
	"backend/internal/common/utils/PathUtil"
	"backend/internal/types"
	"backend/share/base"
	"fmt"
	"strings"
)

// ExcelUnsWrapDto is a wrapper for UNS data during Excel processing.
type ExcelUnsWrapDto struct {
	CheckSuccess  bool `json:"checkSuccess"`
	Uns           types.CreateTopicDto
	TemplateAlias string                `json:"templateAlias"`
	Labels        []string              `json:"labels"`
	Refers        []types.InstanceField `json:"refers"`
}

func NewExcelUnsWrapDto(uns types.CreateTopicDto) *ExcelUnsWrapDto {
	return &ExcelUnsWrapDto{
		CheckSuccess: true,
		Uns:          uns,
	}
}

func (e *ExcelUnsWrapDto) GetFlagNo() string {
	return e.Uns.GainBatchIndex()
}

func (e *ExcelUnsWrapDto) AddLabel(label string) {
	e.Labels = append(e.Labels, label)
}

// ExcelTemplateDto
// TODO: Custom validation for name needs to be implemented in the service layer.
type ExcelTemplateDto struct {
	Batch       int    `json:"batch"`
	Index       int    `json:"index"`
	Name        string `json:"name" validate:"required,max=63"`
	Alias       string `json:"alias"`
	Fields      string `json:"fields"` // JSON string of FieldDefine array
	Description string `json:"description"`
}

func (e *ExcelTemplateDto) GainBatchIndex() string {
	return fmt.Sprintf("%d-%d", e.Batch, e.Index)
}

// ExcelNameSpaceDto
// TODO: Custom validation for path (@TopicNameValidator) and alias (@AliasValidator) needs to be implemented in the service layer.
type ExcelNameSpaceDto struct {
	Batch         int     `json:"batch"`
	Index         int     `json:"index"`
	Path          string  `json:"path" validate:"required"`
	Alias         string  `json:"alias"`
	DisplayName   string  `json:"displayName"`
	TemplateAlias string  `json:"templateAlias"`
	Fields        string  `json:"fields"` // JSON string of FieldDefine array
	Description   *string `json:"description"`
	Refers        string  `json:"refers"`
	Expression    string  `json:"expression"`
}

func (e *ExcelNameSpaceDto) CreateTopic() *types.CreateTopicDto {
	return &types.CreateTopicDto{
		Index:       e.Index,
		Batch:       e.Batch,
		Path:        e.Path,
		Name:        PathUtil.GetName(e.Path),
		Alias:       e.Alias,
		DisplayName: &e.DisplayName,
		Description: e.Description,
	}
}

func (e *ExcelNameSpaceDto) GainBatchIndex() string {
	return fmt.Sprintf("%d-%d", e.Batch, e.Index)
}

// ExcelFolderDto
// TODO: Custom validation for path (@TopicNameValidator) and alias (@AliasValidator) needs to be implemented in the service layer.
type ExcelFolderDto struct {
	Batch         int     `json:"batch"`
	Index         int     `json:"index"`
	Path          string  `json:"path" validate:"required"`
	Alias         string  `json:"alias"`
	DisplayName   string  `json:"displayName"`
	TemplateAlias string  `json:"templateAlias"`
	Fields        string  `json:"fields"` // JSON string of FieldDefine array
	Description   *string `json:"description"`
}

func (e *ExcelFolderDto) CreateTopic() *types.CreateTopicDto {
	return &types.CreateTopicDto{
		Index:       e.Index,
		Batch:       e.Batch,
		Path:        e.Path,
		Name:        PathUtil.GetName(e.Path),
		Alias:       e.Alias,
		DisplayName: &e.DisplayName,
		Description: e.Description,
	}
}

func (e *ExcelFolderDto) GainBatchIndex() string {
	return fmt.Sprintf("%d-%d", e.Batch, e.Index)
}

func (e *ExcelFolderDto) Trim() {
	e.Path = strings.TrimSpace(e.Path)
	e.Alias = strings.TrimSpace(e.Alias)
	e.DisplayName = strings.TrimSpace(e.DisplayName)
	e.TemplateAlias = strings.TrimSpace(e.TemplateAlias)
	if e.Description != nil {
		e.Description = base.V2p(strings.TrimSpace(*e.Description))
	}
}

// ExcelNamespaceBaseDto represents the base DTO for Excel namespace operations
type ExcelNamespaceBaseDto struct {
	Topic    string `json:"topic,omitzero"`    // 主题
	DataType int    `json:"dataType,omitzero"` // 数据类型：1--时序库 2--关系库
	Fields   string `json:"fields,omitzero"`   // 字段定义（JSON 字符串）
}
