// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"backend/internal/common/errors"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetByUuidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetByUuidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByUuidLogic {
	return &GetByUuidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetByUuidLogic) GetByUuid(uuid string) (map[string]any, error) {
	dbJSON, err := grafanautil.GetDashboardByUUID(uuid)
	if err != nil || dbJSON == nil {
		return nil, errors.NewBuzError(400, "uns.dashboard.not.exit")
	}
	return dbJSON, nil
}
