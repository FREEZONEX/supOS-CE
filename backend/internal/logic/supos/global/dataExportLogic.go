package global

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 全局数据导出
func NewDataExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataExportLogic {
	return &DataExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DataExportLogic) DataExport() error {
	// todo: add your logic here and delete this line

	return nil
}
