package postgresql

import (
	"backend/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostgresqlHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewPostgresqlHandle(ctx context.Context, svcCtx *svc.ServiceContext) *PostgresqlHandle {
	return &PostgresqlHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}
