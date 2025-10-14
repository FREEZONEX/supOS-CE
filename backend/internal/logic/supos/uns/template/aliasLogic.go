package template

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AliasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据别名查询模板详情
func NewAliasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AliasLogic {
	return &AliasLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AliasLogic) Alias(req *types.WithID) error {
	// todo: add your logic here and delete this line

	return nil
}
