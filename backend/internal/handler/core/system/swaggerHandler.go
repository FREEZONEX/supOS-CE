// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package system

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	openapiV1Prefix          = "/openapi/v1"
	openapiV1SchemaRefPrefix = "#/components/schemas/"
)

var (
	openapiV1SwaggerOnce sync.Once
	openapiV1SwaggerJSON []byte
	openapiV1SwaggerErr  error
)

func OpenapiV1SwaggerUIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Tier0 Edge OpenAPI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    html, body, #swagger-ui { height: 100%; margin: 0; }
    body { background: #fff; }
    .swagger-ui .topbar { display: none; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.addEventListener('load', function () {
      window.ui = SwaggerUIBundle({
        url: '/openapi/v1/swagger.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        persistAuthorization: true,
        displayRequestDuration: true
      });
    });
  </script>
</body>
</html>`))
	}
}

func OpenapiV1SwaggerJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := openapiV1NormalizedSwaggerJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(data)
	}
}

func SwaggerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readSwaggerJSON()
		if err != nil {
			http.Error(w, "swagger.json not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(data)
	}
}

func readSwaggerJSON() ([]byte, error) {
	candidates := []string{
		filepath.Join("http", "swagger.json"),
		filepath.Join("backend", "http", "swagger.json"),
		filepath.Join("/app", "http", "swagger.json"),
	}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err == nil {
			return data, nil
		}
	}
	return nil, os.ErrNotExist
}

func openapiV1NormalizedSwaggerJSON() ([]byte, error) {
	openapiV1SwaggerOnce.Do(func() {
		raw, err := readSwaggerJSON()
		if err != nil {
			openapiV1SwaggerErr = err
			return
		}
		var spec map[string]any
		if err := json.Unmarshal(raw, &spec); err != nil {
			openapiV1SwaggerErr = err
			return
		}

		spec["info"] = map[string]any{
			"title":       "Tier0 Edge OpenAPI v1",
			"description": "Tier0 Edge OpenAPI",
			"version":     "v1",
		}
		spec["servers"] = []any{
			map[string]any{
				"url":         openapiV1Prefix,
				"description": "OpenAPI v1",
			},
		}
		spec["tags"] = []any{
			map[string]any{
				"name":        "openapi",
				"description": "OpenAPI v1",
			},
		}
		paths := normalizeOpenAPIV1Paths(spec["paths"])
		spec["paths"] = paths
		pruneOpenAPIV1Schemas(spec, paths)

		components, ok := spec["components"].(map[string]any)
		if !ok {
			components = map[string]any{}
			spec["components"] = components
		}
		components["securitySchemes"] = map[string]any{
			"ApiKeyAuth": map[string]any{
				"type":        "apiKey",
				"in":          "header",
				"name":        "x-api-key",
				"description": "OpenAPI API Key, for example sk-service-xxx or sk-agent-xxx",
			},
		}
		spec["security"] = []any{
			map[string]any{"ApiKeyAuth": []any{}},
		}
		applyOpenAPIRequestDescriptions(spec)

		data, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			openapiV1SwaggerErr = err
			return
		}
		openapiV1SwaggerJSON = data
	})
	return openapiV1SwaggerJSON, openapiV1SwaggerErr
}

func normalizeOpenAPIV1Paths(raw any) map[string]any {
	paths, ok := raw.(map[string]any)
	if !ok {
		return map[string]any{}
	}
	out := make(map[string]any)
	for path, item := range paths {
		if !strings.HasPrefix(path, openapiV1Prefix+"/") {
			continue
		}
		normalizedPath := strings.TrimPrefix(path, openapiV1Prefix)
		if normalizedPath == "" {
			normalizedPath = "/"
		}
		if !isOpenAPIV1DocPath(normalizedPath) {
			continue
		}
		normalizeOpenAPIV1PathItem(normalizedPath, item)
		out[normalizedPath] = item
	}
	return out
}

func isOpenAPIV1DocPath(path string) bool {
	return path == "/info" ||
		strings.HasPrefix(path, "/assets") ||
		strings.HasPrefix(path, "/auth/") ||
		strings.HasPrefix(path, "/flow/") ||
		strings.HasPrefix(path, "/flows") ||
		strings.HasPrefix(path, "/uns/")
}

func normalizeOpenAPIV1PathItem(path string, item any) {
	pathItem, ok := item.(map[string]any)
	if !ok {
		return
	}
	for _, operation := range pathItem {
		operationMap, ok := operation.(map[string]any)
		if !ok {
			continue
		}
		operationMap["tags"] = []any{"openapi"}
		if path == "/info" {
			operationMap["security"] = []any{}
		}
	}
}

func pruneOpenAPIV1Schemas(spec map[string]any, paths map[string]any) {
	components, ok := spec["components"].(map[string]any)
	if !ok {
		return
	}
	schemas, ok := components["schemas"].(map[string]any)
	if !ok {
		return
	}
	used := map[string]struct{}{}
	collectOpenAPIV1SchemaRefs(paths, schemas, used)
	for name, schema := range schemas {
		if isOpenAPIV1SchemaName(name) {
			used[name] = struct{}{}
			collectOpenAPIV1SchemaRefs(schema, schemas, used)
		}
	}
	if len(used) == 0 {
		delete(components, "schemas")
		return
	}
	pruned := make(map[string]any, len(used))
	for name, schema := range schemas {
		if _, ok := used[name]; ok {
			pruned[name] = schema
		}
	}
	components["schemas"] = pruned
}

func collectOpenAPIV1SchemaRefs(value any, schemas map[string]any, used map[string]struct{}) {
	switch typed := value.(type) {
	case map[string]any:
		if ref, ok := typed["$ref"].(string); ok && strings.HasPrefix(ref, openapiV1SchemaRefPrefix) {
			name := strings.TrimPrefix(ref, openapiV1SchemaRefPrefix)
			if _, exists := used[name]; !exists {
				used[name] = struct{}{}
				collectOpenAPIV1SchemaRefs(schemas[name], schemas, used)
			}
		}
		for _, child := range typed {
			collectOpenAPIV1SchemaRefs(child, schemas, used)
		}
	case []any:
		for _, child := range typed {
			collectOpenAPIV1SchemaRefs(child, schemas, used)
		}
	}
}

func isOpenAPIV1SchemaName(name string) bool {
	if strings.HasPrefix(name, "Openapi") {
		return true
	}
	switch name {
	case "Envelope",
		"AuthWhoamiResp",
		"InfoResp",
		"ReadResp",
		"WriteResp",
		"BrowseResp",
		"SearchResp",
		"HistoryResp",
		"CreateResp",
		"AssetUploadResp",
		"UnsAttachmentInfo",
		"UnsAttachmentUploadResp",
		"UnsAttachmentListResp",
		"UnsBindFlowResp":
		return true
	default:
		return false
	}
}
