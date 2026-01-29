package appkey

import (
	"context"

	"backend/internal/logic/supos/appkey/service"
	"backend/internal/svc"
	"backend/share/spring"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create() error {
	appKeyService := spring.GetBean[*service.AppKeyService]()
	success, err := appKeyService.CreateSecretKey(l.ctx)
	if err != nil {
		return err
	}

	if !success {
		return errors.NewCodeError(400, "创建密钥失败")
	}

	return nil
}
