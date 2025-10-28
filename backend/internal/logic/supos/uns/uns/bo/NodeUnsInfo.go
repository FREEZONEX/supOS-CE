package bo

import (
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
)

type NodeUnsInfo interface {
	GetId() int64
	GetParentId() *int64
	GetAlias() string
	GetParentAlias() *string
	GetName() string
	GetDisplayName() string
	GetPath() string
	GetDataType() *int16
	GetParentDataType() *int16
	GetPathType() int16
	GetMountType() *int16
	GetMountSource() string
}

var _ NodeUnsInfo = &types.CreateTopicDto{}
var _ NodeUnsInfo = &dao.UnsNamespace{}
var _ NodeUnsInfo = &dao.TreeNodeUns{}
var _ NodeUnsInfo = &dao.UnsLabel{}
