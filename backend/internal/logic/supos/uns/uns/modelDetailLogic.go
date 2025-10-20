package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ModelDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询文件夹详情
func NewModelDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModelDetailLogic {
	return &ModelDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModelDetailLogic) ModelDetail() error {
	// todo: add your logic here and delete this line

	return nil
}
