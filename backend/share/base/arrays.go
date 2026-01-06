package base

import (
	"cmp"
	"fmt"
	"strings"
)

func SearchInt(a []int, x int) int {
	return BinarySearch(len(a), func(i int) int {
		return a[i] - x
	})
}
func SearchLong(a []int64, x int64) int {
	return BinarySearch(len(a), func(i int) int {
		if a[i] > x {
			return 1
		} else if a[i] < x {
			return -1
		}
		return 0
	})
}
func SearchStrings(a []string, x string) int {
	return BinarySearch(len(a), func(i int) int {
		if a[i] > x {
			return 1
		} else if a[i] < x {
			return -1
		}
		return 0
	})
}

// BinarySearch 二分查找算法
func BinarySearch(n int, f func(int) int) int {
	low, high := 0, n-1

	for low <= high {
		var mid = int(uint(low+high) >> 1)
		var cmp = f(mid)

		if cmp < 0 {
			low = mid + 1
		} else if cmp > 0 {
			high = mid - 1
		} else {
			return mid // key found
		}
	}
	return -(low + 1) // key not found. 返回 负的插入位置
}
func BinarySearchCmp[T cmp.Ordered](arr []T, key T) int {
	low, high := 0, len(arr)-1

	for low <= high {
		var mid = int(uint(low+high) >> 1)
		var cmRs = cmp.Compare(arr[mid], key)

		if cmRs < 0 {
			low = mid + 1
		} else if cmRs > 0 {
			high = mid - 1
		} else {
			return mid // key found
		}
	}
	return -(low + 1) // key not found. 返回 负的插入位置
}

// BinarySearchArray in java style
func BinarySearchArray[T any](arr []T, key T, comparator func(a, b T) int) int {
	low, high := 0, len(arr)-1

	for low <= high {
		var mid = int(uint(low+high) >> 1)
		var cmRs = comparator(arr[mid], key)

		if cmRs < 0 {
			low = mid + 1
		} else if cmRs > 0 {
			high = mid - 1
		} else {
			return mid // key found
		}
	}
	return -(low + 1) // key not found. 返回 负的插入位置
}

// BinarySearchLowHigh 二分查找算法
func BinarySearchLowHigh(n int, f func(int) int) (min, max int) {
	low, high := 0, n-1

	for low < high {
		var mid = int(uint(low+high+1) >> 1)
		var cmRs = f(mid)

		if cmRs < 0 {
			low = mid
		} else if cmRs > 0 {
			high = mid - 1
		} else {
			return mid, mid // key found
		}
	}
	min = low
	max = low + 1
	return min, max
}
func Equals[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, e := range a {
		if b[i] != e {
			return false
		}
	}
	return true
}
func EqualsF[T comparable](a, b []T, compare func(a, b T) bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i, e := range a {
		if !compare(b[i], e) {
			return false
		}
	}
	return true
}
func ToString[T comparable](arr []T) string {
	var str = strings.Builder{}
	str.Grow(4 + len(arr)*10)
	str.WriteByte('[')
	for i, e := range arr {
		if i > 0 {
			str.WriteByte(',')
		}
		str.WriteString(fmt.Sprint(e))
	}
	str.WriteByte(']')
	return str.String()
}

// Partition 将切片按指定大小分割成多个子切片
func Partition[E any](arr []E, size int) [][]E {
	length := len(arr)
	if length == 0 || size < 1 {
		return [][]E{}
	}

	// 计算需要分割的子切片数量
	result := make([][]E, 0, (length+size-1)/size)

	for i := 0; i < length; i += size {
		end := i + size
		if end > length {
			end = length
		}
		// 将原始切片的一部分添加到结果中
		result = append(result, arr[i:end])
	}

	return result
}

type LongSlice []int64

func (x LongSlice) Len() int           { return len(x) }
func (x LongSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x LongSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
