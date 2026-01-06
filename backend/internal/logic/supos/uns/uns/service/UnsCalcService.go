package service

import (
	"backend/internal/logic/supos/uns/uns/bo"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"context"
)

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
	dto bo.UnsDetail) {

}
func (s UnsCalcService) detectReferencedCalcInstance(ctx context.Context, files []*dao.UnsNamespace, modelPath string, delFields []*types.FieldDefine) (affectedFiles []*dao.UnsNamespace) {

	return
}
