package bo

type NodeUnsInfo interface {
	GetId() int64
	GetParentId() int64
	GetAlias() string
	GetParentAlias() string
	GetName() string
	GetDisplayName() string
	GetPath() string
	GetDataType() int32
	GetPathType() int32
	GetMountType() int32
	GetMountSource() string
}
