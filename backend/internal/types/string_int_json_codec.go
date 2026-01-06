package types

import (
	"encoding/json"
	"strconv"
)

// Int64 自定义 int64 类型，序列化为字符串
type Int64 int64
type Float64 float64
type Float32 float32

// MarshalJSON 实现 json.Marshaler 接口
func (s Int64) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatInt(int64(s), 10) + `"`), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (s *Int64) UnmarshalJSON(data []byte) error {
	// 处理字符串格式
	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		str := string(data[1 : len(data)-1])
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		*s = Int64(val)
		return nil
	}

	// 处理数字格式（向后兼容）
	var val int64
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	*s = Int64(val)
	return nil
}

// Float64 自定义 float64 类型，序列化为字符串

// MarshalJSON 实现 json.Marshaler 接口
func (s Float64) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatFloat(float64(s), 'f', -1, 64) + `"`), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (s *Float64) UnmarshalJSON(data []byte) error {
	// 处理字符串格式
	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		str := string(data[1 : len(data)-1])
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		*s = Float64(val)
		return nil
	}

	// 处理数字格式（向后兼容）
	var val float64
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	*s = Float64(val)
	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (s Float32) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatFloat(float64(s), 'f', -1, 32) + `"`), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (s *Float32) UnmarshalJSON(data []byte) error {
	// 处理字符串格式
	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		str := string(data[1 : len(data)-1])
		val, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return err
		}
		*s = Float32(val)
		return nil
	}

	// 处理数字格式（向后兼容）
	var val float32
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	*s = Float32(val)
	return nil
}
