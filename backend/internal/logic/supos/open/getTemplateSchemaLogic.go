// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTemplateSchemaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询模板schema 元数据结构
func NewGetTemplateSchemaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTemplateSchemaLogic {
	return &GetTemplateSchemaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTemplateSchemaLogic) GetTemplateSchema() error {
	// todo: add your logic here and delete this line

	return nil
}
