package service

import "backend/internal/common/dto"

type UnsCalcService struct {
}

func (s UnsCalcService) CheckFileField(dto *dto.CreateTopicDto) string {
	return ""
}

func (s UnsCalcService) CheckRefers(unsDto *dto.CreateTopicDto) string {
	return ""
}

func (s UnsCalcService) CheckComplexExpression(unsDto *dto.CreateTopicDto) string {
	return ""
}
