package bo

import (
	"backend/internal/common/dto"
	"backend/internal/repo/relationDB"
)

type UnsPoLabels struct {
	unsPo       *relationDB.UnsNamespace
	labels      []string
	resetLabels bool
	dto         *dto.CreateTopicDto
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
	return u.unsPo.ID
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
func (u *UnsPoLabels) SetDto(d *dto.CreateTopicDto) {
	u.dto = d
	d.LabelIDs = u.labelIds
}
