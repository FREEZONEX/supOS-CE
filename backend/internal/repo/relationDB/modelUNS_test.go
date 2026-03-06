package relationDB

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func TestUnsFieldsDecode(t *testing.T) {
	{
		fieldsJson := `[{"name":"timeStamp","type":"DATETIME","systemField":true},{"name":"job_id","type":"STRING","unique":false,"maxLen":"512"},{"name":"_id","type":"LONG","unique":true,"systemField":true}]`
		var f Fields
		err := json.Unmarshal([]byte(fieldsJson), &f)
		if err != nil {
			t.Error(err) //json: invalid use of ,string struct tag, trying to unmarshal unquoted value into *int32
		} else {
			bs, _ := json.MarshalIndent(f, "", " ")
			t.Log(string(bs))
		}
	}
	{
		fieldsJson := `[{"name":"timeStamp","type":"DATETIME","systemField":true},{"name":"job_id","type":"STRING","unique":false,"maxLen":512},{"name":"_id","type":"LONG","unique":true,"systemField":true}]`
		var f Fields
		err := json.Unmarshal([]byte(fieldsJson), &f)
		if err != nil {
			t.Error(err) //json: invalid use of ,string struct tag, trying to unmarshal unquoted value into *int32
		} else {
			bs, _ := json.MarshalIndent(f, "", " ")
			t.Log(string(bs))
		}
	}
}
func TestUnsFieldsDecodeJsonier(t *testing.T) {
	type FieldDefine struct {
		Name        string      `json:"name"`
		Type        string      `json:"type"`
		Unique      *bool       `json:"unique,optional,omitempty"`
		Index       *string     `json:"index,optional,omitempty"`
		DisplayName *string     `json:"displayName,optional,omitempty"`
		Remark      *string     `json:"remark,optional,omitempty"`
		MaxLen      *int32      `json:"maxLen,optional,omitempty"`
		TbValueName *string     `json:"tbValueName,optional,omitempty"`
		Unit        *string     `json:"unit,optional,omitempty"`
		UpperLimit  *float64    `json:"upperLimit,optional,omitempty"`
		LowerLimit  *float64    `json:"lowerLimit,optional,omitempty"`
		Decimal     *int32      `json:"decimal,optional,omitempty"`
		SystemField *bool       `json:"systemField,optional,omitempty"`
		LastValue   string      `json:"-"`
		LastTime    int64       `json:"-"`
		Uns         interface{} `json:"-"`
	}
	extra.RegisterFuzzyDecoders()
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	{
		fieldsJson := `[{"name":"timeStamp","type":"DATETIME","systemField":true},{"name":"job_id","type":"STRING","unique":false,"maxLen":"512"},{"name":"_id","type":"LONG","unique":true,"systemField":true}]`
		var f []*FieldDefine
		err := json.Unmarshal([]byte(fieldsJson), &f)
		if err != nil {
			t.Error(err) //json: invalid use of ,string struct tag, trying to unmarshal unquoted value into *int32
		} else {
			bs, _ := json.MarshalIndent(f, "", " ")
			t.Log(string(bs))
		}
	}
	{
		fieldsJson := `[{"name":"timeStamp","type":"DATETIME","systemField":true},{"name":"job_id","type":"STRING","unique":false,"maxLen":512},{"name":"_id","type":"LONG","unique":true,"systemField":true}]`
		var f []*FieldDefine
		err := json.Unmarshal([]byte(fieldsJson), &f)
		if err != nil {
			t.Error(err) //json: invalid use of ,string struct tag, trying to unmarshal unquoted value into *int32
		} else {
			bs, _ := json.MarshalIndent(f, "", " ")
			t.Log(string(bs))
		}
	}
}
