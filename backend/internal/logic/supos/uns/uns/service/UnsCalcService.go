package service

import "backend/internal/types"

// 实时计算服务（占位）
type UnsCalcService struct {
}

func (s UnsCalcService) CheckFileField(dto *types.CreateTopicDto) string {
	return ""
}

func (s UnsCalcService) CheckRefers(unsDto *types.CreateTopicDto) string {
	return ""
}

func (s UnsCalcService) CheckComplexExpression(unsDto *types.CreateTopicDto) string {
	return ""
}
func (s UnsCalcService) setRefersAndExpression(fs []*types.InstanceField,
	expression string,
	calculationType *int32,
	protocolMap map[string]interface{},
	dto *types.InstanceDetail) {

}
