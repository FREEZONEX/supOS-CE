package uns

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// Definition mirrors UNS field definition shape used in schemas/templates.
type Definition struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Unit        string `json:"unit"`
	Unique      bool   `json:"unique"`
	SystemField bool   `json:"system_field"`
}

// DefinitionListJSON accepts either a JSON string or an array of Definition.
type DefinitionListJSON []Definition

func (d DefinitionListJSON) Value() (driver.Value, error) {
	if len(d) == 0 {
		return []byte("[]"), nil
	}
	buf, err := json.Marshal([]Definition(d))
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *DefinitionListJSON) Scan(value interface{}) error {
	if value == nil {
		*d = nil
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("DefinitionListJSON: unsupported Scan type %T", value)
	}
	return json.Unmarshal(data, d)
}

// 上面 FilterValid 方法参考前一条回复
func (d DefinitionListJSON) FilterValid() DefinitionListJSON {
	var out DefinitionListJSON
	for _, def := range d {
		if strings.TrimSpace(def.Name) != "" && strings.TrimSpace(def.Type) != "" {
			out = append(out, def)
		}
	}
	return out
}

func (d *DefinitionListJSON) UnmarshalJSON(b []byte) error {
	if d == nil {
		return fmt.Errorf("nil DefinitionListJSON receiver")
	}
	trimmed := bytes.TrimSpace(b)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*d = nil
		return nil
	}
	if trimmed[0] == '"' { // string
		var str string
		if err := json.Unmarshal(trimmed, &str); err != nil {
			return err
		}
		return d.setFromString(str)
	}
	if trimmed[0] != '[' { // must be array
		return fmt.Errorf("schema must be a JSON array")
	}
	var defs []Definition
	if err := json.Unmarshal(trimmed, &defs); err != nil {
		return fmt.Errorf("schema must be a JSON array: %w", err)
	}
	if err := validateDefinitions(defs); err != nil {
		return err
	}
	*d = defs
	return nil
}

func (d DefinitionListJSON) MarshalJSON() ([]byte, error) {
	if d == nil {
		return json.Marshal([]Definition(nil))
	}
	return json.Marshal([]Definition(d))
}

func (d DefinitionListJSON) ToDefinitions() []Definition {
	if len(d) == 0 {
		return nil
	}
	out := make([]Definition, len(d))
	copy(out, d)
	return out
}

func (d DefinitionListJSON) String() string {
	if len(d) == 0 {
		return ""
	}
	b, err := json.Marshal([]Definition(d))
	if err != nil {
		return ""
	}
	return string(b)
}

func (d *DefinitionListJSON) setFromString(str string) error {
	trimmed := strings.TrimSpace(str)
	if trimmed == "" {
		*d = nil
		return nil
	}
	var defs []Definition
	if err := json.Unmarshal([]byte(trimmed), &defs); err != nil {
		return fmt.Errorf("schema must be a JSON array: %w", err)
	}
	if err := validateDefinitions(defs); err != nil {
		return err
	}
	*d = defs
	return nil
}

func validateDefinitions(defs []Definition) error {
	var filtered []Definition
	for _, def := range defs {
		if strings.TrimSpace(def.Name) == "" || strings.TrimSpace(def.Type) == "" {
			continue // 忽略掉 name 或 type 为空的项
		}
		filtered = append(filtered, def)
	}
	// 如果需要返回过滤后的结果，可以在调用处处理；这里只做校验
	return nil
}
