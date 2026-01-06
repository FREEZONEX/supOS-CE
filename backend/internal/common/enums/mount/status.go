package mount

type MountStatus string

const (
	MountStatusOnline  MountStatus = "online"
	MountStatusOffline MountStatus = "offline"
)

// String returns the string representation
func (m MountStatus) String() string {
	return string(m)
}

// GetMountStatusFromValue returns MountStatus from boolean value
func GetMountStatusFromValue(statusValue bool) MountStatus {
	if statusValue {
		return MountStatusOnline
	}
	return MountStatusOffline
}

// GetMountStatusFromString returns MountStatus from status string
func GetMountStatusFromString(status string) MountStatus {
	if status == MountStatusOnline.String() {
		return MountStatusOnline
	}
	return MountStatusOffline
}
