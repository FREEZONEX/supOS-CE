package uns

import (
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

func (l *PasteFolderOrFileLogic) PasteFolderOrFile(req *types.PasteRequestVO) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
