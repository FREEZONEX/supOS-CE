package mount

import "backend/internal/common/enums/mount"

// MountMetaQueryEvent defines an event for querying mount metadata.
type MountMetaQueryEvent struct {
	ID        int64
	QueryType mount.MountMetaQueryType
	Param     any
	Callback  func(any)
}
