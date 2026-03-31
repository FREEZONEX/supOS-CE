package bo

import (
	"backend/internal/repo/relationDB"
	"backend/internal/types"
)

type UnsPoLabels struct {
	UnsPo       *relationDB.UnsNamespace
	labels      []string
	resetLabels bool
}

func NewUnsPoLabels(unsPo *relationDB.UnsNamespace, resetLabels bool, labels []string) *UnsPoLabels {
	unsPo.LabelIds = make(map[int64]string)
	return &UnsPoLabels{
		UnsPo:       unsPo,
		resetLabels: resetLabels,
		labels:      labels,
	}
}
func (u *UnsPoLabels) UnsId() int64 {
	return u.UnsPo.Id
}
func (u *UnsPoLabels) LabelNames() []string {
	return u.labels
}
func (u *UnsPoLabels) IsResetLabels() bool {
	return u.resetLabels
}
func (u *UnsPoLabels) SetLabelId(label string, id int64) {
	u.UnsPo.LabelIds[id] = label
}
func (u *UnsPoLabels) SetDto(d *types.CreateTopicDto) {
	d.LabelIDs = u.UnsPo.LabelIds
}
