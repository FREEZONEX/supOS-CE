package base

func GroupByForEach[T any, G comparable](arr []T, grp func(e T) G, visit func(k G, vs []T)) {
	groupMap := make(map[G][]T, len(arr))
	for _, e := range arr {
		k := grp(e)
		vs, has := groupMap[k]
		if !has {
			vs = make([]T, 0, len(arr)/2)
		}
		vs = append(vs, e)
		groupMap[k] = vs
	}
	for k, vs := range groupMap {
		visit(k, vs)
	}
	groupMap = nil
}
func GroupBy[T any, G comparable](arr []T, grp func(e T) G) map[G][]T {
	groupMap := make(map[G][]T, len(arr))
	for _, e := range arr {
		k := grp(e)
		vs, has := groupMap[k]
		if !has {
			vs = make([]T, 0, len(arr)/2)
		}
		vs = append(vs, e)
		groupMap[k] = vs
	}
	return groupMap
}
