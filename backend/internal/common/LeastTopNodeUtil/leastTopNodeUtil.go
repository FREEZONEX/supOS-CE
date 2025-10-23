package LeastTopNodeUtil

// TreeNode 接口定义
type TreeNode interface {
	GetID() int64
	GetParentID() *int64
}

// 泛型版本的TreeNodes接口
type TreeNodes[T TreeNode] interface {
	Size() int
	Visit(func(T))
}

// 获取最少顶层节点
func GetLeastTopNodes[T TreeNode](treeNodes TreeNodes[T]) []T {
	if treeNodes == nil || treeNodes.Size() == 0 {
		return make([]T, 0)
	}

	// 构建节点映射和子节点树
	nodeMap := make(map[int64]T)
	childrenMap := make(map[int64][]int64)

	treeNodes.Visit(func(node T) {
		nodeMap[node.GetID()] = node
		childrenMap[node.GetID()] = make([]int64, 0)
	})

	treeNodes.Visit(func(node T) {
		parentId := node.GetParentID()
		if parentId != nil && *parentId != 0 {
			if _, exists := nodeMap[*parentId]; exists {
				childrenMap[*parentId] = append(childrenMap[*parentId], node.GetID())
			}
		}
	})

	covered := make(map[int64]bool)
	result := make([]T, 0)

	treeNodes.Visit(func(node T) {
		currentId := node.GetID()
		if covered[currentId] {
			return
		}

		// 向上查找最顶层的未覆盖祖先
		topId := findTopNode(currentId, nodeMap, covered)
		if topId == nil {
			return
		}

		result = append(result, nodeMap[*topId])
		markDescendants(*topId, childrenMap, covered)
	})

	return result
}

func findTopNode[T TreeNode](startId int64, nodeMap map[int64]T, covered map[int64]bool) *int64 {
	currentId := startId
	var topId *int64
	isChainUncovered := true

	for currentId != 0 {
		if covered[currentId] {
			isChainUncovered = false
			break
		}
		topId = &currentId
		current := nodeMap[currentId]
		parentId := current.GetParentID()
		if parentId == nil || *parentId == 0 {
			break
		}
		if _, exists := nodeMap[*parentId]; !exists {
			break
		}
		currentId = *parentId
	}

	if isChainUncovered {
		return topId
	}
	return nil
}

func markDescendants(rootId int64, childrenMap map[int64][]int64, covered map[int64]bool) {
	queue := make([]int64, 0)
	queue = append(queue, rootId)

	for len(queue) > 0 {
		currentId := queue[0]
		queue = queue[1:]

		if covered[currentId] {
			continue
		}
		covered[currentId] = true

		if children, exists := childrenMap[currentId]; exists {
			queue = append(queue, children...)
		}
	}
}
