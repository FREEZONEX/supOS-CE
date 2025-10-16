package label

import (
	"context"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除标签
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.WithID) (resp *types.Empty, err error) {
	if req.ID <= 0 {
		return nil, errors.Parameter.WithMsg("id无效")
	}
	db := relationDB.NewUnsLabelRepo(stores.GetCommonConn(l.ctx))
	if err = db.Delete(l.ctx, req.ID); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}
