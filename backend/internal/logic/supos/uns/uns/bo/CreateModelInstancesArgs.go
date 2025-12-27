package bo

import (
	"backend/internal/common"
	"backend/internal/types"
	"encoding/json"
)

type CreateModelInstancesArgs struct {
	Topics                        []*types.CreateTopicDto     `json:"topics"`
	FromImport                    bool                        `json:"fromImport"`
	RetainTableWhenDeleteInstance bool                        `json:"retainTableWhenDeleteInstance"`
	ThrowModelExistsErr           bool                        `json:"throwModelExistsErr"`
	FlowName                      string                      `json:"flowName"`
	LabelsMap                     map[string][]string         `json:"labelsMap"`
	SkipWhenExists                bool                        `json:"skipWhenExists"` // true: 已存在的 UNS 不更新
	StatusConsumer                func(*common.RunningStatus) `json:"-"`              // 使用json忽略标记
}

func (c CreateModelInstancesArgs) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}
