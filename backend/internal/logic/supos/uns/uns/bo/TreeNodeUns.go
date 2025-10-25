package bo

type TreeNodeUns struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	DisplayName   string  `json:"displayName"`
	Path          string  `json:"path"`
	Alias         string  `json:"alias"`
	ParentAlias   *string `json:"parentAlias"`
	ParentID      *int64  `json:"parentId"`
	PathType      int16   `json:"pathType"`
	DataType      *int16  `json:"dataType"`
	MountType     *int16  `json:"mountType"`
	MountSource   string  `json:"mountSource"`
	CountChildren string  `json:"countChildren"`
}

var _ NodeUnsInfo = &TreeNodeUns{}

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

func (t *TreeNodeUns) GetPathType() int16 {
	return t.PathType
}

func (t *TreeNodeUns) GetMountType() *int16 {
	return t.MountType
}

func (t *TreeNodeUns) GetMountSource() string {
	return t.MountSource
}
