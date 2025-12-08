package main

import (
	_ "backend/internal/adapters/grafana"
	_ "backend/internal/adapters/msg_consumer" // 手动导入 adapter
	"backend/internal/common/event"
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/logic/supos/uns/system"
	_ "backend/internal/logic/supos/uns/topology/service" // 导入触发 init() 注册
	_ "backend/internal/logic/supos/uns/uns/service"
	"os"

	"backend/internal/svc"
	"backend/share/spring"
	"context"
	"flag"
	"fmt"

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
	c.RestConf.MaxBytes = max(c.RestConf.MaxBytes, 1<<30) //http body最大限制最少1G

	server := rest.MustNewServer(c.RestConf)
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
	server.Start()
}
