package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type NameDuplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 校验指定文件夹夹是否已存在文件夹、文件名称
func NewNameDuplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NameDuplicationLogic {
	return &NameDuplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NameDuplicationLogic) NameDuplication() error {
	// todo: add your logic here and delete this line

	return nil
}
