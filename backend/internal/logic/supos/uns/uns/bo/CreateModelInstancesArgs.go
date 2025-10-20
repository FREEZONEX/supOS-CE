package bo

import (
	"backend/internal/common"
	"backend/internal/common/dto"
	"encoding/json"
)

type CreateModelInstancesArgs struct {
	Topics                        []*dto.CreateTopicDto       `json:"topics"`
	FromImport                    bool                        `json:"fromImport"`
	RetainTableWhenDeleteInstance bool                        `json:"retainTableWhenDeleteInstance"`
	ThrowModelExistsErr           bool                        `json:"throwModelExistsErr"`
	FlowName                      string                      `json:"flowName"`
	LabelsMap                     map[string][]string         `json:"labelsMap"`
	StatusConsumer                func(*common.RunningStatus) `json:"-"` // 使用json忽略标记
}

func (c *CreateModelInstancesArgs) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}
