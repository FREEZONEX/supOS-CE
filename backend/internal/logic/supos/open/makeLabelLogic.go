// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/uns/label/service"
	dao "backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type MakeLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量文件打标签
func NewMakeLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MakeLabelLogic {
	return &MakeLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MakeLabelLogic) MakeLabel(req *types.MakeLabelDtoArray) (resp *types.ResultVO, err error) {
	unsLabelService := spring.GetBean[*service.UnsLabelService]()

	// 批量创建标签（如果不存在）
	allLabelNames := make([]string, 0)
	for _, item := range req.Items {
		allLabelNames = append(allLabelNames, item.LabelNames...)
	}

	// 去重
	labelNameSet := make(map[string]bool)
	for _, name := range allLabelNames {
		labelNameSet[name] = true
	}
	uniqueLabelNames := make([]string, 0, len(labelNameSet))
	for name := range labelNameSet {
		uniqueLabelNames = append(uniqueLabelNames, name)
	}

	// 创建标签
	_, err = unsLabelService.CreateBatch(l.ctx, uniqueLabelNames)
	if err != nil {
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.label.batch.create.failed") + ": " + err.Error(),
		}, nil
	}

	// 为每个文件打标签
	var unsMapper dao.UnsNamespaceRepo
	for _, item := range req.Items {
		// 通过 alias 查询 uns
		db := dao.GetDb(l.ctx)
		uns, err := unsMapper.GetByAlias(db, item.FileAlias)
		if err != nil {
			l.Errorf("查询文件失败: %v", err)
			continue
		}
		if uns == nil {
			l.Errorf("文件不存在: %s", item.FileAlias)
			continue
		}

		// 获取标签ID映射
		labelIdMap, err := unsLabelService.CreateBatch(l.ctx, item.LabelNames)
		if err != nil {
			l.Errorf("创建标签失败: %v", err)
			continue
		}

		// 转换为 labelId -> labelName 映射
		labelIds := make([]int64, 0, len(labelIdMap))
		for _, labelId := range labelIdMap {
			labelIds = append(labelIds, labelId)
		}

		// 获取现有的标签
		var labelRefRepo dao.UnsLabelRefRepo
		var labelRepo dao.UnsLabelRepo
		existingLabels, err := labelRepo.ListByUnsId(db, uns.Id)
		if err != nil {
			l.Errorf("查询现有标签失败: %v", err)
			continue
		}

		// 合并标签
		existingLabelIdMap := make(map[int64]string)
		for _, label := range existingLabels {
			existingLabelIdMap[label.ID] = label.LabelName
		}

		// 添加新标签
		for labelName, labelId := range labelIdMap {
			if _, exists := existingLabelIdMap[labelId]; !exists {
				existingLabelIdMap[labelId] = labelName
				// 插入关联关系
				err = labelRefRepo.Insert(db, &dao.UnsLabelRef{
					UnsID:   uns.Id,
					LabelID: labelId,
				})
				if err != nil {
					l.Errorf("插入标签关联失败: %v", err)
				}
			}
		}

		// 更新 uns 的标签映射
		updatePo := &dao.UnsNamespace{
			Id:       uns.Id,
			LabelIds: existingLabelIdMap,
		}
		err = unsMapper.Update(db, updatePo)
		if err != nil {
			l.Errorf("更新文件标签失败: %v", err)
		}
	}

	return &types.ResultVO{
		Code: 200,
		Msg:  "ok",
	}, nil
}
