package bits

type littleEndian byte

var LittleEndian ByteOrder = littleEndian(1)

func (littleEndian) String() string { return "LittleEndian" }

func (littleEndian) GoString() string { return "bits.LittleEndian" }

func (o littleEndian) PutShort(bb []byte, x int16) {
	putShortL(bb, x)
}
func (o littleEndian) GetShort(bb []byte) int16 {
	return getShortL(bb)
}
func (o littleEndian) PutInt(bb []byte, x int) {
	putIntL(bb, x)
}
func (o littleEndian) GetInt(bb []byte) int {
	return getIntL(bb)
}
func (o littleEndian) PutLong(bb []byte, x int64) {
	putLongL(bb, x)
}
func (o littleEndian) GetLong(bb []byte) int64 {
	return getLongL(bb)
}
func (littleEndian) Uint16(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[0]) | uint16(b[1])<<8
}

func (littleEndian) PutUint16(b []byte, v uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}
func (littleEndian) Uint32(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func (littleEndian) PutUint32(b []byte, v uint32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}
func (littleEndian) Uint64(b []byte) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func (littleEndian) PutUint64(b []byte, v uint64) {
	_ = b[7] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

// -- short
func putShortL(bb []byte, x int16) {
	bb[1] = byte(x >> 8)
	bb[0] = byte(x)
}
func getShortL(bb []byte) int16 {
	return makeShort(bb[1],
		bb[0])
}

// -- int
func putIntL(bb []byte, x int) {
	bb[3] = byte(x >> 24)
	bb[2] = byte(x >> 16)
	bb[1] = byte(x >> 8)
	bb[0] = byte(x)
}
func getIntL(bb []byte) int {
	return makeInt(bb[3],
		bb[2],
		bb[1],
		bb[0])
}

// -- long
func putLongL(bb []byte, x int64) {
	bb[7] = byte(x >> 56)
	bb[6] = byte(x >> 48)
	bb[5] = byte(x >> 40)
	bb[4] = byte(x >> 32)
	bb[3] = byte(x >> 24)
	bb[2] = byte(x >> 16)
	bb[1] = byte(x >> 8)
	bb[0] = byte(x)
}

func getLongL(bb []byte) int64 {
	return makeLong(bb[7],
		bb[6],
		bb[5],
		bb[4],
		bb[3],
		bb[2],
		bb[1],
		bb[0])
}
