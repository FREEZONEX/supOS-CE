package bits

type ByteOrder interface {
	PutShort(bb []byte, x int16)
	GetShort(bb []byte) int16

	PutInt(bb []byte, x int)
	GetInt(bb []byte) int

	PutLong(bb []byte, x int64)
	GetLong(bb []byte) int64
	//
	Uint16([]byte) uint16
	Uint32([]byte) uint32
	Uint64([]byte) uint64
	PutUint16([]byte, uint16)
	PutUint32([]byte, uint32)
	PutUint64([]byte, uint64)
	String() string
}

func makeShort(b1, b0 byte) int16 {
	return (int16(int8(b1)) << 8) | (int16(int8(b0)) & 0xff)
}

func makeInt(b3, b2, b1, b0 byte) int {
	return ((int(int8(b3)) << 24) |
		((int(int8(b2)) & 0xff) << 16) |
		((int(int8(b1)) & 0xff) << 8) |
		(int(int8(b0)) & 0xff))
}
func makeLong(b7, b6, b5, b4,
	b3, b2, b1, b0 byte) int64 {
	return ((int64(int8(b7)) << 56) |
		((int64(int8(b6)) & 0xff) << 48) |
		((int64(int8(b5)) & 0xff) << 40) |
		((int64(int8(b4)) & 0xff) << 32) |
		((int64(int8(b3)) & 0xff) << 24) |
		((int64(int8(b2)) & 0xff) << 16) |
		((int64(int8(b1)) & 0xff) << 8) |
		(int64(int8(b0)) & 0xff))
}
