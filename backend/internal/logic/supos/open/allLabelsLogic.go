// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"
	"strconv"

	dao "backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllLabelsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标签列表
func NewAllLabelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLabelsLogic {
	return &AllLabelsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllLabelsLogic) AllLabels(req *types.LabelQueryDto) (resp *types.LabelPageResult, err error) {
	var labelRepo dao.UnsLabelRepo
	db := dao.GetDb(l.ctx)

	var labels []*dao.UnsLabel
	var total int64

	// 如果有分页参数
	if req.PageNo > 0 && req.PageSize > 0 {
		labels, err = labelRepo.ListAll(db, req.PageNo, req.PageSize)
		if err != nil {
			return &types.LabelPageResult{
				Code: 500,
			}, nil
		}

		// 获取总数（简化处理）
		total = int64(len(labels))
	} else {
		// 不分页，查询所有
		labels, err = labelRepo.ListAll(db, 1, 10000)
		if err != nil {
			return &types.LabelPageResult{
				Code: 500,
			}, nil
		}
		total = int64(len(labels))
	}

	// 转换为 LabelOpenVo
	labelVos := make([]types.LabelOpenVo, len(labels))
	for i, label := range labels {
		labelVos[i] = types.LabelOpenVo{
			Id:         strconv.FormatInt(label.ID, 10),
			LabelName:  label.LabelName,
			CreateTime: label.CreateAt.UnixMilli(),
		}
	}

	return &types.LabelPageResult{
		Code:     200,
		Data:     labelVos,
		PageNo:   int64(req.PageNo),
		PageSize: int64(req.PageSize),
		Total:    total,
	}, nil
}
