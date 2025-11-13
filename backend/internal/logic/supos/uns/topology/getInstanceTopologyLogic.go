// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package topology

import (
	"backend/internal/common/utils/topologylog"
	"backend/internal/svc"
	"backend/internal/types"
	"context"

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

func (l *GetInstanceTopologyLogic) GetInstanceTopology(req *types.GetInstanceTopologyReq) (*types.GetInstanceTopologyResp, error) {
	return l.createDefaultTopologyData(), nil
}

func (l *GetInstanceTopologyLogic) createDefaultTopologyData() *types.GetInstanceTopologyResp {
	topologyDatas := make([]types.InstanceTopologyData, len(topologylog.TopologyNodes))
	for i, node := range topologylog.TopologyNodes {
		topologyDatas[i] = types.InstanceTopologyData{
			TopologyNode: node,
			EventCode:    topologylog.EventCodeSuccess,
		}
	}
	return &types.GetInstanceTopologyResp{Data: topologyDatas}
}
