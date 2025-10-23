package bo

type NodeUnsInfo interface {
	GetId() int64
	GetParentId() *int64
	GetAlias() string
	GetParentAlias() string
	GetName() string
	GetDisplayName() string
	GetPath() string
	GetDataType() *int16
	GetPathType() int16
	GetMountType() *int16
	GetMountSource() string
}
