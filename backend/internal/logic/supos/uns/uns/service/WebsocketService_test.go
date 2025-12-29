package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"testing"
	"time"
)

func Test_processWsMsg(t *testing.T) {
	def := &types.CreateTopicDto{Alias: "ac", DataSrcID: 2, Fields: []*types.FieldDefine{
		{Name: constants.SysFieldCreateTime, Type: types.FieldTypeDatetime},
		{Name: "json", Type: types.FieldTypeString},
	}}
	msg := serviceApi.WebsocketMessage{
		Def: def,
		Data: map[string]any{
			constants.SysFieldCreateTime: time.Now(),
			"json":                       `{"debug":1}`,
		},
	}
	rs := processWsMsg(msg)
	t.Log(string(rs))
}
