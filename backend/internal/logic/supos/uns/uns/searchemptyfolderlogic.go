package uns

import (
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

func (l *SearchEmptyFolderLogic) SearchEmptyFolder(req *types.EmptyFolderReq) (resp []types.CreateTopicDto, err error) {
	// todo: add your logic here and delete this line

	return
}
