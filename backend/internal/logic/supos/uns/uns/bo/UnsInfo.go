package bo

import (
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
)

type UnsInfo interface {
	GetId() int64
	GetParentId() *int64
	GetAlias() string
	GetParentAlias() *string
	GetName() string
	GetDisplayName() string
	GetPath() string
	GetLayRec() string
	GetDataPath() string
	GetProtocolMap() map[string]interface{}
	GetDescription() string
	GetExpression() string
	GetCalculationType() *int32
	//GetReadWriteMode() string
	GetExtend() map[string]interface{}
	GetLabelIds() map[int64]string
	GetModelId() *int64
	GetDataType() *int16
	GetParentDataType() *int16
	GetPathType() int16
	GetRefers() []*types.InstanceField
	GetFields() []*types.FieldDefine
	GetRefUns() map[int64]int
	GetCreateAt() int64
	GetUpdateAt() int64
	GetFlags() *int32
	GetExtendFieldFlags() *int32
	GetMountType() *int16
	GetMountSource() string
	GetSrcJdbcType() types.SrcJdbcType
	GetStatus() *int16
}

var _ UnsInfo = &types.CreateTopicDto{}
var _ UnsInfo = &dao.UnsNamespace{}
