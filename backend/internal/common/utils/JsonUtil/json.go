package JsonUtil

import (
	"backend/internal/common/utils/integerutil"
	"encoding/json"
	"strings"
)

// ToJson serializes an object to a JSON string.
func ToJson(obj any) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ToJsonBytes serializes an object to a JSON byte slice.
func ToJsonBytes(obj any) ([]byte, error) {
	return json.Marshal(obj)
}

// FromJson deserializes a JSON string into an object of the specified type.
func FromJson(jsonStr string, v any) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// ConvertFlatMapToNested converts a flat map with dot-notation keys
// into a nested map structure. It also handles array notation like `key[0]`.
// This implementation is a direct port of the logic from JsonMapConvertUtils.java.
func ConvertFlatMapToNested(source map[string]any, ignoreEmpty bool) map[string]any {
	rs := make(map[string]any)
	for key, value := range source {
		valStr, isString := value.(string)
		if value == nil || (isString && strings.TrimSpace(valStr) == "") {
			if ignoreEmpty {
				continue
			}
		}

		if isString {
			trimmedVal := strings.TrimSpace(valStr)
			if strings.HasPrefix(trimmedVal, "[") && strings.HasSuffix(trimmedVal, "]") {
				var arr []any
				// If it looks like a JSON array, try to parse it.
				if err := json.Unmarshal([]byte(trimmedVal), &arr); err == nil {
					value = arr
				} else {
					value = trimmedVal // Use the trimmed string if parsing fails
				}
			} else {
				value = trimmedVal
			}
		}

		parts := strings.Split(key, ".")
		prev := rs
		var prevList []any
		prevIndex := -1
		prevProp := ""

		for i, child := range parts {
			if i > 0 { // In Go, we prepare the next level inside the loop
				if prevIndex >= 0 { // The previous level was an array
					if prevList != nil && len(prevList) > prevIndex {
						oldV := prevList[prevIndex]
						if oldV == nil {
							newMap := make(map[string]any)
							prevList[prevIndex] = newMap
							prev = newMap
						} else if castedMap, ok := oldV.(map[string]any); ok {
							prev = castedMap
						} else {
							// Handle error: expected a map at this array index, but found something else.
							// For now, we overwrite it to match the Java version's behavior.
							newMap := make(map[string]any)
							prevList[prevIndex] = newMap
							prev = newMap
						}
					}
				} else { // The previous level was a map
					if existing, ok := prev[prevProp]; !ok {
						newMap := make(map[string]any)
						prev[prevProp] = newMap
						prev = newMap
					} else if castedMap, ok := existing.(map[string]any); ok {
						prev = castedMap
					} else {
						// Handle error: expected a map for this key, but found something else.
						// Overwrite to match Java logic.
						newMap := make(map[string]any)
						prev[prevProp] = newMap
						prev = newMap
					}
				}
			}

			qt := strings.Index(child, "[")
			ed := -1
			if qt > 0 {
				ed = strings.Index(child[qt+1:], "]")
				if ed > 0 {
					ed += qt + 1
				}
			}

			var prop string
			isArray := false
			if qt > 0 && ed > qt {
				isArray = true
				prop = child[:qt]
				indexStr := strings.TrimSpace(child[qt+1 : ed])
				if indexNum := integerutil.ParseInt(indexStr); indexNum != nil {
					prevIndex = int(*indexNum)
				} else {
					prevIndex = 0 // Default to 0 if parsing fails
				}
			} else {
				prop = child
				prevIndex = -1 // Not an array
			}
			prevProp = prop

			if isArray {
				var list []any
				if existing, ok := prev[prop]; ok {
					if castedList, ok := existing.([]any); ok {
						list = castedList
					}
				}
				if list == nil {
					list = make([]any, 0)
					prev[prop] = list
				}
				// Ensure list is large enough
				for len(list) <= prevIndex {
					list = append(list, nil)
				}
				prev[prop] = list // Re-assign in case list was re-allocated
				prevList = list
			}
		}

		// After the loop, set the final value
		if prevIndex >= 0 {
			if prevList != nil && len(prevList) > prevIndex {
				prevList[prevIndex] = value
			}
		} else {
			prev[prevProp] = value
		}
	}
	return rs
}

// ensureListAndSet is a helper to place a value in a list within a map.
func ensureListAndSet(m *map[string]any, prop string, index int, value any) {
	// Ensure the list exists
	if _, ok := (*m)[prop]; !ok {
		(*m)[prop] = make([]any, 0)
	}
	list, ok := (*m)[prop].([]any)
	if !ok {
		// If the key exists but is not a list, we have a conflict.
		// For simplicity, we overwrite it. A more robust implementation might return an error.
		list = make([]any, 0)
	}
	// Ensure the list is large enough
	for len(list) <= index {
		list = append(list, nil)
	}
	list[index] = value
	(*m)[prop] = list
}

// ensureListAndGetMap is a helper to get a map from within a list, creating it if necessary.
func ensureListAndGetMap(m *map[string]any, prop string, index int) map[string]any {
	// Ensure the list exists
	if _, ok := (*m)[prop]; !ok {
		(*m)[prop] = make([]any, 0)
	}
	list, _ := (*m)[prop].([]any)
	// Ensure the list is large enough
	for len(list) <= index {
		list = append(list, nil)
	}
	// Ensure the map exists at the index
	var childMap map[string]any
	if list[index] == nil {
		childMap = make(map[string]any)
		list[index] = childMap
	} else {
		childMap, _ = list[index].(map[string]any)
	}
	(*m)[prop] = list
	return childMap
}
