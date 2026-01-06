// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/logic/supos/uns/uns/service"
	"backend/share/spring"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchEmptyFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询空文件夹
func NewSearchEmptyFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchEmptyFolderLogic {
	return &SearchEmptyFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchEmptyFolderLogic) SearchEmptyFolder(req *types.EmptyFolderReq) (resp *types.EmptyFolderResp, err error) {
	resp, err = spring.GetBean[*service.UnsQueryService]().SearchEmptyFolder(l.ctx, req)
	return
}
