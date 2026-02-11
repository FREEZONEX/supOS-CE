package relationDB

import (
	"context"

	"gorm.io/gorm/clause"
)

type NodeFlowModelMapper struct{}

func (NodeFlowModelMapper) SaveBatch(ctx context.Context, list []*NoderedFlowNode) error {
	db := GetDb(ctx)
	return db.Model(&NoderedFlowNode{}).Clauses(clause.OnConflict{DoNothing: true}).Save(list).Error
}
