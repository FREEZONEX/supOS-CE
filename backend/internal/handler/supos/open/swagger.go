package open

import (
	"net/http"
	"path/filepath"
	"runtime"
)

// SwaggerYamlHandler 提供 swagger.yaml 文件的访问
func SwaggerYamlHandler(w http.ResponseWriter, r *http.Request) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	http.ServeFile(w, r, filepath.Join(dir, "Tier0.openapi.yaml"))
}

// SwaggerUIHandler 提供 Swagger UI 静态资源的访问
func SwaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	// Swagger UI 静态文件位于项目根目录的 dist 文件夹下
	swaggerUIPath := filepath.Join(dir, "../../../../dist")

	// 处理根路径请求，返回 index.html
	if r.URL.Path == "/swagger-ui/" || r.URL.Path == "/swagger-ui" {
		http.ServeFile(w, r, filepath.Join(swaggerUIPath, "index.html"))
		return
	}

	// 处理其他静态资源请求
	http.StripPrefix("/swagger-ui", http.FileServer(http.Dir(swaggerUIPath))).ServeHTTP(w, r)
}
