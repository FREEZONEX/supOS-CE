package msg_consumer

import (
	"backend/internal/common/serviceApi"
	"backend/internal/types"
)

// UnsRealtimeCalcService 实时计算服务
type UnsRealtimeCalcService struct {
}

func (c UnsRealtimeCalcService) TryCalculate(
	defService serviceApi.IUnsDefinitionService,
	def *types.CreateTopicDto,
	data map[string]any) (calcDef *types.CreateTopicDto, calcData map[string]any, errMsg string) {

	if def == nil || len(def.RefUns) == 0 {
		return
	}
	//TODO  实时计算
	return
}
