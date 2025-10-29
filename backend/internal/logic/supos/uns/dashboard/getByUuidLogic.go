// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetByUuidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据 UID 获取 Dashboard
func NewGetByUuidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByUuidLogic {
	return &GetByUuidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetByUuidLogic) GetByUuid(req *types.UuidRequest) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
