package service

import (
	"backend/internal/common/utils/PathUtil"
	dao "backend/internal/repo/relationDB"
	"backend/share/base"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/zeromicro/go-zero/core/logx"
)

type saveOrUpdate struct {
	insertList map[int64]*dao.UnsNamespace
	updateList map[int64]*dao.UnsNamespace
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
		if node.ParentId == nil {
			if node.Path == "" && node.Name != "" {
				node.Path = node.Name
			}
			rootNodes = append(rootNodes, node)
		} else {
			childrenMap[*node.ParentId] = append(childrenMap[*node.ParentId], node)
		}
	}

	// 需要更新的节点集合
	nodesToInsert := make(map[int64]*dao.UnsNamespace)
	nodesToUpdate := make(map[int64]*dao.UnsNamespace)

	// 处理路径和名称
	processPathName(rootNodes, addFiles, dbFiles)
	for _, children := range childrenMap {
		processPathName(children, addFiles, dbFiles)
	}

	// 分类节点
	recorder := func(po *dao.UnsNamespace) bool {
		id := po.Id
		if _, inDB := dbFiles[id]; !inDB {
			// 新增节点
			return base.PutIfAbsent(nodesToInsert, po.Id, po)
		} else {
			// 更新节点
			po.UpdateAt = updateTime
			return base.PutIfAbsent(nodesToUpdate, po.Id, po)
		}
	}
	list := base.MapValues(allNodes)
	list, erList := base.SorByDependency(list, func(a, b *dao.UnsNamespace) bool {
		return a.Id < b.Id
	}, func(t *dao.UnsNamespace) int64 {
		return t.Id
	}, func(t *dao.UnsNamespace) int64 {
		return base.P2vWithDefault(t.ParentId, -1)
	})
	logx.Debug("setLayRecAndPath: list=", list)
	if len(erList) > 0 {
		logx.Errorf("setLayRecAndPath: 存在循环依赖: %+v", base.Map[*dao.UnsNamespace, int64](erList, func(e *dao.UnsNamespace) int64 {
			return e.Id
		}))
	}
	// 处理所有节点
	for _, node := range list {
		id := node.Id
		proc := addFiles[id] != nil
		if !proc {
			if dbPo := dbFiles[id]; dbPo != nil && (node.LayRec == "" || !equalsInt64(node.ParentId, dbPo.ParentId)) {
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
	return &saveOrUpdate{insertList: nodesToInsert, updateList: nodesToUpdate}
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
	if node.ParentId == nil { // 根节点
		node.LayRec = strconv.FormatInt(node.Id, 10)
		node.Path = node.PathName
	} else {
		parent := allNodes[*node.ParentId]
		if parent == nil { // 处理异常情况
			node.LayRec = strconv.FormatInt(node.Id, 10)
			node.Path = node.Name
			return
		}

		// 递归生成父节点路径
		if parent.LayRec == "" {
			generatePath(parent, allNodes)
		}

		node.LayRec = parent.LayRec + "/" + strconv.FormatInt(node.Id, 10)
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
		if children, exists := childrenMap[current.Id]; exists {
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
func processPathName(siblings []*dao.UnsNamespace, addFiles map[int64]*dao.UnsNamespace, dbFiles map[int64]*dao.UnsNamespace) {
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
				return group[i].Id < group[j].Id
			})
		}
		for i, node := range group {
			dbPo := dbFiles[node.Id]
			if dbPo != nil && dbPo.Name == node.Name {
				node.PathName = PathUtil.GetName(dbPo.Path)
			}
			if node.LayRec == "" && base.MapContainsKey(addFiles, node.Id) {
				xp := strings.LastIndex(name, "-")
				if xp > 0 && xp < len(name)-1 && unicode.IsDigit(rune(name[xp+1])) {
					name = name[:xp+1] + "0" + name[xp+1:]
				}
				index := int64(i) + node.CountExistsSiblings
				if index > 0 {
					node.PathName = name + "-" + strconv.FormatInt(index, 10)
				} else {
					node.PathName = name
				}
			}
			if node.PathName == "" {
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
		if (c != '-' && !isIdentifierPart(c)) || c == '$' {
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

func parseChildrenCount(layRecList []*dao.LayRecCc) map[int64][]int {
	// 创建映射存储节点ID到其直接子节点列表
	childrenMap := make(map[int64][]int64)
	// 创建映射存储节点ID到其直接子节点数量
	directChildrenCountMap := make(map[int64][]int)
	// 存储所有节点ID
	allNodes := make(map[int64]bool)

	// 解析每个LayRecCc对象，构建树结构
	for _, layRecCc := range layRecList {
		pathParts := strings.Split(layRecCc.LayRec, "/")
		nodeId, _ := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
		allNodes[nodeId] = true

		cc := make([]int, 8)
		parseTypeCount(layRecCc.CountChildren, cc)
		directChildrenCountMap[nodeId] = cc

		// 如果不是根节点，添加到父节点的子节点列表
		if len(pathParts) > 1 {
			parentId, _ := strconv.ParseInt(pathParts[len(pathParts)-2], 10, 64)
			if _, exists := childrenMap[parentId]; !exists {
				childrenMap[parentId] = make([]int64, 0)
			}
			childrenMap[parentId] = append(childrenMap[parentId], nodeId)
		}
	}

	// 创建结果映射
	result := make(map[int64][]int)

	// 找到所有根节点（没有父节点的节点）
	for nodeId := range allNodes {
		// 从根节点开始计算子孙节点总数
		calculateDescendants(nodeId, childrenMap, directChildrenCountMap, result)
	}

	return result
}

func calculateDescendants(nodeId int64,
	childrenMap map[int64][]int64,
	directChildrenCountMap map[int64][]int,
	result map[int64][]int) []int {

	// 如果已经计算过，直接返回
	if rs, exists := result[nodeId]; exists {
		return rs
	}

	totalDescendants := directChildrenCountMap[nodeId]
	if totalDescendants == nil {
		totalDescendants = make([]int, 8)
	}

	// 递归计算所有子节点的子孙节点数
	if vs, exists := childrenMap[nodeId]; exists {
		for _, childId := range vs {
			if childId != nodeId {
				rs := calculateDescendants(childId, childrenMap, directChildrenCountMap, result)
				for i := 0; i < len(totalDescendants); i++ {
					totalDescendants[i] += rs[i]
				}
			}
		}
	}

	result[nodeId] = totalDescendants
	return totalDescendants
}

func parseTypeCount(countChildren string, rs []int) {
	if countChildren != "" {
		segments := strings.Split(countChildren, ",")
		for _, kv := range segments {
			sp := strings.Index(kv, ":")
			if sp > 0 {
				pathType, _ := strconv.Atoi(kv[:sp])
				count, _ := strconv.Atoi(kv[sp+1:])
				if pathType < len(rs) {
					rs[pathType] = count
				}
			}
		}
	}
}
