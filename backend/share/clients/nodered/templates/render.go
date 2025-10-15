package templates

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FS is defined in embed.go via //go:embed
// var FS embed.FS

// Load loads a template file content by name.
// Priority: external dir UNS_FLOW_TEMPLATE_DIR/name -> embedded FS.
func Load(name string) (string, error) {
	if dir := os.Getenv("UNS_FLOW_TEMPLATE_DIR"); dir != "" {
		p := filepath.Join(dir, name)
		if b, err := os.ReadFile(p); err == nil {
			return string(b), nil
		}
	}
	b, err := FS.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("load template %s: %w", name, err)
	}
	return string(b), nil
}

var idTokenRe = regexp.MustCompile(`\{\{ID(\d+)\}\}`)
var varTokenRe = regexp.MustCompile(`\{\{([A-Z0-9_]+)\}\}`)
var dollarIdRe = regexp.MustCompile(`\$id_[a-zA-Z0-9_]+`)

// Render replaces placeholders in tpl with variables and generates values for ID tokens.
// Supported tokens:
//   - {{ID1}}, {{ID2}}, ...: will be replaced by values from idGen in order of appearance
//   - {{VAR}}: replaced by vars["VAR"] if present
func Render(tpl string, vars map[string]string, idGen func() string) string {
	// First pass: replace variable tokens
	out := varTokenRe.ReplaceAllStringFunc(tpl, func(m string) string {
		sub := varTokenRe.FindStringSubmatch(m)
		if len(sub) == 2 {
			if v, ok := vars[sub[1]]; ok {
				return v
			}
		}
		return m
	})
	// Second pass: replace IDs in a stable order
	// We scan line by line to keep ordering deterministic
	var sb strings.Builder
	sc := bufio.NewScanner(strings.NewReader(out))
	for sc.Scan() {
		line := sc.Text()
		line = idTokenRe.ReplaceAllStringFunc(line, func(m string) string {
			return idGen()
		})
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	if err := sc.Err(); err == nil {
		s := sb.String()
		if strings.HasSuffix(s, "\n") {
			return s[:len(s)-1]
		}
		return s
	}
	// fallback single shot
	return idTokenRe.ReplaceAllStringFunc(out, func(m string) string { return idGen() })
}

// RenderDollar supports Java-style templates with $ placeholders, e.g.:
//   - $id_inject, $id_model_selector, $id_func, $id_mqtt (stable per token)
//   - $uns_path, $alias_path_topic, $payload, $disabled, $clientid, ...
func RenderDollar(tpl string, vars map[string]string, idGen func() string) string {
	// map each unique $id_* token to a generated id
	ids := map[string]string{}
	out := dollarIdRe.ReplaceAllStringFunc(tpl, func(tok string) string {
		if v, ok := ids[tok]; ok {
			return v
		}
		v := idGen()
		ids[tok] = v
		return v
	})
	// replace variable tokens
	// iterate vars to avoid regex pitfalls with $ in replacement
    for k, v := range vars {
        // ensure newline-safe inside JSON string literals
        vSafe := strings.ReplaceAll(v, "\n", "\\n")
        out = strings.ReplaceAll(out, "$"+k, vSafe)
    }
    // make it a single-line string: drop real newlines/tabs (keep spaces)
    out = strings.ReplaceAll(out, "\r", "")
    out = strings.ReplaceAll(out, "\n", "")
    out = strings.ReplaceAll(out, "\t", "")
    return out
}
