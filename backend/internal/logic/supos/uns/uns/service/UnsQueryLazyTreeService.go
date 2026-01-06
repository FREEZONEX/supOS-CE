package service

import (
	"backend/internal/common/constants"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
)

// LazyTree 懒加载的树查询
func (l *UnsQueryService) LazyTree(ctx context.Context, params *types.UnsTreeCondition) (resp *types.UnsTreePageResp, err error) {
	if params.PageNo < 1 {
		params.PageNo = 1
	}
	if params.PageSize < 1 {
		params.PageSize = constants.DefaultPageSize
	} else if params.PageSize > 1000 {
		params.PageSize = constants.MaxPageSize
	}
	pageNo := params.PageNo
	pageSize := params.PageSize
	query := &dao.UnsTreeNextLevelQuery{UnsTreeCondition: *params}

	db := dao.GetDb(ctx)
	ctx = dao.SetDb(ctx, db)
	var parentId = params.ParentId
	if parentId != nil && *parentId != 0 {
		var parent *dao.UnsNamespace
		parent, err = l.unsMapper.SelectById(db, *parentId)
		if err != nil || parent == nil {
			resp = emptyPage(params)
			if err != nil {
				resp.Code = 500
			}
			return
		}
		query.LayRecPrev = parent.LayRec
	}
	if params.SearchType == 1 &&
		params.Keyword == "" &&
		params.PathType == nil &&
		params.DataType == nil &&
		params.SubscribeEnable == nil {
		// 不考虑parentId，无其他条件的简单搜索
		return l.simpleTree(ctx, parentId, query.LayRecPrev, pageNo, pageSize)
	}

	total := int64(0)
	var treeResultList []*types.TopicTreeResult
	var list []*dao.TreeNodeUns
	list, err = l.unsMapper.NextLevelPagedQueryList(db, query, &total)
	if err != nil {
		resp = emptyPage(params)
		resp.Code = 500
		return
	}
	if len(list) > 0 {
		rsTypes := make([]int, 8)
		treeResultList = make([]*types.TopicTreeResult, 0, len(list))

		for _, po := range list {
			result := UnsConverter.Dto2TreeResult(po)

			var folderCount, fileCount int
			if po.PathType == constants.PathTypeDir {
				countChildren := po.CountChildren
				parseTypeCount(countChildren, rsTypes)
				folderCount = rsTypes[constants.PathTypeDir]
				fileCount = rsTypes[constants.PathTypeFile]
			}

			result.CountChildren = &fileCount
			result.HasChildren = folderCount+fileCount > 0
			treeResultList = append(treeResultList, result)
		}
	}

	resp = &types.UnsTreePageResp{PageResultDTO: types.PageResultDTO{
		PageNo:   int64(pageNo),
		PageSize: int64(pageSize),
		Total:    total,
		Code:     200,
		Msg:      "Normal DB Search",
	}, Data: treeResultList}
	return
}
func (l *UnsQueryService) simpleTree(ctx context.Context, parentId *int64, layRecPrev string, pageNo, pageSize int) (resp *types.UnsTreePageResp, err error) {
	db := dao.GetDb(ctx)
	countChildrenList, er := l.unsMapper.ListCountChildren(db, layRecPrev)
	if er != nil {
		return nil, er
	}
	ccMap := parseChildrenCount(countChildrenList)
	total := int64(0)
	pageRs, er := l.unsMapper.ParentIdPagedQueryList(db, parentId, pageNo, pageSize, &total)
	if er != nil {
		return nil, er
	}
	treeResultList := base.Map[*dao.UnsNamespace, *types.TopicTreeResult](pageRs, func(e *dao.UnsNamespace) *types.TopicTreeResult {
		result := UnsConverter.Dto2TreeResult(e)
		folderCount, fileCount := 0, 0
		if e.PathType == constants.PathTypeDir {
			if rsTypes := ccMap[e.Id]; len(rsTypes) > 2 {
				folderCount = rsTypes[constants.PathTypeDir]
				fileCount = rsTypes[constants.PathTypeFile]
			}
		}
		result.CountChildren = &fileCount
		result.HasChildren = folderCount+fileCount > 0
		return result
	})
	resp = &types.UnsTreePageResp{
		PageResultDTO: types.PageResultDTO{
			PageNo:   int64(pageNo),
			PageSize: int64(pageSize),
			Total:    total,
			Code:     200,
			Msg:      "Simple DB Search",
		}, Data: treeResultList}
	return
}
func emptyPage(params *types.UnsTreeCondition) *types.UnsTreePageResp {
	return &types.UnsTreePageResp{
		PageResultDTO: types.PageResultDTO{
			Code:     200,
			PageNo:   int64(params.PageNo),
			PageSize: int64(params.PageSize),
		},
	}
}
