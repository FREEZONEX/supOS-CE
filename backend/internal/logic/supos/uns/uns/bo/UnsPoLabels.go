package bo

import (
	"backend/internal/repo/relationDB"
	"backend/internal/types"
)

type UnsPoLabels struct {
	unsPo       *relationDB.UnsNamespace
	labels      []string
	resetLabels bool
	dto         *types.CreateTopicDto
	labelIds    map[int64]string
}

func NewUnsPoLabels(unsPo *relationDB.UnsNamespace, resetLabels bool, labels []string) *UnsPoLabels {
	labelIds := make(map[int64]string)
	unsPo.LabelIds = labelIds
	return &UnsPoLabels{
		labelIds:    labelIds,
		unsPo:       unsPo,
		resetLabels: resetLabels,
		labels:      labels,
	}
}
func (u *UnsPoLabels) UnsId() int64 {
	return u.unsPo.Id
}
func (u *UnsPoLabels) LabelNames() []string {
	return u.labels
}
func (u *UnsPoLabels) IsResetLabels() bool {
	return u.resetLabels
}
func (u *UnsPoLabels) SetLabelId(label string, id int64) {
	u.labelIds[id] = label
}
func (u *UnsPoLabels) SetDto(d *types.CreateTopicDto) {
	u.dto = d
	d.LabelIDs = u.labelIds
}
