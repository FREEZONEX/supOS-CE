package base

import (
	"cmp"
	"sort"
)

func SorByDependency[T any, R cmp.Ordered](
	deps []T,
	less func(T, T) bool,
	keyFunc func(T) R,
	valueFunc func(T) R,
) (okValues, errorValues []T) {
	gf := buildReverseGraph(deps, keyFunc, valueFunc)
	lm := CalculateLevels(gf)

	if len(lm) < len(gf) {
		errorValues = make([]T, 0, len(deps)-len(lm))
		okValues = make([]T, 0, len(lm))
		kM := make(map[R]T)
		for _, d := range deps {
			kM[keyFunc(d)] = d
		}
		for k := range lm {
			v, has := kM[k]
			delete(kM, k)
			if has {
				okValues = append(okValues, v)
			}
		}
		for _, d := range kM {
			errorValues = append(errorValues, d)
		}
		if len(okValues) == 0 {
			return
		}
	} else {
		okValues = deps
	}
	arr := dgSlice{
		size: len(okValues),
		swap: func(i, j int) {
			okValues[i], okValues[j] = okValues[j], okValues[i]
		},
		less: func(i, j int) (cmp bool) {
			m, n := okValues[i], okValues[j]
			a, b := keyFunc(m), keyFunc(n)
			rs := lm[a] - lm[b]
			if rs != 0 {
				cmp = rs < 0
			} else {
				cmp = less(m, n)
			}
			return cmp
		},
	}
	sort.Sort(arr)
	return
}

type dgSlice struct {
	size int
	less func(i, j int) bool
	swap func(i, j int)
}

func (x dgSlice) Len() int           { return x.size }
func (x dgSlice) Less(i, j int) bool { return x.less(i, j) }
func (x dgSlice) Swap(i, j int)      { x.swap(i, j) }

// 构建反向依赖图
func buildReverseGraph[T any, R cmp.Ordered](
	deps []T,
	keyFunc func(T) R,
	valueFunc func(T) R,
) map[R][]R {
	graph := make(map[R][]R)
	nodes := make(map[R]struct{})

	for _, dep := range deps {
		from, to := keyFunc(dep), valueFunc(dep)
		nodes[from] = struct{}{}
		nodes[to] = struct{}{}
		graph[to] = append(graph[to], from)
	}

	// 确保所有节点都在图中
	for node := range nodes {
		if _, ok := graph[node]; !ok {
			graph[node] = nil
		}
	}
	return graph
}

func CalculateLevels[R cmp.Ordered](
	graph map[R][]R,
) map[R]int {
	levelMap := make(map[R]int)

	// 第一步：找到所有根节点（没有父节点的节点）
	// 构建子节点到父节点的反向映射
	childToParents := make(map[R][]R)
	allNodes := make(map[R]bool)

	// 收集所有节点
	for parent, children := range graph {
		allNodes[parent] = true
		for _, child := range children {
			allNodes[child] = true
			childToParents[child] = append(childToParents[child], parent)
		}
	}

	// 找到根节点（没有出现在childToParents中的节点，或者没有父节点的节点）
	var roots []R
	for node := range allNodes {
		if len(childToParents[node]) == 0 {
			roots = append(roots, node)
			levelMap[node] = 0 // 根节点层级为0
		}
	}

	// 第二步：使用BFS计算每个节点的层级
	visited := make(map[R]bool)
	queue := make([]R, 0)

	// 将根节点加入队列
	for _, root := range roots {
		queue = append(queue, root)
		visited[root] = true
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		currentLevel := levelMap[current]

		// 遍历当前节点的所有子节点
		for _, child := range graph[current] {
			if !visited[child] {
				visited[child] = true
				// 子节点的层级 = 父节点层级 + 1
				levelMap[child] = currentLevel + 1
				queue = append(queue, child)
			}
		}
	}

	return levelMap
}
