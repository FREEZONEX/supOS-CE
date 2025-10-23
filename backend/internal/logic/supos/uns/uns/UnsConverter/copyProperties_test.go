package UnsConverter

import (
	"backend/internal/common/dto"
	"backend/internal/common/enums"
	"backend/internal/types"
	"encoding/json"
	"testing"

	"gitee.com/unitedrhino/share/errors"
	"github.com/jinzhu/copier"
)

func TestCopyFields(t *testing.T) {
	src := types.CreateTopicDto{Id: 123, Name: "test123", PathType: 2, DataType: 1, Fields: []types.FieldDefine{
		{
			Name: "id", Type: "LONG",
		}, {
			Name: "ts", Type: "DATETIME",
		},
	}, Refers: []types.InstanceField{
		{Id: 10001, Alias: "A1-1"},
	}, Extend: map[string]interface{}{
		"Debug": true,
	},
	}
	target := dto.CreateTopicDto{}
	err := copier.CopyWithOption(&target, src, copier.Option{IgnoreEmpty: true, Converters: []copier.TypeConverter{
		{
			SrcType: copier.String,
			DstType: enums.FieldTypeInteger,
			Fn: func(src interface{}) (dst interface{}, err error) {
				if rs, ok := enums.GetFieldTypeByNameIgnoreCase(src.(string)); ok {
					return rs, nil
				}
				return nil, errors.Default
			},
		},
	}})
	bs, _ := json.MarshalIndent(target, "", " ")
	t.Logf("Copy:%v, rs: %s", err, string(bs))
}
