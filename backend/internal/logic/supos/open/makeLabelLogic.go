// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
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

// 批量文件打标签
func NewMakeLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MakeLabelLogic {
	return &MakeLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MakeLabelLogic) MakeLabel(req *types.MakeLabelDtoArray) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
