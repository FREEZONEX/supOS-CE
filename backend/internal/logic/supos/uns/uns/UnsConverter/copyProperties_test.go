package UnsConverter

import (
	"backend/internal/types"
	"backend/share/base"
	"encoding/json"
	"testing"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"github.com/jinzhu/copier"
)

func TestCopyFields(t *testing.T) {
	dataType := int16(1)
	src := types.CreateTopicDto{Id: 123, Name: "test123", PathType: 2, DataType: &dataType, Fields: []*types.FieldDefine{
		{
			Name: "id", Type: "LONG",
		}, {
			Name: "ts", Type: "DATETIME",
		},
	}, Refers: []*types.InstanceField{
		{Id: 10001, Alias: "A1-1"},
	}, Extend: map[string]interface{}{
		"Debug": true,
	},
	}
	target := types.CreateTopicDto{}
	options := []copier.TypeConverter{
		{
			SrcType: copier.String,
			DstType: types.FieldTypeInteger,
			Fn: func(src interface{}) (dst interface{}, err error) {
				if rs, ok := types.GetFieldTypeByNameIgnoreCase(src.(string)); ok {
					return rs, nil
				}
				return nil, errors.Default
			},
		}}
	err := copier.CopyWithOption(&target, src, copier.Option{IgnoreEmpty: true, Converters: options})
	bs, _ := json.MarshalIndent(target, "", " ")
	t.Logf("Copy:%v, rs: %s", err, string(bs))

	srcList := []*types.CreateTopicDto{&src}
	tarList := make([]*types.CreateTopicDto, len(srcList))
	err = copier.CopyWithOption(&tarList, srcList, copier.Option{IgnoreEmpty: true, Converters: options})
	bs, _ = json.MarshalIndent(tarList, "", " ")
	t.Logf("Copy:%v, rs: %s", err, string(bs))

}
func TestCopyTime2Long(t *testing.T) {
	type User struct {
		Id          int64     `json:"id,omitzero"`
		Name        string    `gorm:"column:name;not null" json:"name"`
		Description *string   `json:"description,optional,omitzero" validate:"max=255"`
		CreateAt    time.Time `json:"createAt,omitzero"`
	}
	type UserVo struct {
		Id          string  `json:"id,omitzero"`
		Name        string  `gorm:"column:name;not null" json:"name"`
		Description *string `gorm:"column:description" json:"description"`
		CreateAt    int64   `json:"createAt,omitzero"`
	}
	po := &User{Id: 123, CreateAt: time.Now(), Name: "", Description: base.V2p("")}
	vo := &UserVo{Name: "X", Description: base.V2p("des")}
	err := copier.CopyWithOption(vo, po, apiConvertOptions)
	if err != nil {
		t.Log(err)
	}
	t.Logf("t0=%d,t1=%d, desc=%s\n", po.CreateAt.UnixMilli(), vo.CreateAt, *vo.Description)
	t.Logf("vo=%+v\n", *vo)
}
