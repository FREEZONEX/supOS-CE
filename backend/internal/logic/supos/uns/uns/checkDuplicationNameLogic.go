// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckDuplicationNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 校验指定文件夹夹是否已存在文件夹、文件名称
func NewCheckDuplicationNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckDuplicationNameLogic {
	return &CheckDuplicationNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckDuplicationNameLogic) CheckDuplicationName(req *types.CheckDuplicationNameReq) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
