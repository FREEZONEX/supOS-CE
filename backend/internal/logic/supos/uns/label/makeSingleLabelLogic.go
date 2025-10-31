// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package label

import (
	"backend/internal/logic/supos/uns/label/service"
	"backend/share/spring"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MakeSingleLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件打单个标签
func NewMakeSingleLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MakeSingleLabelLogic {
	return &MakeSingleLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MakeSingleLabelLogic) MakeSingleLabel(req *types.MakeSingleLabelReq) (resp *types.BaseResult, err error) {
	sv := spring.GetBean[*service.UnsLabelService]()
	return sv.MakeSingleLabel(l.ctx, req)
}
