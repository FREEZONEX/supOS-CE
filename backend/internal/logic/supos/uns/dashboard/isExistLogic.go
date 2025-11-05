package dashboard

import (
	"backend/internal/common/errors"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsExistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIsExistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsExistLogic {
	return &IsExistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsExistLogic) IsExist(alias string) (map[string]any, error) {
	uuid := grafanautil.GetDashboardUUIDByAlias(alias)
	dbJSON, err := grafanautil.GetDashboardByUUID(uuid)
	if err != nil || dbJSON == nil {
		return nil, errors.NewBuzError(400, "uns.dashboard.not.exit")
	}
	return dbJSON, nil
}
