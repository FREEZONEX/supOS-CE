package mount

// MountEventType represents mount event types
type MountEventType string

const (
	MountEventCollectorModify MountEventType = "COLLECTOR_MODIFY"
	MountEventCollectorDelete MountEventType = "COLLECTOR_DELETE"

	MountEventSourceAdd    MountEventType = "SOURCE_ADD"
	MountEventSourceModify MountEventType = "SOURCE_MODIFY"
	MountEventSourceDelete MountEventType = "SOURCE_DELETE"

	MountEventTagAdd    MountEventType = "TAG_ADD"
	MountEventTagModify MountEventType = "TAG_MODIFY"
	MountEventTagDelete MountEventType = "TAG_DELETE"

	MountEventVideoTagAdd    MountEventType = "VIDEO_TAG_ADD"
	MountEventVideoTagModify MountEventType = "VIDEO_TAG_MODIFY"
	MountEventVideoTagDelete MountEventType = "VIDEO_TAG_DELETE"
)

// String returns the string representation
func (m MountEventType) String() string {
	return string(m)
}

// CanBatchGroup checks if two event types can be batched together
func CanBatchGroup(type1, type2 MountEventType) bool {
	switch {
	case type1 == MountEventTagAdd && type2 == MountEventTagAdd:
		return true
	case type1 == MountEventTagAdd && type2 == MountEventTagModify:
		return true
	case type1 == MountEventTagModify && type2 == MountEventTagAdd:
		return true
	case type1 == MountEventTagModify && type2 == MountEventTagModify:
		return true
	case type1 == MountEventTagDelete && type2 == MountEventTagDelete:
		return true
	default:
		return false
	}
}
