package PathUtil

import (
	"regexp"
	"testing"
)

func TestIsAliasFormatOK(t *testing.T) {
	alias := []string{"123", "4AIgc", "bug@34%$ASDFsd@#中文", "Normal123", "aigc.comon"}
	fieldNamePattern = regexp.MustCompile(`[A-Za-z0-9_]+$`)
	for _, v := range alias {
		t.Log(v, " Ok ? ", fieldNamePattern.MatchString(v))
	}
}
