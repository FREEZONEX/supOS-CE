package dto

import "time"

type UnsSearchCondition struct {
	SearchType       int                    `json:"searchType"`
	Deep             *int                   `json:"deep,omitempty"`
	ParentId         *int64                 `json:"parentId,omitempty"`
	PathType         *int16                 `json:"pathType,omitempty"`
	Keyword          string                 `json:"keyword,omitempty"`
	Alias            string                 `json:"alias,omitempty"`
	ParentAlias      string                 `json:"parentAlias,omitempty"`
	ParentAliasList  []string               `json:"parentAliasList,omitempty"`
	AliasList        []string               `json:"aliasList,omitempty"`
	Name             string                 `json:"name,omitempty"`
	DisplayName      string                 `json:"displayName,omitempty"`
	Path             string                 `json:"path,omitempty"`
	LayRec           string                 `json:"layRec,omitempty"`
	PathList         []string               `json:"pathList,omitempty"`
	Description      string                 `json:"description,omitempty"`
	TemplateName     string                 `json:"templateName,omitempty"`
	TemplateAlias    string                 `json:"templateAlias,omitempty"`
	TemplateID       *int64                 `json:"templateId,omitempty"`
	DataType         *int                   `json:"dataType,omitempty"`
	LabelName        string                 `json:"labelName,omitempty"`
	UpdateStartTime  *time.Time             `json:"updateStartTime,omitempty"`
	UpdateEndTime    *time.Time             `json:"updateEndTime,omitempty"`
	CreateStartTime  *time.Time             `json:"createStartTime,omitempty"`
	CreateEndTime    *time.Time             `json:"createEndTime,omitempty"`
	Extend           map[string]interface{} `json:"extend,omitempty"`
	WithValues       *bool                  `json:"withValues,omitempty"`
	ReturnParentInfo *bool                  `json:"returnParentInfo,omitempty"`
	ShowRec          *bool                  `json:"-"`
	FilterFolder     *bool                  `json:"-"`
}
