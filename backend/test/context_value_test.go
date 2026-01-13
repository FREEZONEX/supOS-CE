package test

import (
	"backend/internal/common/utils/datetimeutils"
	"backend/share/base"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestJson(t *testing.T) {
	type FieldDefine struct {
		Name        string      `json:"name"`
		Type        string      `json:"type"`
		Unique      *bool       `json:"unique,optional,omitempty"`
		Index       *string     `json:"index,optional,omitempty"`
		DisplayName *string     `json:"displayName,optional,omitempty"`
		Remark      *string     `json:"remark,optional,omitempty"`
		MaxLen      *int        `json:"maxLen,optional,omitempty,string"`
		TbValueName *string     `json:"tbValueName,optional,omitempty"`
		Unit        *string     `json:"unit,optional,omitempty"`
		UpperLimit  *float64    `json:"upperLimit,optional,omitempty"`
		LowerLimit  *float64    `json:"lowerLimit,optional,omitempty"`
		Decimal     *int        `json:"decimal,optional,omitempty,string"`
		SystemField *bool       `json:"systemField,optional,omitempty"`
		LastValue   interface{} `json:"-,optional"`
		LastTime    int64       `json:"-,optional"`
		Uns         interface{} `json:"-"`
	}
	var def FieldDefine
	json.Unmarshal([]byte(`{
                              "name": "job_id",
                              "type": "STRING",
                              "unique": false,
                              "maxLen": "512"
                            }`), &def)
	t.Log(def.Name, base.P2v(def.MaxLen))
}
func TestCtxValue(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "test", 1)
	ctx = context.WithValue(ctx, "debug", true)

	t.Log(ctx.Value("test"))
	t.Log(ctx.Value("debug"))

	curT := "1234555.45"
	ct, err := strconv.ParseFloat(fmt.Sprint(curT), 64)
	t.Log(ct, err)
	tm := time.UnixMilli(int64(ct))
	t.Log(tm, tm.IsZero())
	str := fmt.Sprint(time.Now().UnixMilli())
	t.Log(len(str), ", ", str)
	dates := []string{str, "2025-12-25T12:52:27.001Z", "2025/12/25 20:52:27[UTC+8]"}
	for _, v := range dates {
		dt, er := datetimeutils.ParseDate(v)
		t.Logf("yeay=%d, dt=%s, er=%v", dt.Year(), dt, er)
	}

	t.Log(datetimeutils.DateTimeUTC(base.V2p(int64(1766998230001))))
}
