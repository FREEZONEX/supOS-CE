package bits

type bigEndian byte

var BigEndian ByteOrder = bigEndian(0)

func (bigEndian) String() string { return "BigEndian" }

func (bigEndian) GoString() string { return "bits.BigEndian" }

func (o bigEndian) PutShort(bb []byte, x int16) {
	putShortB(bb, x)
}
func (o bigEndian) GetShort(bb []byte) int16 {
	return getShortB(bb)
}
func (o bigEndian) PutInt(bb []byte, x int) {
	putIntB(bb, x)
}
func (o bigEndian) GetInt(bb []byte) int {
	return getIntB(bb)
}
func (o bigEndian) PutLong(bb []byte, x int64) {
	putLongB(bb, x)
}
func (o bigEndian) GetLong(bb []byte) int64 {
	return getLongB(bb)
}
func (bigEndian) Uint16(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[1]) | uint16(b[0])<<8
}

func (bigEndian) PutUint16(b []byte, v uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 8)
	b[1] = byte(v)
}

func (bigEndian) AppendUint16(b []byte, v uint16) []byte {
	return append(b,
		byte(v>>8),
		byte(v),
	)
}

func (bigEndian) Uint32(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

func (bigEndian) PutUint32(b []byte, v uint32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}

func (bigEndian) AppendUint32(b []byte, v uint32) []byte {
	return append(b,
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

func (bigEndian) Uint64(b []byte) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

func (bigEndian) PutUint64(b []byte, v uint64) {
	_ = b[7] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
}

func (bigEndian) AppendUint64(b []byte, v uint64) []byte {
	return append(b,
		byte(v>>56),
		byte(v>>48),
		byte(v>>40),
		byte(v>>32),
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

// -- short
func putShortB(bb []byte, x int16) {
	bb[0] = byte(x >> 8)
	bb[1] = byte(x)
}
func getShortB(bb []byte) int16 {
	return makeShort(bb[0],
		bb[1])
}

// -- int
func putIntB(bb []byte, x int) {
	bb[0] = byte(x >> 24)
	bb[1] = byte(x >> 16)
	bb[2] = byte(x >> 8)
	bb[3] = byte(x)
}

func getIntB(bb []byte) int {
	return makeInt(bb[0],
		bb[1],
		bb[2],
		bb[3])
}

// -- long
func putLongB(bb []byte, x int64) {
	bb[0] = byte(x >> 56)
	bb[1] = byte(x >> 48)
	bb[2] = byte(x >> 40)
	bb[3] = byte(x >> 32)
	bb[4] = byte(x >> 24)
	bb[5] = byte(x >> 16)
	bb[6] = byte(x >> 8)
	bb[7] = byte(x)
}

func getLongB(bb []byte) int64 {
	return makeLong(bb[0],
		bb[1],
		bb[2],
		bb[3],
		bb[4],
		bb[5],
		bb[6],
		bb[7])
}
