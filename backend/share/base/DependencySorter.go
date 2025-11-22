package base

import (
	"cmp"
	"sort"
)

func SorByDependency[T any, R cmp.Ordered](
	deps []T,
	keyFunc func(T) R,
	valueFunc func(T) R,
) {
	gf := BuildReverseGraph(deps, keyFunc, valueFunc)
	lm := CalculateLevels(gf)
	arr := dgSlice{
		size: len(deps),
		swap: func(i, j int) {
			deps[i], deps[j] = deps[j], deps[i]
		},
		less: func(i, j int) bool {
			a, b := keyFunc(deps[i]), keyFunc(deps[j])
			return lm[a] > lm[b]
		},
	}
	sort.Sort(arr)
}

type dgSlice struct {
	size int
	less func(i, j int) bool
	swap func(i, j int)
}

func (x dgSlice) Len() int           { return x.size }
func (x dgSlice) Less(i, j int) bool { return x.less(i, j) }
func (x dgSlice) Swap(i, j int)      { x.swap(i, j) }

// 计算层级（记忆化递归）
func CalculateLevels[R cmp.Ordered](
	graph map[R][]R,
) map[R]int {
	levels := make(map[R]int)
	for node := range graph {
		getLevel(node, graph, levels)
	}
	return levels
}

// 构建反向依赖图
func BuildReverseGraph[T any, R cmp.Ordered](
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

// 内部递归计算层级
func getLevel[R cmp.Ordered](
	node R,
	graph map[R][]R,
	levels map[R]int,
) int {
	level, ok := levels[node]
	if !ok {
		parents := graph[node]
		maxLevel := -1
		for _, parent := range parents {
			if parentLevel := getLevel[R](parent, graph, levels); parentLevel > maxLevel {
				maxLevel = parentLevel
			}
		}
		level = maxLevel + 1
		levels[node] = level
	}
	return level
}
