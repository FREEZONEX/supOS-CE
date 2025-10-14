package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SystemConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取系统配置
func NewSystemConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SystemConfigLogic {
	return &SystemConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SystemConfigLogic) SystemConfig() error {
	// todo: add your logic here and delete this line

	return nil
}
