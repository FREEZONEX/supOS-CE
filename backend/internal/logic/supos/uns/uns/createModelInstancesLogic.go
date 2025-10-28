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

type CreateModelInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量创建文件夹和文件
func NewCreateModelInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateModelInstancesLogic {
	return &CreateModelInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateModelInstancesLogic) CreateModelInstances(req *types.BatchCreateReq) (resp *types.ResultVO, err error) {
	errTipMap := spring.GetBean[*service.UnsAddService]().CreateModelAndInstance(l.ctx, req.List, req.FromImport)
	if len(errTipMap) == 0 {
		return &types.ResultVO{Code: 200, Msg: "ok"}, nil
	}
	return &types.ResultVO{Code: 206, Msg: "ok", Data: errTipMap}, nil
}
