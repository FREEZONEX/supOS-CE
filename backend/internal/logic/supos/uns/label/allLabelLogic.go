package label

import (
	"context"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/stores"
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
	db := relationDB.NewUnsLabelRepo(stores.GetCommonConn(l.ctx))
	dbFilter := relationDB.UnsLabelFilter{
		LabelName: req.Key,
	}
	list, err := db.FindByFilter(l.ctx, dbFilter, nil)
	if err != nil {
		return nil, err
	}
	resp = &types.UnsLabelListResp{
		List: make([]*types.UnsLabel, 0),
	}
	for _, item := range list {
		resp.List = append(resp.List, &types.UnsLabel{
			ID:        item.ID,
			LabelName: item.LabelName,
			CreateAt:  item.CreateAt.UnixMilli(),
		})
	}
	return
}
