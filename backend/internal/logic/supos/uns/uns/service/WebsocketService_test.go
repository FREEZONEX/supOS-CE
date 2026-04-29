package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"backend/share/base"
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

func TestProcessWsMsgJsonbUnwrapsJsonbPayload(t *testing.T) {
	def := &types.UnsDefinition{CreateTopicDto: types.CreateTopicDto{
		Alias:     "ac",
		DataSrcID: 2,
		DataType:  base.V2p(constants.JsonbType),
		Fields: []*types.FieldDefine{
			{Name: constants.SysFieldCreateTime, Type: types.FieldTypeDatetime, SystemField: base.OptionalTrue},
			{Name: constants.JsonbField, Type: types.FieldTypeString},
		},
	}}
	msg := serviceApi.WebsocketMessage{
		Def: def,
		Data: map[string]string{
			constants.SysFieldCreateTime: strconv.FormatInt(time.Now().UnixMilli(), 10),
			constants.JsonbField:         `{"debug":1}`,
		},
	}
	rs := processWsMsg(msg)

	var got struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(rs, &got); err != nil {
		t.Fatalf("processWsMsg returned invalid JSON: %v, body=%s", err, string(rs))
	}
	if _, has := got.Data[constants.JsonbField]; has {
		t.Fatalf("jsonb payload should be unwrapped for display, got %#v", got.Data)
	}
	if got.Data["debug"] != float64(1) {
		t.Fatalf("expected unwrapped json payload, got %#v", got.Data)
	}
}
