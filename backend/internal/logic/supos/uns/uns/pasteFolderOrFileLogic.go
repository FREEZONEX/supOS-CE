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

type PasteFolderOrFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 粘贴文件夹和文件
func NewPasteFolderOrFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PasteFolderOrFileLogic {
	return &PasteFolderOrFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PasteFolderOrFileLogic) PasteFolderOrFile(req *types.PasteRequestVO) (resp *types.CreateUnsResp, err error) {
	resp, err = spring.GetBean[*service.UnsAddService]().PasteFolderOrFile(l.ctx, req)
	return
}
