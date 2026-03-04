package relationDB

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/clause"
)

type NodeFlowModelMapper struct{}

func (NodeFlowModelMapper) SaveBatch(ctx context.Context, list []*NoderedFlowNode) error {
	db := GetDb(ctx)
	return db.Model(&NoderedFlowNode{}).Clauses(clause.OnConflict{DoNothing: true}).Save(list).Error
}
func (NodeFlowModelMapper) ListUnsExistsFlows(ctx context.Context, aliasList []string) (aliasMap map[string]bool) {
	var refs []*NoderedFlowNode
	db := GetDb(ctx)
	err := db.Where("alias IN ?", aliasList).Select("alias").Find(&refs).Error
	if err != nil {
		logx.Errorf("failed to select flows refs by aliases: %v", err)
		return nil
	} else if len(refs) == 0 {
		return nil
	}
	aliasMap = make(map[string]bool, len(refs))
	for _, ref := range refs {
		aliasMap[ref.Alias] = true
	}
	return aliasMap
}
