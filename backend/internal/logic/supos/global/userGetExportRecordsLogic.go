package global

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetExportRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取导出记录
func NewUserGetExportRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetExportRecordsLogic {
	return &UserGetExportRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGetExportRecordsLogic) UserGetExportRecords() error {
	// todo: add your logic here and delete this line

	return nil
}
