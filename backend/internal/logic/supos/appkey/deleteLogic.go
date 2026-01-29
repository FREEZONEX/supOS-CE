package appkey

import (
	"context"

	"backend/internal/logic/supos/appkey/service"
	"backend/internal/svc"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(id int64) error {
	appKeyService := spring.GetBean[*service.AppKeyService]()
	err := appKeyService.DeleteSecretKey(l.ctx, id)
	if err != nil {
		return err
	}

	return nil
}
