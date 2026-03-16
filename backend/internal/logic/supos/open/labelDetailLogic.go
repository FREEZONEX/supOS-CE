// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/common/I18nUtils"
	dao "backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LabelDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标签详情
func NewLabelDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LabelDetailLogic {
	return &LabelDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LabelDetailLogic) LabelDetail(id int64) (resp *types.ResultVO, err error) {
	var labelRepo dao.UnsLabelRepo
	db := dao.GetDb(l.ctx)

	label, err := labelRepo.SelectById(db, id)
	if err != nil {
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.operation.failed") + ": " + err.Error(),
		}, nil
	}

	if label == nil {
		return &types.ResultVO{
			Code: 404,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.label.not.exists"),
		}, nil
	}

	return &types.ResultVO{
		Code: 200,
		Msg:  "ok",
		Data: types.LabelVo{
			ID:         label.ID,
			LabelName:  label.LabelName,
			CreateTime: label.CreateAt.UnixMilli(),
		},
	}, nil
}
