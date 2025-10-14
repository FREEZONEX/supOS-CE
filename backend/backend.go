package main

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/svc"
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
	utils.ConfMustLoad("etc/backend.yaml", &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	server.PrintRoutes()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
