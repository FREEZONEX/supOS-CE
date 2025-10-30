// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package label

import (
	"backend/internal/logic/supos/uns/label/service"
	"backend/share/spring"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageListUnsByLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取标签下的文件列表
func NewPageListUnsByLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageListUnsByLabelLogic {
	return &PageListUnsByLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageListUnsByLabelLogic) PageListUnsByLabel(req *types.LabelPageReq) (resp *types.UnsByLabelPageResp, err error) {
	sv := spring.GetBean[*service.UnsLabelService]()
	return sv.PageListUnsByLabel(l.ctx, req)
}
