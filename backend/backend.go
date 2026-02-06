package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "backend/internal/adapters/grafana"
	_ "backend/internal/adapters/msg_consumer" // 手动导入 adapter
	"backend/internal/common/event"
	"backend/internal/config"
	"backend/internal/handler"
	service "backend/internal/logic/supos/appkey/service"
	_ "backend/internal/logic/supos/uns/dashboard/service"
	"backend/internal/logic/supos/uns/system"
	_ "backend/internal/logic/supos/uns/topology/service" // 导入触发 init() 注册
	_ "backend/internal/logic/supos/uns/uns/service"
	"backend/internal/svc"
	"backend/share/spring"

	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	_ "github.com/zeromicro/go-zero/core/proc" //开启pprof采集 https://mp.weixin.qq.com/s/yYFM3YyBbOia3qah3eRVQA
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	defer utils.Recover(context.Background())
	flag.Parse()
	logx.DisableStat()
	var c config.Config
	var confFile = "etc/backend.yaml"
	if info, er := os.Stat("../deploy/"); er == nil && info.IsDir() {
		confFile = "etc/backend-dev.yaml"
	}
	utils.ConfMustLoad(confFile, &c)
	// 提供 Swagger UI 静态资源
	var opts []rest.RunOption
	opts = append(opts, rest.WithFileServer("/files/", http.Dir("/app/go-edge")))

	// 判断是否是本地开发环境
	if info, er := os.Stat("../deploy/"); er == nil && info.IsDir() {
		// 本地开发环境
		opts = append(opts, rest.WithFileServer("/swagger-ui", http.Dir("swagger/dist")))
	} else {
		// 生产环境
		opts = append(opts, rest.WithFileServer("/swagger-ui", http.Dir("/app/swagger/dist")))
	}

	server := rest.MustNewServer(c.RestConf, opts...)
	defer server.Stop()

	system.SetLogLevel(c.Log.Level)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	handler.RegisterExtHandlers(server, ctx)
	server.PrintRoutes()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	spring.RegisterBean[*svc.ServiceContext](ctx)

	// 初始化 AppKeyService
	appKeyService := spring.GetBean[*service.AppKeyService]()
	if err := appKeyService.Init(c); err != nil {
		logx.Errorf("Failed to initialize AppKeyService: %v", err)
		os.Exit(-1)
	}

	spring.RefreshBeanContext()
	_ = spring.PublishEvent(&event.ContextRefreshedEvent{SvcContext: ctx})
	defer func() {
		_ = spring.PublishEvent(&event.ContextClosedEvent{SvcContext: ctx})
	}()
	fmt.Printf("Started server at %s:%d...\n", c.Host, c.Port)
	server.StartWithOpts(withSwaggerBinaryContentType())
}

// withSwaggerBinaryContentType forces the swagger yaml downloads to use application/octet-stream.
func withSwaggerBinaryContentType() rest.StartOption {
	return func(svr *http.Server) {
		if svr == nil || svr.Handler == nil {
			return
		}

		handler := svr.Handler
		svr.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isSwaggerBinaryFile(r.URL.Path) {
				w.Header().Set("Content-Type", "application/octet-stream")
			}
			handler.ServeHTTP(w, r)
		})
	}
}

func isSwaggerBinaryFile(path string) bool {
	switch strings.ToLower(path) {
	case "/files/system/resource/swagger/supos.yaml", "/files/system/resource/swagger/supos-en.yaml":
		return true
	default:
		return false
	}
}
