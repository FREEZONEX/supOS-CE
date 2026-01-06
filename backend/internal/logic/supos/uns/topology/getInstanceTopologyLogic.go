// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package topology

import (
	"backend/internal/logic/supos/uns/topology/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInstanceTopologyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInstanceTopologyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInstanceTopologyLogic {
	return &GetInstanceTopologyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInstanceTopologyLogic) GetInstanceTopology(req *types.GetInstanceTopologyReq) (*types.JsonResult, error) {
	topologyService := spring.GetBean[*service.UnsTopologyService]()

	// Get topology data for the specific instance.
	// Currently returns default data with all nodes in SUCCESS state
	topologyData := topologyService.GainTopologyDataOfFile(req.UnsId)

	return &types.JsonResult{
		Code: http.StatusOK,
		Msg:  "success",
		Data: topologyData,
	}, nil
}
