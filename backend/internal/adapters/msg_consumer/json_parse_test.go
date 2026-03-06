package msg_consumer

import (
	"strconv"
	"strings"
	"testing"
	"unicode"
)

func Test_jsoniterMapStr(t *testing.T) {
	jsonStr := []string{`{"timeStamp":1768288504446, "wet":168.931,"qos":0}`, `[{"type":"pg"},{"type":"tsdb"}]`, `3.14159`, `[{"type":"pg"},3.1415]`}
	for _, v := range jsonStr {
		payload := []byte(v)
		data, err := parseJsonList(payload)
		t.Log(err, data)
	}
	strs := []string{`"1.5"`, `3.14159`}
	for _, str := range strs {
		if str[0] == '"' {
			str = strings.TrimRightFunc(str[1:], func(r rune) bool {
				return unicode.IsSpace(r) || r == '"'
			})
		}
		Double, err := strconv.ParseFloat(str, 64)
		t.Log(Double, err, str)
	}

}
