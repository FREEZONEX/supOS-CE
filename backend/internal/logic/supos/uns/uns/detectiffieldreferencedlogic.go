package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetectIfFieldReferencedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 预先判断是否有属性关联
func NewDetectIfFieldReferencedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetectIfFieldReferencedLogic {
	return &DetectIfFieldReferencedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetectIfFieldReferencedLogic) DetectIfFieldReferenced(req *types.UpdateModeRequestVo) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
