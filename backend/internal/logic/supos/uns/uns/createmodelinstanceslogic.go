package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateModelInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量创建文件夹和文件
func NewCreateModelInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateModelInstancesLogic {
	return &CreateModelInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateModelInstancesLogic) CreateModelInstances(req *types.BatchCreateReq) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
