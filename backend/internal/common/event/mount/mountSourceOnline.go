package mount

import "backend/internal/common/enums/mount"

// MountSourceOnlineEvent defines an event related to a mount source's online/offline status.
type MountSourceOnlineEvent struct {
	SourceType  mount.MountSourceType
	SourceAlias string
	Callback    func(bool)
}
