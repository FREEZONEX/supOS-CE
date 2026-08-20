package repo

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type flowNode struct {
	ParentID int64 `gorm:"column:parent_id;default:0;index:idx_uns_flow_node_parent"`
	NodeID   int64 `gorm:"column:node_id;primaryKey"`
	OnlyTime
}

func (flowNode) TableName() string { return "uns_nodered_flow_node" }

var ErrFlowFolderCannotDeploy = errors.New("flow folder cannot deploy")





func (r *FlowRepo) FlowNodeIDs(ctx context.Context, flowID int64) ([]int64, error) {
	var out []int64
	err := r.db.WithContext(ctx).Table("uns_nodered_flow_node").
		Where("parent_id = ?", flowID).Order("node_id").Pluck("node_id", &out).Error
	return out, err
}

func (r *FlowRepo) BindFlowIDMapByNodeIDs(ctx context.Context, nodeIDs []int64) (map[int64]int64, error) {
	out := make(map[int64]int64, len(nodeIDs))
	if len(nodeIDs) == 0 {
		return out, nil
	}
	cleanNodeIDs := make([]int64, 0, len(nodeIDs))
	seen := make(map[int64]struct{}, len(nodeIDs))
	for _, nodeID := range nodeIDs {
		if nodeID <= 0 {
			continue
		}
		if _, ok := seen[nodeID]; ok {
			continue
		}
		seen[nodeID] = struct{}{}
		cleanNodeIDs = append(cleanNodeIDs, nodeID)
	}
	if len(cleanNodeIDs) == 0 {
		return out, nil
	}
	var rows []struct {
		NodeID int64 `gorm:"column:node_id"`
		FlowID int64 `gorm:"column:flow_id"`
	}
	err := r.db.WithContext(ctx).
		Table("uns_nodered_flow_node AS m").
		Select("m.node_id, m.parent_id AS flow_id").
		Joins("JOIN uns_nodered_flow AS f ON f.id = m.parent_id").
		Where("m.node_id IN ?", cleanNodeIDs).
		Where("f.deleted_time = 0").
		Order("m.node_id ASC, m.created_time DESC, m.parent_id DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row.NodeID == 0 || row.FlowID == 0 {
			continue
		}
		if _, ok := out[row.NodeID]; ok {
			continue
		}
		out[row.NodeID] = row.FlowID
	}
	return out, nil
}

func (r *FlowRepo) FlowIDsByNodeIDs(ctx context.Context, nodeIDs []int64, flowType int) ([]int64, error) {
	if len(nodeIDs) == 0 {
		return []int64{}, nil
	}
	cleanNodeIDs := make([]int64, 0, len(nodeIDs))
	seen := make(map[int64]struct{}, len(nodeIDs))
	for _, nodeID := range nodeIDs {
		if nodeID <= 0 {
			continue
		}
		if _, ok := seen[nodeID]; ok {
			continue
		}
		seen[nodeID] = struct{}{}
		cleanNodeIDs = append(cleanNodeIDs, nodeID)
	}
	if len(cleanNodeIDs) == 0 {
		return []int64{}, nil
	}
	q := r.db.WithContext(ctx).
		Table("uns_nodered_flow_node AS m").
		Joins("JOIN uns_nodered_flow AS f ON f.id = m.parent_id").
		Where("m.node_id IN ?", cleanNodeIDs).
		Where("f.deleted_time = 0").
		Where("f.node_type <> 1")
	if flowType != 0 {
		q = q.Where("f.flow_type = ?", flowType)
	}
	var out []int64
	err := q.Distinct("f.id").Order("f.id").Pluck("f.id", &out).Error
	return out, err
}

func replaceFlowNodes(db *gorm.DB, flowID int64, nodeIDs []int64, now int64) error {
	cleanNodeIDs := make([]int64, 0, len(nodeIDs))
	seen := make(map[int64]struct{}, len(nodeIDs))
	for _, nodeID := range nodeIDs {
		if nodeID == 0 {
			continue
		}
		if _, ok := seen[nodeID]; ok {
			continue
		}
		seen[nodeID] = struct{}{}
		cleanNodeIDs = append(cleanNodeIDs, nodeID)
	}
	if err := db.Where("parent_id = ?", flowID).Delete(&flowNode{}).Error; err != nil {
		return err
	}
	if len(cleanNodeIDs) > 0 {
		if err := db.Where("node_id IN ?", cleanNodeIDs).Delete(&flowNode{}).Error; err != nil {
			return err
		}
	}
	for _, nodeID := range cleanNodeIDs {
		row := flowNode{ParentID: flowID, NodeID: nodeID}
		ts := repoTimeFromMilli(now)
		row.CreatedTime = ts
		row.UpdatedTime = ts
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "node_id"}},
			DoUpdates: clause.Assignments(map[string]any{"updated_time": ts}),
		}).Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}
