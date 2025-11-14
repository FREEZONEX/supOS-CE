// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package topology

import (
	"backend/internal/common/utils/topologylog"
	"backend/internal/logic/supos/uns/topology/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type MockInstanceTopologyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMockInstanceTopologyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MockInstanceTopologyLogic {
	return &MockInstanceTopologyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MockInstanceTopologyLogic) MockInstanceTopology(req *types.MockInstanceTopologyReq) error {
	topologylog.Log(req.UnsId, req.Node, topologylog.EventCodeError, "sd")

	topologyService := spring.GetBean[*service.UnsTopologyService]()
	topologyService.UpdateTopologyState(req.Node, string(topologylog.EventCodeError))

	return nil
}
