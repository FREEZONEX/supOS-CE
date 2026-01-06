package base

import (
	"cmp"
	"reflect"
	"unsafe"
)

type SortedString string

func (s *SortedString) Size() int {
	var str = string(*s)
	if len(str) < 2 {
		return 0
	}
	return int(str[1])
}
func (s *SortedString) Entries() [][2]string {
	var str = s2b(string(*s))
	if len(str) < 2 {
		return nil
	}
	var size = int(str[1])
	if size == 0 {
		return nil
	}
	rs := make([][2]string, size)
	INDEX_SIZE := (size + 1) << 1
	for index := 0; index < size; index++ {
		i := (index + 1) << 1
		kStart := INDEX_SIZE
		if index > 0 {
			kStart += int(str[i-2] + str[i-1])
		}
		vStart := INDEX_SIZE + int(str[i])
		vEnd := int(str[i+1])
		k := b2s(str[kStart:vStart])
		v := b2s(str[vStart : vStart+vEnd])
		rs[index] = [2]string{k, v}
	}
	return rs
}
func (s *SortedString) Get(k string) (v string, has bool) {
	var str = string(*s)
	size := int(str[1])
	val, i := s.get(k, size)
	return val, i >= 0
}
func (s *SortedString) Iterator() EntryIterator[string, string] {
	return &ssEntryIterator{array: s.Entries()}
}

type ssEntryIterator struct {
	array [][2]string
	i     int
}

func (a *ssEntryIterator) HasNext() bool {
	return a.i < len(a.array)
}
func (a *ssEntryIterator) Next() (k, v string) {
	rs := a.array[a.i]
	a.i++
	return rs[0], rs[1]
}
func (s *SortedString) ContainsKey(k string) bool {
	var str = string(*s)
	size := int(str[1])
	_, i := s.get(k, size)
	return i >= 0
}

//	func (s *SortedString) Remove(k string) bool {
//		var str = s2b(string(*s))
//		size := int(str[1])
//		v, index := s.get(k, size)
//		if index >= 0 {
//			i := (index + 1) << 1
//			INDEX_SIZE := (size + 1) << 1
//			kStart := INDEX_SIZE
//			if index > 0 {
//				kStart += int(str[i-2] + str[i-1])
//			}
//			//vStart := INDEX_SIZE + int(str[i])
//			//vEnd := int(str[i+1])
//			//k = b2s(str[kStart:vStart])
//			//v = b2s(str[vStart : vStart+vEnd])
//			for j, m := 0, len(k)+len(v); j < m; j++ {
//				str[kStart+j] = str[kStart+j+m]
//			}
//			return true
//		}
//		return false
//	}
func (s *SortedString) Put(kStr, vStr string) {
	var str = s2b(string(*s))
	var size = 0
	if len(str) > 1 {
		size = int(str[1])
	}
	vOld, index := s.get(kStr, size)
	k, v := s2b(kStr), s2b(vStr)
	if index >= 0 { // key found
		sub := len(vOld) - len(v)
		INDEX_SIZE := (size + 1) << 1
		indexPos := (index + 1) << 1
		if sub == 0 && vStr != vOld { //原地修改，不用更新索引
			i := (index + 1) << 1
			vStart := INDEX_SIZE + int(str[i])
			copy(str[vStart:], v)
			*s = SortedString(str)
		} else if sub < 0 { //值变大，原值及右侧的数据右移，并更新索引
			for dif := sub; dif < 0; dif++ {
				str = append(str, 0)
			}
			kStart := 0
			if index > 0 {
				kStart = int(str[indexPos-2] + str[indexPos-1])
			}
			dataInsertPos := INDEX_SIZE + kStart + len(k)
			// 数据右移 -sub 位
			mvoff := -sub
			for i := len(str) - 1; i >= dataInsertPos+mvoff; i-- {
				str[i] = str[i-mvoff]
			}
			copy(str[dataInsertPos:], v) //写值
			// 其后的索引值增加偏移量
			addOffset := i2b(-sub)
			for i := index + 1; i < size; i++ {
				kStartPos := (i + 1) << 1
				str[kStartPos] += addOffset
			}
			str[indexPos+1] = i2b(len(v)) //更新 value.length
			*s = SortedString(str)
		} else if sub > 0 { //值变小，原值右侧的数据左移，并更新索引
			vEnd := INDEX_SIZE + int(str[indexPos]) + len(v)
			// 原值右侧的数据左移
			for i, Max := vEnd, len(str)-sub; i < Max; i++ {
				//sd := string([]byte{str[i], str[i+sub]})
				//fmt.Printf("change[%d] %v->%v\n", i, sd[:1], sd[1:])
				str[i] = str[i+sub]
			}
			str[indexPos+1] = i2b(len(v)) //更新 value.length
			// 其后的索引值减少偏移量
			subOffset := i2b(sub)
			for i := index + 1; i < size; i++ {
				kStartPos := (i + 1) << 1
				str[kStartPos] -= subOffset
			}
			kStart := 0
			if index > 0 {
				kStart = int(str[indexPos-2] + str[indexPos-1])
			}
			dataInsertPos := INDEX_SIZE + kStart + len(k)
			copy(str[dataInsertPos:], v) //写值
			*s = SortedString(str[:len(str)-sub])
		}
	} else { // key not found
		newStr := appendEntry(str, k, v, index, size)
		*s = SortedString(b2s(newStr))
	}
}

func (s *SortedString) get(k string, size int) (v string, i int) {
	var str = s2b(string(*s))
	midK, midV := "", ""
	i = BinarySearch(size, func(mid int) int {
		midK, midV = getByIndex(str, mid, size)
		return cmp.Compare(midK, k)
	})
	if i >= 0 {
		v = midV
	}
	return v, i
}
func (s *SortedString) String() string {
	return string(*s)
}

func getByIndex(str []byte, index, size int) (k, v string) {
	// 先查索引
	i := (index + 1) << 1
	INDEX_SIZE := (size + 1) << 1
	kStart := INDEX_SIZE
	if index > 0 {
		kStart += int(str[i-2] + str[i-1])
	}
	vStart := INDEX_SIZE + int(str[i])
	vEnd := int(str[i+1])
	// 再定位数据
	k = b2s(str[kStart:vStart])
	v = b2s(str[vStart : vStart+vEnd])
	return k, v
}
func appendEntry(str, k, v []byte, index int, size int) []byte {
	index = -index - 1
	if size > 0 && index >= size { // 非空且末尾追加
		prevend := 0
		var INDEX_SIZE = (size + 1) << 1
		if size > 0 {
			prevend = int(str[INDEX_SIZE-2] + str[INDEX_SIZE-1])
		}
		dataseg := str[INDEX_SIZE:]
		indexseg := str[:INDEX_SIZE]
		str = append(dataseg, k...)
		str = append(str, v...)
		indexseg = append(indexseg, i2b(prevend+len(k)), i2b(len(v)))
		str = append(indexseg, str...)
		str[1] = i2b(size + 1) // 更新size
		return str
	} else if size == 0 { //空的情况
		str = make([]byte, 2, 255)
		str[0] = '$'
		str[1] = 1 // 初始size
		lenK, lenV := i2b(len(k)), i2b(len(v))
		str = append(str, lenK, lenV)
		str = append(str, k...)
		str = append(str, v...)
		return str
	}
	// sortedstring 如 {"city":"hangzhou","am","09"},存储的[]byte形如："[$][2] [2][2][8][8]am09cityHangzhou"
	var INDEX_SIZE = (size + 1) << 1
	str = append(str, 0, 0)
	str = append(str, k...)
	str = append(str, v...)
	indexInsertPos := (index + 1) << 1
	kStart := 0
	if index > 0 {
		kStart = int(str[indexInsertPos-2] + str[indexInsertPos-1])
	}
	dataInsertPos := INDEX_SIZE + kStart
	// 插入点右边的数据右移(2+ len(k) + len(v) )位
	mvoff := 2 + len(k) + len(v)
	for i := len(str) - 1; i >= dataInsertPos+mvoff; i-- {
		str[i] = str[i-mvoff]
	}
	// 写数据
	copy(str[2+dataInsertPos:], k)
	copy(str[2+dataInsertPos+len(k):], v)
	// 插入点左边的右移2位
	for i, Min := dataInsertPos+1, INDEX_SIZE+2; i >= Min; i-- {
		str[i] = str[i-2]
	}
	// 其后的索引值增加偏移量
	addOffset := i2b(len(k) + len(v))
	for i := index; i < size; i++ { //这里还没移动索引，所以初始值从 index 开始
		kStartPos := (i + 1) << 1
		str[kStartPos] += addOffset
	}
	//索引右边移动两位
	for i := size - 1; i >= index; i-- {
		b2 := (i + 1) << 1
		str[b2+3] = str[b2+1]
		str[b2+2] = str[b2]
	}

	// 写索引
	str[indexInsertPos] = i2b(kStart + len(k))
	str[indexInsertPos+1] = i2b(len(v))

	str[1] = i2b(size + 1) // 更新size
	return str
}
func i2b(n int) byte {
	if n > 255 {
		panic("越界!")
	}
	return byte(n)
}
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}
func InsertTo[T any](a []T, k T, insertPos int) []T {
	if insertPos >= 0 {
		return a
	}
	size := len(a)
	insertPos = -insertPos - 1
	a = append(a, k)
	if insertPos < size && insertPos >= 0 {
		for i := len(a) - 1; i > insertPos; i-- {
			a[i] = a[i-1]
		}
		a[insertPos] = k
	}
	return a
}
