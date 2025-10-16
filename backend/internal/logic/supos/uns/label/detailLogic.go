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

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询详情
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.WithID) (resp *types.UnsLabel, err error) {
	if req.ID <= 0 {
		return nil, errors.Parameter.WithMsg("id无效")
	}
	db := relationDB.NewUnsLabelRepo(stores.GetCommonConn(l.ctx))
	item, err := db.FindOne(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &types.UnsLabel{
		ID:        item.ID,
		LabelName: item.LabelName,
		CreateAt:  item.CreateAt.UnixMilli(),
	}, nil
}
