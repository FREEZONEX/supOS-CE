package main

import (
	"backend/internal/config"
	"backend/internal/handler"
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
	utils.ConfMustLoad(confFile, &c)

	/* 下面是使用示例
	msg := i18ns.LocalizeMsg("nodered.protocol.unsupported", "vewwrfw3")
	logx.Info(msg) // 输出:  Unsupported protocol: vewwrfw3.
	*/

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	server.PrintRoutes()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	spring.RegisterBean[*svc.ServiceContext](ctx)
	spring.RefreshBeanContext()
	server.Start()
}
