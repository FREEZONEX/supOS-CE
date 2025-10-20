// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataSrc2UnsFieldsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 外部数据源表的字段定义转uns字段定义
func NewDataSrc2UnsFieldsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataSrc2UnsFieldsLogic {
	return &DataSrc2UnsFieldsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DataSrc2UnsFieldsLogic) DataSrc2UnsFields(req *types.DbFieldsInfoVo) (resp *types.DbFieldsInfoVoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
