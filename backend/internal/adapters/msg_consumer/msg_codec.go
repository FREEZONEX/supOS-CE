package msg_consumer

import (
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"backend/share/base/bits"
	"context"
)

func encodeMsg(ctx context.Context, unsData []serviceApi.TopicMessage) []byte {
	size := 4
	for _, m := range unsData {
		size += 8 + 2 + 4
		for _, vs := range m.Data {
			size += 2
			for k, v := range vs {
				size += 2 + len(k) + 2 + len(v)
			}
		}

	}
	bs := make([]byte, size)
	pos := 0
	bits.BigEndian.PutUint32(bs, uint32(len(unsData)))
	pos += 4
	for _, m := range unsData {
		bits.BigEndian.PutLong(bs[pos:], m.UnsId)
		pos += 8
		bits.BigEndian.PutShort(bs[pos:], int16(m.DataSrcId))
		pos += 2
		bits.BigEndian.PutUint32(bs[pos:], uint32(len(m.Data)))
		pos += 4
		for _, vs := range m.Data {
			bits.BigEndian.PutUint16(bs[pos:], uint16(len(vs)))
			pos += 2
			for k, v := range vs {
				//if strings.HasPrefix(k, "double") {
				//	_, Nerr := strconv.ParseFloat(v, 64)
				//	if Nerr != nil {
				//		payload := ctx.Value("payload")
				//		logx.Errorf("源头浮点格式错误:%v, str=%s, payload=%v", Nerr, v, payload)
				//	}
				//}
				bits.BigEndian.PutUint16(bs[pos:], uint16(len(k)))
				pos += 2
				copy(bs[pos:], k)
				pos += len(k)

				bits.BigEndian.PutUint16(bs[pos:], uint16(len(v)))
				pos += 2
				copy(bs[pos:], v)
				pos += len(v)
			}
		}
	}
	return bs
}
func decodeMsg(data []byte, unsData *[]serviceApi.TopicMessage) {
	size := bits.BigEndian.Uint32(data)
	if size < 0 {
		return
	}
	pos := 4
	*unsData = make([]serviceApi.TopicMessage, size)
	for i := 0; i < int(size); i++ {
		msg := serviceApi.TopicMessage{}
		msg.UnsId = bits.BigEndian.GetLong(data[pos:])
		pos += 8
		msg.DataSrcId = types.SrcJdbcType(bits.BigEndian.GetShort(data[pos:]))
		pos += 2
		countMap := bits.BigEndian.Uint32(data[pos:])
		pos += 4
		msg.Data = make([]map[string]string, countMap)
		for j := 0; j < int(countMap); j++ {
			mapSize := bits.BigEndian.Uint16(data[pos:])
			pos += 2
			vmap := make(map[string]string, mapSize)
			for k := 0; k < int(mapSize); k++ {
				kLen := bits.BigEndian.Uint16(data[pos:])
				pos += 2
				key := data[pos : pos+int(kLen)]
				pos += int(kLen)

				vLen := bits.BigEndian.Uint16(data[pos:])
				pos += 2
				value := data[pos : pos+int(vLen)]
				pos += int(vLen)

				kStr, valueStr := b2s(key), b2s(value)
				vmap[kStr] = valueStr
			}
			msg.Data[j] = vmap
		}
		(*unsData)[i] = msg
	}
}
