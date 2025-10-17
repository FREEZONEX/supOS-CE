package base

type ReadonlyMap[K any, V any] interface {
	Get(k K) (V, bool)
	ContainsKey(key K) bool
	Iterator() EntryIterator[K, V]
}
