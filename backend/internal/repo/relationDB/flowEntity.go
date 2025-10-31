package relationDB

// FlowEntity defines the behaviour shared by source/event Node-RED flow records.
type FlowEntity interface {
	GetID() int64
	SetID(int64)
	GetFlowID() string
	SetFlowID(string)
	GetFlowName() string
	SetFlowName(string)
	GetFlowData() string
	SetFlowData(string)
	GetFlowStatus() string
	SetFlowStatus(string)
	GetTemplate() string
	SetTemplate(string)
	GetDescription() string
	SetDescription(string)
	// GetFlowType() int32
	// SetFlowType(int32)
}

// --- NoderedSourceFlow implements FlowEntity ---

func (f *NoderedSourceFlow) GetID() int64 {
	if f == nil {
		return 0
	}
	return f.ID
}

func (f *NoderedSourceFlow) SetID(id int64) {
	if f != nil {
		f.ID = id
	}
}

func (f *NoderedSourceFlow) GetFlowID() string {
	if f == nil {
		return ""
	}
	return f.FlowID
}

func (f *NoderedSourceFlow) SetFlowID(flowID string) {
	if f != nil {
		f.FlowID = flowID
	}
}

func (f *NoderedSourceFlow) GetFlowName() string {
	if f == nil {
		return ""
	}
	return f.FlowName
}

func (f *NoderedSourceFlow) SetFlowName(name string) {
	if f != nil {
		f.FlowName = name
	}
}

func (f *NoderedSourceFlow) GetFlowData() string {
	if f == nil {
		return ""
	}
	return f.FlowData
}

func (f *NoderedSourceFlow) SetFlowData(data string) {
	if f != nil {
		f.FlowData = data
	}
}

func (f *NoderedSourceFlow) GetFlowStatus() string {
	if f == nil {
		return ""
	}
	return f.FlowStatus
}

func (f *NoderedSourceFlow) SetFlowStatus(status string) {
	if f != nil {
		f.FlowStatus = status
	}
}

func (f *NoderedSourceFlow) GetTemplate() string {
	if f == nil {
		return ""
	}
	return f.Template
}

func (f *NoderedSourceFlow) SetTemplate(tpl string) {
	if f != nil {
		f.Template = tpl
	}
}

func (f *NoderedSourceFlow) GetDescription() string {
	if f == nil {
		return ""
	}
	return f.Description
}

func (f *NoderedSourceFlow) SetDescription(desc string) {
	if f != nil {
		f.Description = desc
	}
}

// func (f *NoderedSourceFlow) GetFlowType() int32 {
// 	if f == nil {
// 		return 0
// 	}
// 	return f.FlowType
// }

// func (f *NoderedSourceFlow) SetFlowType(t int32) {
// 	if f != nil {
// 		f.FlowType = t
// 	}
// }

// --- NoderedEventFlow implements FlowEntity ---

func (f *NoderedEventFlow) GetID() int64 {
	if f == nil {
		return 0
	}
	return f.ID
}

func (f *NoderedEventFlow) SetID(id int64) {
	if f != nil {
		f.ID = id
	}
}

func (f *NoderedEventFlow) GetFlowID() string {
	if f == nil {
		return ""
	}
	return f.FlowID
}

func (f *NoderedEventFlow) SetFlowID(flowID string) {
	if f != nil {
		f.FlowID = flowID
	}
}

func (f *NoderedEventFlow) GetFlowName() string {
	if f == nil {
		return ""
	}
	return f.FlowName
}

func (f *NoderedEventFlow) SetFlowName(name string) {
	if f != nil {
		f.FlowName = name
	}
}

func (f *NoderedEventFlow) GetFlowData() string {
	if f == nil {
		return ""
	}
	return f.FlowData
}

func (f *NoderedEventFlow) SetFlowData(data string) {
	if f != nil {
		f.FlowData = data
	}
}

func (f *NoderedEventFlow) GetFlowStatus() string {
	if f == nil {
		return ""
	}
	return f.FlowStatus
}

func (f *NoderedEventFlow) SetFlowStatus(status string) {
	if f != nil {
		f.FlowStatus = status
	}
}

func (f *NoderedEventFlow) GetTemplate() string {
	if f == nil {
		return ""
	}
	return f.Template
}

func (f *NoderedEventFlow) SetTemplate(tpl string) {
	if f != nil {
		f.Template = tpl
	}
}

func (f *NoderedEventFlow) GetDescription() string {
	if f == nil {
		return ""
	}
	return f.Description
}

func (f *NoderedEventFlow) SetDescription(desc string) {
	if f != nil {
		f.Description = desc
	}
}

// func (f *NoderedEventFlow) GetFlowType() int32 {
// 	if f == nil {
// 		return 0
// 	}
// 	return f.FlowType
// }

// func (f *NoderedEventFlow) SetFlowType(t int32) {
// 	if f != nil {
// 		f.FlowType = t
// 	}
// }
