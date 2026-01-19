package msg_consumer

import (
	"bytes"
	"unicode"

	"github.com/buger/jsonparser"
)

const singleValueKey = "value"

func parseJsonList(payload []byte) (data []map[string]string, err error) {
	switch payload[0] {
	case '{':
		vm := parseMap(payload)
		data = []map[string]string{vm}
	case '[':
		jsonparser.ArrayEach(payload, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			switch dataType {
			case jsonparser.Object:
				data = append(data, parseMap(value))
			case jsonparser.Array:
				rs, pErr := parseJsonList(value)
				if pErr != nil {
					err = pErr
					return
				}
				data = append(data, rs...)
			default:
				vm := make(map[string]string, 2)
				vm[singleValueKey] = b2s(value)
				data = append(data, vm)
			}
		})
	default:
		if payload[0] == '"' {
			payload = bytes.TrimRightFunc(payload[1:], func(r rune) bool {
				return unicode.IsSpace(r) || r == '"'
			})
		}
		vm := make(map[string]string, 2)
		vm[singleValueKey] = string(payload)
		data = append(data, vm)
	}
	return data, err
}
func parseMap(payload []byte) map[string]string {
	vm := make(map[string]string, 8)
	jsonparser.ObjectEach(payload, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		//k := make([]byte, len(key))
		//v := make([]byte, len(value))
		//copy(k, key)
		//copy(v, value)
		vm[b2s(key)] = b2s(value)
		return nil
	})
	return vm
}
