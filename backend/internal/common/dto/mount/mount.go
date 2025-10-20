package mount

// MountSourceDto corresponds to the Java MountSourceDto class.
type MountSourceDto struct {
	SourceAlias string           `json:"sourceAlias"`
	SourceName  string           `json:"sourceName"`
	Devices     []MountDeviceDto `json:"devices"`
}

// MountDto corresponds to the Java MountDto class.
type MountDto struct {
	TargetAlias string         `json:"targetAlias"`
	SourceType  string         `json:"sourceType"` // e.g., "collector", "videoCollector"
	Extend      MountSourceDto `json:"extend"`
	DataType    int            `json:"dataType"`
	Persistence bool           `json:"persistence"`
	Dashboard   bool           `json:"dashboard"`
	SyncMeta    bool           `json:"syncMeta"`
}

// MountDeviceDto corresponds to the Java MountDeviceDto class.
type MountDeviceDto struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
}
