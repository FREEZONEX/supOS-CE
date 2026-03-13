// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLabelSchemaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询标签schema 元数据结构
func NewGetLabelSchemaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLabelSchemaLogic {
	return &GetLabelSchemaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLabelSchemaLogic) GetLabelSchema() error {
	// todo: add your logic here and delete this line

	return nil
}
