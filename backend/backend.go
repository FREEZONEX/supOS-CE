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
		confFile = "etc/backend-local.yaml"
	}
	utils.ConfMustLoad(confFile, &c)
	server := rest.MustNewServer(c.RestConf, rest.WithFileServer("/files/", http.Dir("/app/go-edge")))
	defer server.Stop()

	if lv := c.LoggerLevel; lv != "" {
		system.SetLogLevel(lv)
	}

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	handler.RegisterExtHandlers(server, ctx)
	server.PrintRoutes()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	spring.RegisterBean[*svc.ServiceContext](ctx)
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
