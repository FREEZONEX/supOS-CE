package service

import (
	"backend/internal/common/dto"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"math"
	"strconv"
	"strings"

	"gitee.com/unitedrhino/share/stores"
)

// SearchTree 旧的整树查询
func (l *UnsQueryService) SearchTree(ctx context.Context, req *types.SearchTreeReq) (resp *types.SearchTreeResp, err error) {
	switch req.SearchType {
	case 1:
		resp, err = l.searchTree(ctx, req)
	case 2:
		resp, err = l.searchByTag(ctx, req)
	case 3:
		resp, err = l.searchByTemplate(ctx, req)
	}
	return
}
func (l *UnsQueryService) searchTree(ctx context.Context, req *types.SearchTreeReq) (resp *types.SearchTreeResp, err error) {
	db := dao.GetDb(ctx)
	pid, pType := req.ParentId, req.PathType
	query := &dto.UnsSearchCondition{Keyword: req.Keyword, SearchType: req.SearchType}
	if pid > 0 {
		query.ParentId = &pid
	}
	if pType >= 0 {
		query.PathType = &pType
	}
	pageRs, er := l.unsMapper.PageListByConditions(db, query,
		&stores.PageInfo{Page: 1, Size: math.MaxInt64, Orders: []stores.OrderBy{{Field: "id", Sort: stores.OrderAsc}}},
	)
	resp = &types.SearchTreeResp{}
	resp.Code = 500
	if er != nil || len(pageRs.Data) == 0 {
		err = er
		if er != nil {
			resp.Msg = er.Error()
		} else {
			resp.Code = 200
			resp.Msg = "NoData"
		}
		return
	}
	list := pageRs.Data
	layRecList := base.NewEmptySet[string](32)
	for _, uns := range list {
		layRec := uns.LayRec
		lastSlash := strings.LastIndex(layRec, "/")
		if lastSlash < 0 {
			layRecList.Add(layRec)
		} else {
			layRecSub := layRec[:lastSlash]
			if strings.Contains(layRecSub, "/") {
				layRecList.AddAll(strings.Split(layRecSub, "/"))
			} else {
				layRecList.Add(layRecSub)
			}
		}
	}
	layRecLongSet := base.FilterAndMap[string, int64](layRecList.Values(), func(entry string) (id int64, ok bool) {
		var numErr error
		id, numErr = strconv.ParseInt(entry, 10, 64)
		ok = numErr == nil && id > 0
		return
	})
	allNamespaces, er := l.unsMapper.SelectByIds(db, layRecLongSet)
	if er != nil {
		resp.Msg = er.Error()
		err = er
		return
	}
	allNamespaces = append(allNamespaces, list...)
	treeResults, er := l.getTopicTreeResults(allNamespaces, list, req.ShowRec)
	if er != nil {
		resp.Msg = er.Error()
		err = er
		return
	}
	resp.Data = treeResults
	resp.Code = 200
	resp.Msg = "ok"
	return
}
func (l *UnsQueryService) searchByTag(ctx context.Context, req *types.SearchTreeReq) (resp *types.SearchTreeResp, err error) {
	db := dao.GetDb(ctx)
	var list []*dao.UnsNamespace
	list, err = l.unsMapper.ListLabeledUnsByKeyword(db, req.Keyword)
	if err != nil {
		resp = &types.SearchTreeResp{}
		resp.Code = 500
		return
	}
	resp = &types.SearchTreeResp{Data: l.uns2treeList(req, list)}
	resp.Code = 200
	resp.Msg = "ok"
	return
}
func (l *UnsQueryService) searchByTemplate(ctx context.Context, req *types.SearchTreeReq) (resp *types.SearchTreeResp, err error) {
	db := dao.GetDb(ctx)
	resp = &types.SearchTreeResp{}
	resp.Code = 500
	var list []*dao.UnsNamespace
	list, err = l.unsMapper.ListInTemplate(db, req.Keyword)
	if err != nil {
		resp = &types.SearchTreeResp{}
		resp.Code = 500
		return
	}
	resp = &types.SearchTreeResp{Data: l.uns2treeList(req, list)}
	resp.Code = 200
	resp.Msg = "ok"
	return
}

func (l *UnsQueryService) uns2treeList(req *types.SearchTreeReq, list []*dao.UnsNamespace) (treeResults []*types.TopicTreeResult) {
	treeResults = base.FilterAndMap[*dao.UnsNamespace, *types.TopicTreeResult](list, func(uns *dao.UnsNamespace) (v *types.TopicTreeResult, match bool) {
		v = &types.TopicTreeResult{
			Id:             strconv.FormatInt(uns.Id, 10),
			PathType:       uns.PathType,
			Protocol:       uns.ProtocolType,
			ParentDataType: uns.ParentDataType,
			Path:           uns.Path,
			Name:           PathUtil.GetName(uns.Path),
		}
		if l.mountService != nil {
			v.Mount = l.mountService.ParseMountDetail(uns, true)
		}
		return v, true
	})
	return
}

// GetTopicTreeResults 获取主题树结果（带错误处理）
func (l *UnsQueryService) getTopicTreeResults(all []*dao.UnsNamespace, list []*dao.UnsNamespace, showRec bool) ([]*types.TopicTreeResult, error) {
	if len(list) == 0 {
		return []*types.TopicTreeResult{}, nil
	}

	nodeMap := make(map[string]*types.TopicTreeResult)

	if len(all) == 0 {
		all = list
	}

	// 构建节点映射
	for _, po := range all {
		if po == nil {
			continue
		}

		path := po.Path
		if path == "" {
			continue
		}

		rs := UnsConverter.Dto2TreeResult(po)

		//if po.PathType == 2 {
		//	info, exists := topicLastMessages[path]
		//	if exists && info != nil {
		//		rs.Value = info.MessageCount
		//		if showRec {
		//			rs.LastUpdateTime = info.LastUpdateTime
		//		}
		//	}
		//}

		// 处理字段
		if po.Fields != nil {
			fields := po.Fields
			for _, f := range fields {
				f.Index = nil
			}
			rs.Fields = fields
		}
		if mountService := l.mountService; mountService != nil {
			rs.Mount = mountService.ParseMountDetail(po, true)
		}
		nodeMap[path] = rs
	}

	return buildTreeStructure(nodeMap, list), nil
}

// buildTreeStructure 构建树形结构
func buildTreeStructure(nodeMap map[string]*types.TopicTreeResult, list []*dao.UnsNamespace) []*types.TopicTreeResult {
	rootNodes := make(map[string]*types.TopicTreeResult)
	childrenMap := make(map[string]map[string]bool)

	for _, po := range list {
		if po == nil || po.Path == "" {
			continue
		}

		path := po.Path
		currentNode, exists := nodeMap[path]
		if !exists {
			continue
		}

		parentPath := PathUtil.SubParentPath(path)

		// 当前节点就是根节点
		if parentPath == "" {
			ensureChildMap(childrenMap, parentPath)
			if _, has := childrenMap[parentPath][path]; !has {
				childrenMap[parentPath][path] = true
				rootNodes[path] = currentNode
			}
			continue
		}

		// 构建从当前节点到根节点的路径
		addRoot, tempParentNode := buildPathToRoot(nodeMap, childrenMap, currentNode, path)

		if addRoot && tempParentNode != nil {
			ensureChildMap(childrenMap, parentPath)
			childMap := childrenMap[parentPath]
			if _, has := childMap[tempParentNode.Path]; !has {
				rootNodes[tempParentNode.Path] = tempParentNode
			}
		}
	}

	// set countChildren
	for _, node := range rootNodes {
		node.GetCountChildren()
	}
	result := base.MapValues(rootNodes)
	types.SortUnsList(result)
	return result
}

// ensureChildMap 确保子节点映射存在
func ensureChildMap(childrenMap map[string]map[string]bool, key string) {
	if _, exists := childrenMap[key]; !exists {
		childrenMap[key] = make(map[string]bool)
	}
}

// buildPathToRoot 构建到根节点的路径
func buildPathToRoot(nodeMap map[string]*types.TopicTreeResult, childrenMap map[string]map[string]bool, currentNode *types.TopicTreeResult, path string) (addRoot bool, tempParentNode *types.TopicTreeResult) {
	tempCurrentNode := currentNode
	parentPath := PathUtil.SubParentPath(path)
	addRoot = true
	for parentPath != "" {
		var exists bool
		tempParentNode, exists = nodeMap[parentPath]
		if !exists {
			name := PathUtil.GetName(parentPath)
			tempParentNode = &types.TopicTreeResult{Name: name, Path: parentPath, PathType: 0}
			nodeMap[parentPath] = tempParentNode
		}

		ensureChildMap(childrenMap, parentPath)
		childMap := childrenMap[parentPath]
		if _, has := childMap[tempCurrentNode.Path]; !has {
			childMap[tempCurrentNode.Path] = true
			tempParentNode.AddChild(tempCurrentNode)
		} else {
			addRoot = false
			break
		}

		parentPath = PathUtil.SubParentPath(parentPath)
		tempCurrentNode = tempParentNode
	}
	return addRoot, tempParentNode
}
