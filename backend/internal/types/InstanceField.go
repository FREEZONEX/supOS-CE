package types

type InstanceFieldInfo interface {
	GetId() int64
	GetAlias() string
	GetPath() string
	GetField() string
	GetUts() bool
}

// GetId 获取 Id
func (i *InstanceField) GetId() int64 {
	return i.Id
}

// GetAlias 获取 Alias
func (i *InstanceField) GetAlias() string {
	return i.Alias
}

// GetPath 获取 Path
func (i *InstanceField) GetPath() string {
	return i.Path
}

// GetField 获取 Field
func (i *InstanceField) GetField() string {
	return i.Field
}

// GetUts 获取 Uts
func (i *InstanceField) GetUts() bool {
	return i.Uts
}
