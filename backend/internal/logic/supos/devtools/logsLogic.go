package devtools

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// logs
func NewLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogsLogic {
	return &LogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogsLogic) Logs() error {
	// todo: add your logic here and delete this line

	return nil
}
