package appkey

import (
	"context"

	"backend/internal/logic/supos/appkey/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List() (*types.ListAppKeyResp, error) {
	appKeyService := spring.GetBean[*service.AppKeyService]()
	list, err := appKeyService.GetSecretKeyList(l.ctx)
	if err != nil {
		return nil, err
	}

	return &types.ListAppKeyResp{
		Code: 200,
		Msg:  "success",
		Data: list,
	}, nil
}
