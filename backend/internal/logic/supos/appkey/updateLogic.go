package appkey

import (
	"context"

	"backend/internal/logic/supos/appkey/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateAppKeyReq) error {
	appKeyService := spring.GetBean[*service.AppKeyService]()
	err := appKeyService.UpdateSecretKey(l.ctx, req)
	if err != nil {
		return err
	}

	return nil
}
