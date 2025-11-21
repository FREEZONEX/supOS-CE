package postgresql

import (
	"backend/internal/common"
	"backend/internal/common/constants"
	"backend/internal/types"
	"backend/share/base"
	"fmt"
	"testing"
)

func TestGetCreateTableSQL(t *testing.T) {
	fields := []*types.FieldDefine{
		{
			Name:   "id",
			Type:   types.FieldTypeInteger,
			Unique: base.OptionalTrue,
		},
		{
			Name:   "dev_id",
			Type:   types.FieldTypeLong,
			Unique: base.OptionalTrue,
		},
		{
			Name: "name",
			Type: types.FieldTypeString,
		},
		{
			Name: constants.SysFieldCreateTime,
			Type: types.FieldTypeDatetime,
		}}
	common.InitSnowflake(123)
	tableName := "public.users"
	sql := getCreateTableSQL(false, tableName, fields)
	fmt.Println(sql)
}
