package asset

import (
	"mime"
	"strings"
)

// activeContentTypes can execute markup or script when rendered inline by a
// browser. Downloads served from a same-origin URL (the public capability
// route) must force Content-Disposition: attachment for these types so an
// uploaded HTML/SVG/JS payload cannot run in the deployment's origin
// (stored XSS).
var activeContentTypes = map[string]struct{}{
	"text/html":                {},
	"application/xhtml+xml":    {},
	"image/svg+xml":            {},
	"text/xml":                 {},
	"application/xml":          {},
	"application/javascript":   {},
	"application/x-javascript": {},
	"text/javascript":          {},
	"application/ecmascript":   {},
	"text/ecmascript":          {},
}

// IsActiveContentType reports whether contentType is an active (scriptable)
// type that must not be rendered inline. Parameters (charset etc.) are
// ignored; empty or unparseable types are treated as inert. Inert types
// (images, application/pdf, text/plain, ...) stay inline.
func IsActiveContentType(contentType string) bool {
	mediaType := strings.ToLower(strings.TrimSpace(contentType))
	if mediaType == "" {
		return false
	}
	if parsed, _, err := mime.ParseMediaType(mediaType); err == nil {
		mediaType = parsed
	} else if idx := strings.Index(mediaType, ";"); idx >= 0 {
		mediaType = strings.TrimSpace(mediaType[:idx])
	}
	if _, ok := activeContentTypes[mediaType]; ok {
		return true
	}
	// Any *+xml type (xhtml/rss/atom/soap variants) can carry active markup.
	return strings.HasSuffix(mediaType, "+xml")
}
