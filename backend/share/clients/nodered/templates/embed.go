package templates

import "embed"

// FS embeds all template files in this directory.
//go:embed *.tpl
var FS embed.FS

