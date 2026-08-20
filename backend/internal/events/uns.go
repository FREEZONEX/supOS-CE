package events

import "backend/internal/repo"

type UnsNodeSaving struct {
	Action string
	Node   repo.UnsNode
	UserID int64
}

type UnsNodeDeleting struct {
	NodeID int64
	UserID int64
}

type UnsNodeRestoring struct {
	NodeID  int64
	UserID  int64
	Confirm bool
}

type UnsNodeForceDeleting struct {
	NodeID     int64
	UserID     int64
	DeleteFlow bool
}

type UnsNodeCreated struct {
	Node repo.UnsNode
}

type UnsNodesCreated struct {
	Nodes []repo.UnsNode
}

type UnsNodeUpdated struct {
	Node repo.UnsNode
}

type UnsNodeDeleted struct {
	RootID int64
	UserID int64
	Nodes  []repo.UnsNode
}

type UnsNodeRestored struct {
	RootID int64
	UserID int64
	Nodes  []repo.UnsNode
}

type UnsNodeForceDeleted struct {
	RootID     int64
	UserID     int64
	DeleteFlow bool
	Nodes      []repo.UnsNode
}
