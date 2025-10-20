package label

import (
	"context"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"strings"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改标签
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UnsLabel) (resp *types.Empty, err error) {
	if req.ID <= 0 {
		return nil, errors.Parameter.WithMsg("id无效")
	}
	if strings.TrimSpace(req.LabelName) == "" {
		return nil, errors.Parameter.WithMsg("labelName不能为空")
	}
	db := relationDB.NewUnsLabelRepo(l.ctx)
	item, err := db.FindOne(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}
	item.LabelName = req.LabelName
	if err = db.Update(l.ctx, item); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}
