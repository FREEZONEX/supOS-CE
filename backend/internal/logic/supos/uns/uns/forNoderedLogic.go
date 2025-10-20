package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ForNoderedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量创建文件夹和文件(node-red导入专用)
func NewForNoderedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForNoderedLogic {
	return &ForNoderedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForNoderedLogic) ForNodered() error {
	// todo: add your logic here and delete this line

	return nil
}
