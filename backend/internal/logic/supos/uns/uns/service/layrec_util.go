package service

import (
	dao "backend/internal/repo/relationDB"
	"backend/share/base"
	"fmt"
	"sort"
	"strconv"
	"time"
	"unicode"
)

type saveOrUpdate struct {
	insertList []*dao.UnsNamespace
	updateList []*dao.UnsNamespace
}

func (su *saveOrUpdate) String() string {
	return fmt.Sprintf("{insert=%v, update=%v}", su.insertList, su.updateList)
}

// setLayRecAndPath 实现层级路径计算和节点分类
func setLayRecAndPath(updateTime time.Time, addFiles map[int64]*dao.UnsNamespace, dbFiles map[int64]*dao.UnsNamespace) *saveOrUpdate {
	// 合并所有节点（新增覆盖已有）
	allNodes := make(map[int64]*dao.UnsNamespace)
	for k, v := range dbFiles {
		allNodes[k] = v
	}
	for k, v := range addFiles {
		allNodes[k] = v
	}

	// 构建父子关系图
	childrenMap := make(map[int64][]*dao.UnsNamespace)
	rootNodes := make([]*dao.UnsNamespace, 0)

	for _, node := range allNodes {
		if node.ParentID == nil {
			if node.Path == "" && node.Name != "" {
				node.Path = node.Name
			}
			rootNodes = append(rootNodes, node)
		} else {
			childrenMap[*node.ParentID] = append(childrenMap[*node.ParentID], node)
		}
	}

	// 需要更新的节点集合
	nodesToInsert := make(map[int64]*dao.UnsNamespace)
	nodesToUpdate := make(map[int64]*dao.UnsNamespace)

	// 处理路径和名称
	processPathName(rootNodes, addFiles)
	for _, children := range childrenMap {
		processPathName(children, addFiles)
	}

	// 分类节点
	recorder := func(po *dao.UnsNamespace) bool {
		id := po.ID
		if _, inDB := dbFiles[id]; !inDB {
			// 新增节点
			return base.PutIfAbsent(nodesToInsert, po.ID, po)
		} else {
			// 更新节点
			po.UpdateAt = updateTime
			return base.PutIfAbsent(nodesToUpdate, po.ID, po)
		}
	}

	// 处理所有节点
	for _, node := range allNodes {
		id := node.ID
		proc := addFiles[id] != nil
		if !proc {
			if dbPo := dbFiles[id]; dbPo != nil && (node.LayRec == "" || !equalsInt64(node.ParentID, dbPo.ParentID)) {
				proc = true
			}
		}
		if proc {
			// 生成当前节点的层级路径
			generatePath(node, allNodes)
			// 收集当前节点及其所有子节点用于更新
			collectAffectedNodes(node, childrenMap, allNodes, recorder)
		}
	}

	return &saveOrUpdate{insertList: base.MapValues(nodesToInsert), updateList: base.MapValues(nodesToUpdate)}
}
func equalsInt64(a, b *int64) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// generatePath 生成单个节点的层级路径（递归向上查找）
func generatePath(node *dao.UnsNamespace, allNodes map[int64]*dao.UnsNamespace) {
	if node.ParentID == nil { // 根节点
		node.LayRec = strconv.FormatInt(node.ID, 10)
		node.Path = node.PathName
	} else {
		parent := allNodes[*node.ParentID]
		if parent == nil { // 处理异常情况
			node.LayRec = strconv.FormatInt(node.ID, 10)
			node.Path = node.Name
			return
		}

		// 递归生成父节点路径
		if parent.LayRec == "" {
			generatePath(parent, allNodes)
		}

		node.LayRec = parent.LayRec + "/" + strconv.FormatInt(node.ID, 10)
		node.Path = parent.Path + "/" + node.PathName
	}
}

// collectAffectedNodes 收集受影响的所有子节点（BFS遍历）
func collectAffectedNodes(changedNode *dao.UnsNamespace, childrenMap map[int64][]*dao.UnsNamespace, allNodes map[int64]*dao.UnsNamespace, result func(*dao.UnsNamespace) bool) {
	queue := make([]*dao.UnsNamespace, 0, 32)
	queue = append(queue, changedNode)
	result(changedNode)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// 查找所有子节点
		if children, exists := childrenMap[current.ID]; exists {
			for _, node := range children {
				if result(node) { // 避免重复处理
					queue = append(queue, node)
					// 重新生成子节点路径
					generatePath(node, allNodes)
				}
			}
		}
	}
}

// processPathName 处理同名兄弟节点的路径
func processPathName(siblings []*dao.UnsNamespace, addFiles map[int64]*dao.UnsNamespace) {
	if len(siblings) == 0 {
		return
	}

	// 按名称分组
	nameGroup := make(map[string][]*dao.UnsNamespace)
	for _, node := range siblings {
		name := escapeName(node.Name)
		nameGroup[name] = append(nameGroup[name], node)
	}

	// 对每个分组按ID排序并设置pathName
	for name, group := range nameGroup {
		if len(group) > 1 {
			sort.Slice(group, func(i, j int) bool {
				return group[i].ID < group[j].ID
			})
		}
		for i, node := range group {
			if base.MapContainsKey(addFiles, node.ID) && node.CountExistsSiblings > 0 {
				node.PathName = name + "-" + strconv.FormatInt(node.CountExistsSiblings+int64(i), 10)
			} else {
				node.PathName = name
			}
		}
	}
}

// escapeName 处理名称中的特殊字符
func escapeName(name string) string {
	cs := []rune(name)
	changed := false
	for i, c := range cs {
		if c == '-' || !isIdentifierPart(c) || c == '$' {
			changed = true
			cs[i] = '_'
		}
	}
	if changed {
		return string(cs)
	}
	return name
}
func isIdentifierPart(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'
}
