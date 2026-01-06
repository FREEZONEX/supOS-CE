package label

import (
	"backend/internal/logic/supos/uns/label/service"
	"backend/share/spring"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标签列表
func NewAllLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLabelLogic {
	return &AllLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllLabelLogic) AllLabel(req *types.UnsLabelListReq) (resp *types.UnsLabelListResp, err error) {
	sv := spring.GetBean[*service.UnsLabelService]()
	return sv.AllLabel(l.ctx, req)
}
