package service

import (
	"backend/internal/types"
	"context"
)

type UnsUpdateService struct {
}

func (l *UnsUpdateService) UpdateDetail(ctx context.Context, req *types.UpdateUnsDto) (resp *types.StringResult, err error) {

	return
}

func (l *UnsUpdateService) UpdateName(ctx context.Context, req *types.UpdateNameVo) (resp *types.StringResult, err error) {
	return
}

func (l *UnsUpdateService) SubscribeModel(ctx context.Context, req *types.SubscribeModelReq) (resp *types.ResultVO, err error) {
	return
}
