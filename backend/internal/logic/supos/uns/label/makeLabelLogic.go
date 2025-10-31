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

type MakeLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件打标签
func NewMakeLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MakeLabelLogic {
	return &MakeLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MakeLabelLogic) MakeLabel(req *types.MakeLabelReq) (resp *types.BaseResult, err error) {
	sv := spring.GetBean[*service.UnsLabelService]()
	err = sv.ClearAndMakeLabels(l.ctx, req.UnsId, req.LabelList)
	if err == nil {
		resp = &types.BaseResult{Code: 200, Msg: "ok"}
	}
	return
}
