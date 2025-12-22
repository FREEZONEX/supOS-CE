package dashboard

import (
	unsservice "backend/internal/logic/supos/uns/uns/service"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"net/http"

	"gitee.com/unitedrhino/share/i18ns"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetByUnsLogic struct {
	logx.Logger
	ctx                context.Context
	svcCtx             *svc.ServiceContext
	dashboardRefMapper *relationDB.DashboardRefMapper
	unsQueryService    *unsservice.UnsQueryService
}

func NewGetByUnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByUnsLogic {
	return &GetByUnsLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		unsQueryService: spring.GetBean[*unsservice.UnsQueryService](),
	}
}

func (l *GetByUnsLogic) GetByUns(unsAlias string) (*types.JsonResult, error) {
	// TODO: As identified before, the DTO from GetModelDefinition lacks the 'Refers' field.
	// This logic is simplified until the UNS service provides the necessary details.
	db := relationDB.GetDb(l.ctx)
	dashboard, err := l.dashboardRefMapper.GetByUns(db, unsAlias)
	if err != nil {
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, nil
	}
	if dashboard == nil {
		return &types.JsonResult{
			Code: http.StatusBadRequest,
			Msg:  i18ns.LocalizeMsg("uns.dashboard.not.exit"),
		}, nil
	}
	return &types.JsonResult{
		Code: http.StatusOK,
		Msg:  "success",
		Data: dashboard,
	}, nil
}
