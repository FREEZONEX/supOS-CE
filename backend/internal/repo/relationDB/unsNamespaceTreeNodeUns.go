package relationDB

type TreeNodeUns struct {
	ID             int64   `gorm:"column:id;primaryKey" json:"id"`
	Name           string  `gorm:"column:name;" json:"name"`
	DisplayName    string  `gorm:"column:display_name" json:"displayName"`
	Path           string  `gorm:"column:path;" json:"path"`
	Alias          string  `gorm:"column:alias;" json:"alias"`
	ParentAlias    *string `gorm:"column:parent_alias" json:"parentAlias"`
	ParentID       *int64  `gorm:"column:parent_id" json:"parentId"`
	PathType       int16   `gorm:"column:path_type;" json:"pathType"`
	DataType       *int16  `gorm:"column:data_type;" json:"dataType"`
	ParentDataType *int16  `gorm:"column:parent_data_type" json:"parent_data_type"`
	MountType      *int16  `gorm:"column:mount_type;" json:"mountType"`
	MountSource    string  `gorm:"column:mount_source;" json:"mountSource"`
	CountChildren  string  `gorm:"column:count_children;" json:"countChildren"`
}

// 实现NodeUnsInfo接口的方法
func (t *TreeNodeUns) GetId() int64 {
	return t.ID
}

func (t *TreeNodeUns) GetParentId() *int64 {
	return t.ParentID
}

func (t *TreeNodeUns) GetAlias() string {
	return t.Alias
}

func (t *TreeNodeUns) GetParentAlias() *string {
	return t.ParentAlias
}

func (t *TreeNodeUns) GetName() string {
	return t.Name
}

func (t *TreeNodeUns) GetDisplayName() string {
	return t.DisplayName
}

func (t *TreeNodeUns) GetPath() string {
	return t.Path
}

func (t *TreeNodeUns) GetDataType() *int16 {
	return t.DataType
}
func (t *TreeNodeUns) GetParentDataType() *int16 {
	return t.ParentDataType
}
func (t *TreeNodeUns) GetPathType() int16 {
	return t.PathType
}

func (t *TreeNodeUns) GetMountType() *int16 {
	return t.MountType
}

func (t *TreeNodeUns) GetMountSource() string {
	return t.MountSource
}
