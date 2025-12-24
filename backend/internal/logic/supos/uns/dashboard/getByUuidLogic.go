// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

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

func (l *GetByUuidLogic) GetByUuid(uuid string) (*types.JsonResult, error) {
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
