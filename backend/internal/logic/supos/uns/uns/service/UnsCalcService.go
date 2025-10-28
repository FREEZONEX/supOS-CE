package service

import "backend/internal/types"

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
