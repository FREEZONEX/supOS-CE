package bo

type UnsLabels interface {
	UnsId() int64
	LabelNames() []string
	IsResetLabels() bool
	SetLabelId(label string, id int64)
}
