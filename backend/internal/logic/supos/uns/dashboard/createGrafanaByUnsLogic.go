package dashboard

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGrafanaByUnsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// createGrafanaByUns
func NewCreateGrafanaByUnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGrafanaByUnsLogic {
	return &CreateGrafanaByUnsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGrafanaByUnsLogic) CreateGrafanaByUns() error {
	// todo: add your logic here and delete this line

	return nil
}
