package mount

// MountTargetType represents mount target types
type MountTargetType string

const (
	MountTargetFolder MountTargetType = "folder" // 挂载到文件夹
	MountTargetFile   MountTargetType = "file"   // 挂载到文件
)

// Type returns the type string
func (m MountTargetType) Type() string {
	return string(m)
}

// String returns the string representation
func (m MountTargetType) String() string {
	return string(m)
}

// GetMountTargetTypeFromType returns MountTargetType from type string
func GetMountTargetTypeFromType(typeStr string) (MountTargetType, bool) {
	switch MountTargetType(typeStr) {
	case MountTargetFolder, MountTargetFile:
		return MountTargetType(typeStr), true
	default:
		return "", false
	}
}
