package PathUtil

import "testing"

func TestIsAliasFormatOK(t *testing.T) {
	alias := []string{"123", "bug@34%$ASDFsd@#中文", "Normal123"}
	for _, v := range alias {
		t.Log(v, " Ok ? ", IsAliasFormatOK(v))
	}
}
