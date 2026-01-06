package dashboard

import (
	"backend/internal/common/utils/grafanautil"
	dao "backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"net/http"

	"gitee.com/unitedrhino/share/i18ns"
	"github.com/zeromicro/go-zero/core/logx"
)

type IsExistLogic struct {
	logx.Logger
	ctx           context.Context
	svcCtx        *svc.ServiceContext
	dashRefMapper dao.DashboardRefMapper
}

func NewIsExistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsExistLogic {
	return &IsExistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsExistLogic) IsExist(alias string) (*types.JsonResult, error) {
	db := dao.GetDb(l.ctx)
	ref, err := l.dashRefMapper.SelectByUnsAlias(db, alias)
	if err != nil || ref == nil {
		return &types.JsonResult{
			Code: http.StatusBadRequest,
			Msg:  i18ns.LocalizeMsg("uns.dashboard.not.exit"),
		}, err
	}
	uuid := ref.DashboardID
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
