package relationDB

import (
	"time"
)

type NoderedSourceFlow struct {
	ID       int64  `gorm:"column:id;primaryKey;type:bigint;"`
	FlowID   string `gorm:"column:flow_id;size:64;index;comment:node-red flow id"`
	FlowName string `gorm:"column:flow_name;size:128;uniqueIndex:idx_uns_node_flow_name;comment:名称唯一"`
	// Use cross-DB compatible type. Postgres has no LONGTEXT, so use TEXT.
	FlowData    string `gorm:"column:flow_data;type:text;comment:节点json(不含tab)"`
	FlowStatus  string `gorm:"column:flow_status;size:32;comment:状态"`
	Template    string `gorm:"column:template;size:64;comment:模板来源"`
	Description string `gorm:"column:description;size:512;comment:描述"`
	// FlowType    int32  `gorm:"column:flow_type;comment:'1-source 2-event'"`
	// GORM many2many association with UnsNamespaceNodeInfo via uns_node_flow_model
	Nodes      []UnsNamespace `gorm:"many2many:supos_node_flow_models;foreignKey:ID;joinForeignKey:ParentID;References:Alias;joinReferences:Alias"`
	CreateTime time.Time      `gorm:"column:create_time;default:now()" json:"create_time"`
	UpdateTime time.Time      `gorm:"column:update_time" json:"update_time"`
}

func (NoderedSourceFlow) TableName() string {
	return "supos_node_flows"
}

type NoderedSourceFlowNode struct {
	ParentID   int64     `gorm:"column:parent_id;comment:flow表id"`
	Alias      string    `gorm:"column:alias;"`
	Topic      string    `gorm:"column:topic;"`
	CreateTime time.Time `gorm:"column:create_time;default:now()" json:"create_time"`
}

func (NoderedSourceFlowNode) TableName() string {
	return "supos_node_flow_models"
}

type NoderedEventFlow struct {
	ID       int64  `gorm:"column:id;primaryKey;type:bigint;"`
	FlowID   string `gorm:"column:flow_id;size:64;index;comment:node-red flow id"`
	FlowName string `gorm:"column:flow_name;size:128;uniqueIndex:idx_uns_node_flow_name;comment:名称唯一"`
	// Use cross-DB compatible type. Postgres has no LONGTEXT, so use TEXT.
	FlowData    string `gorm:"column:flow_data;type:text;comment:节点json(不含tab)"`
	FlowStatus  string `gorm:"column:flow_status;size:32;comment:状态"`
	Template    string `gorm:"column:template;size:64;comment:模板来源"`
	Description string `gorm:"column:description;size:512;comment:描述"`
	// FlowType    int32  `gorm:"column:flow_type;comment:'1-source 2-event'"`
	// GORM many2many association with UnsNamespaceNodeInfo via uns_node_flow_model
	Nodes      []UnsNamespace `gorm:"many2many:supos_event_flow_models;foreignKey:ID;joinForeignKey:ParentID;References:Alias;joinReferences:Alias"`
	CreateTime time.Time      `gorm:"column:create_time;default:now()" json:"create_time"`
	UpdateTime time.Time      `gorm:"column:update_time" json:"update_time"`
}

func (NoderedEventFlow) TableName() string {
	return "supos_event_flows"
}

type NoderedEventFlowNode struct {
	ParentID   int64     `gorm:"column:parent_id;comment:flow表id"`
	Alias      string    `gorm:"column:alias;"`
	Topic      string    `gorm:"column:topic;"`
	CreateTime time.Time `gorm:"column:create_time;default:now()" json:"create_time"`
}

func (NoderedEventFlowNode) TableName() string {
	return "supos_event_flow_models"
}
