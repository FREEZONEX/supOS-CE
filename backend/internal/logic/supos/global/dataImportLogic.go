package global

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 全局数据导入
func NewDataImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataImportLogic {
	return &DataImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DataImportLogic) DataImport() error {
	// todo: add your logic here and delete this line

	return nil
}
