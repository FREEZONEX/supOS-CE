// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"backend/internal/logic/supos/uns/label/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type MakeLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量文件打标签
func NewMakeLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MakeLabelLogic {
	return &MakeLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MakeLabelLogic) MakeLabel(req *types.MakeLabelDtoArray) (resp *types.ResultVO, err error) {
	resp = &types.ResultVO{Code: 200, Msg: "ok"}
	unsLabelService := spring.GetBean[*service.UnsLabelService]()
	rs, err := unsLabelService.BatchMakeLabels(l.ctx, req.Items)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	} else if rs != nil {
		resp.Code = rs.Code
		resp.Msg = rs.Msg
		resp.Data = rs.Data
	}
	return resp, nil
}
