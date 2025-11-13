// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package topology

import (
	"backend/internal/logic/supos/uns/topology/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"encoding/json"

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
	topologyService := spring.GetBean[*service.UnsTopologyService]()

	// Parse the JSON string and return the typed response
	jsonStr := topologyService.GetLastMsg()

	resp := &types.GetInstanceTopologyResp{}
	if err := json.Unmarshal([]byte(jsonStr), resp); err != nil {
		logx.Errorf("failed to unmarshal topology data: %v", err)
		// Return empty response on error
		return &types.GetInstanceTopologyResp{Data: []types.InstanceTopologyData{}}, nil
	}

	return resp, nil
}
