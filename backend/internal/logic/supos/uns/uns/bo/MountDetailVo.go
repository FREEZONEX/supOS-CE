package bo

type MountDetailVo struct {
	// MountType 挂载类型
	MountType *int `json:"mountType,omitempty"`

	// MountSource 挂载源
	MountSource string `json:"mountSource,omitempty"`

	// DisplayName 挂载源显示名
	DisplayName string `json:"displayName,omitempty"`
}
