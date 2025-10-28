package global

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserExportRecordConfirmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 确认导出记录
func NewUserExportRecordConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserExportRecordConfirmLogic {
	return &UserExportRecordConfirmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserExportRecordConfirmLogic) UserExportRecordConfirm() error {
	// todo: add your logic here and delete this line

	return nil
}
