package dashboard

import (
	"backend/internal/common/utils/grafanautil"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"net/http"

	"gitee.com/unitedrhino/share/i18ns"
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

func (l *IsExistLogic) IsExist(alias string) (*types.JsonResult, error) {
	uuid := grafanautil.GetDashboardUUIDByAlias(alias)
	dbJSON, err := grafanautil.GetDashboardByUUID(l.ctx, uuid)
	if err != nil || dbJSON == nil {
		return &types.JsonResult{
			Code: http.StatusBadRequest,
			Msg:  i18ns.LocalizeMsg("uns.dashboard.not.exit"),
		}, nil
	}
	return &types.JsonResult{
		Code: http.StatusOK,
		Msg:  "success",
		Data: dbJSON,
	}, nil
}
