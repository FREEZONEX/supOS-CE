package relationDB

import (
	"time"
)

// NoderedFlow represents both source and event flows stored in the shared table.
type NoderedFlow struct {
	ID       int64  `gorm:"column:id;primaryKey;type:bigint;"`
	FlowID   string `gorm:"column:flow_id;size:64;index;comment:node-red flow id"`
	FlowName string `gorm:"column:flow_name;size:128;uniqueIndex:idx_uns_node_flow_name;comment:名称唯一"`
	// Use cross-DB compatible type. Postgres has no LONGTEXT, so use TEXT.
	FlowData    string `gorm:"column:flow_data;type:text;comment:节点json(不含tab)"`
	FlowStatus  string `gorm:"column:flow_status;size:32;comment:状态"`
	Template    string `gorm:"column:template;size:64;comment:模板来源"`
	Description string `gorm:"column:description;size:512;comment:描述"`
	// GORM many2many association with UnsNamespace via supos_node_flow_models.
	Nodes      []UnsNamespace `gorm:"many2many:supos_node_flow_models;foreignKey:ID;joinForeignKey:ParentID;References:Alias;joinReferences:Alias"`
	CreateTime time.Time      `gorm:"column:create_time;autoCreateTime;" json:"create_time"`
	UpdateTime time.Time      `gorm:"column:update_time;autoUpdateTime;" json:"update_time"`
	Creator    string         `gorm:"column:creator;comment:创建者"`
}

func (NoderedFlow) TableName() string {
	return "supos_node_flows"
}

type NoderedFlowNode struct {
	ParentID   int64     `gorm:"column:parent_id;comment:flow表id"`
	Alias      string    `gorm:"column:alias;"`
	Topic      string    `gorm:"column:topic;"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
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
