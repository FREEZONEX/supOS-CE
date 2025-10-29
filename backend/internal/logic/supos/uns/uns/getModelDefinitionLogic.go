// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/logic/supos/uns/uns/service"
	"backend/share/spring"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetModelDefinitionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询文件夹详情
func NewGetModelDefinitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetModelDefinitionLogic {
	return &GetModelDefinitionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetModelDefinitionLogic) GetModelDefinition(req *types.ModelDetailReq) (resp *types.ModelDetailResp, err error) {
	resp, err = spring.GetBean[*service.UnsQueryService]().GetModelDefinition(l.ctx, req, "")
	return
}
