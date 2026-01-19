package msg_consumer

import (
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"encoding/json"
	"testing"
)

func TestTopicMsgCodec(t *testing.T) {
	msg1 := serviceApi.TopicMessage{UnsId: 1, DataSrcId: types.SrcJdbcTypePostgresql, Data: []map[string]string{
		{
			"ts":   "1",
			"type": "pg",
		}, {
			"ts":   "2",
			"type": "rs",
		},
	}}
	msg2 := serviceApi.TopicMessage{UnsId: 2, DataSrcId: types.SrcJdbcTypeTimeScaleDB, Data: []map[string]string{
		{
			"ts":      "3",
			"double1": "-52.99193072175547",
		}, {
			"ts":  "4",
			"wet": "93.9692603150808",
		},
	}}
	params := []serviceApi.TopicMessage{msg1, msg2}
	bs := encodeMsg(t.Context(), params)
	jsonBs, _ := json.Marshal(params)
	t.Log(len(bs), len(jsonBs), string(jsonBs))

	var rs []serviceApi.TopicMessage
	decodeMsg(bs, &rs)
	t.Log(rs)
}
