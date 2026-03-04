package relationDB

import (
	"time"
)

// NoderedFlow represents both source and event flows stored in the shared table.
type NoderedFlow struct {
	ID         int64  `gorm:"column:id;primaryKey;type:bigint;" json:"id"`
	GroupId    *int64 `gorm:"column:group_id;comment:分组ID" json:"-"`
	ExportType string `gorm:"-" json:"exportType,omitempty"` // 导出类型：group-分组 默认：flow-流程
	FlowID     string `gorm:"column:flow_id;size:64;index;comment:node-red flow id" json:"flowId,omitempty"`
	FlowName   string `gorm:"column:flow_name;size:128;uniqueIndex:idx_uns_node_flow_name;comment:名称唯一" json:"name,omitempty"`
	// Use cross-DB compatible type. Postgres has no LONGTEXT, so use TEXT.
	FlowData    string `gorm:"column:flow_data;type:text;comment:节点json(不含tab)" json:"flowData,omitempty"`
	FlowStatus  string `gorm:"column:flow_status;size:32;comment:状态" json:"flowStatus,omitempty"`
	Template    string `gorm:"column:template;size:64;comment:模板来源" json:"template,omitempty"`
	Description string `gorm:"column:description;size:512;comment:描述" json:"description,omitempty"`
	// GORM many2many association with UnsNamespace via supos_node_flow_models.
	Nodes      []UnsNamespace `gorm:"many2many:supos_node_flow_models;foreignKey:ID;joinForeignKey:ParentID;References:Alias;joinReferences:Alias" json:"-"`
	CreateTime time.Time      `gorm:"column:create_time;autoCreateTime;" json:"createTime,omitzero"`
	UpdateTime time.Time      `gorm:"column:update_time;autoUpdateTime;" json:"updateTime,omitzero"`
	Creator    string         `gorm:"column:creator;comment:创建者" json:"creator,omitempty"`
	Children   []*NoderedFlow `gorm:"-" json:"children,omitempty"`
	Sort       int32          `gorm:"-" json:"sort,omitzero"`
}

const (
	TemplateTypeSrcFlow   = "node-red"
	TemplateTypeEventFlow = "event-flow"
)

func (NoderedFlow) TableName() string {
	return "supos_node_flows"
}

type NoderedFlowNode struct {
	ParentID   int64     `gorm:"column:parent_id;comment:flow表id" json:"pid"`
	Alias      string    `gorm:"column:alias;" json:"alias,omitempty"`
	Topic      string    `gorm:"column:topic;" json:"-"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"-"`
}

func (NoderedFlowNode) TableName() string {
	return "supos_node_flow_models"
}

// Backward-compatible aliases for legacy references.
type (
	NoderedSourceFlow     = NoderedFlow
	NoderedSourceFlowNode = NoderedFlowNode
)

type NoderedFlowTop struct {
	ID         int64     `gorm:"column:id;primaryKey;type:bigint;"`
	UserID     string    `gorm:"column:user_id;primaryKey;size:128;"`
	Mark       int       `gorm:"column:mark;default:1"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime;" json:"update_time"`
	MarkTime   time.Time `gorm:"column:mark_time;type:timestamptz;autoCreateTime"`
}

func (NoderedFlowTop) TableName() string {
	return "supos_node_flow_top_recodes"
}
