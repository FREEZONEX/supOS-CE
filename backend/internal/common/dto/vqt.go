package dto

import (
	"backend/internal/types"
	"encoding/json"
)

// VQT represents Value-Quality-Timestamp data structure
type VQT[T any] struct {
	Timestamp int64 `json:"timeStamp"` // 时间戳
	Quality   int64 `json:"quality"`   // 质量
	Value     T     `json:"value"`     // 值
}

// NewVQT creates a new VQT instance
func NewVQT[T any](timestamp, quality int64, value T) *VQT[T] {
	return &VQT[T]{
		Timestamp: timestamp,
		Quality:   quality,
		Value:     value,
	}
}

// GetVQTJSON generates VQT JSON string from topic metadata
func GetVQTJSON(t, q int64, v any, metaDef *types.CreateTopicDto) (string, error) {
	vqtMap := GetVQTMap(t, q, v, metaDef)
	data, err := json.Marshal(vqtMap)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetVQTMap generates VQT map from topic metadata
func GetVQTMap(t, q int64, v any, metaDef *types.CreateTopicDto) map[string]any {
	vqt := make(map[string]any)

	timestampField := metaDef.GetTimestampField()
	qualityField := metaDef.GetQualityField()

	vqt[timestampField] = t
	vqt[qualityField] = q

	// Find the first non-system field and set value
	for _, fd := range metaDef.Fields {
		if fd.Name != timestampField && fd.Name != qualityField {
			vqt[fd.Name] = v
			break
		}
	}

	return vqt
}
