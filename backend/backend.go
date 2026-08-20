// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package main

import (
	"backend/internal/bootstrap"
	"backend/internal/config"
	"backend/internal/gateway"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repo"
	"backend/internal/svc"
	"context"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	_ "github.com/zeromicro/go-zero/core/proc" //开启pprof采集 https://mp.weixin.qq.com/s/yYFM3YyBbOia3qah3eRVQA

	"github.com/zeromicro/go-zero/rest"
)

func main() {
	migrateOnly := flag.Bool("migrate-only", false, "run database migrations and exit")
	flag.Parse()
	logx.DisableStat()
	var c config.Config
	conf.MustLoad("etc/backend.yaml", &c)
	c.Normalize()

	if *migrateOnly {
		if err := runMigrateOnly(context.Background(), c); err != nil {
			panic(err)
		}
		fmt.Println("database migration complete")
		return
	}

	app, err := bootstrap.New(context.Background(), c)
	if err != nil {
		panic(err)
	}
	defer app.Close()

	server := newRestServer(c, app)
	defer server.Stop()


	server.PrintRoutes()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.StartWithOpts(gateway.WithHTTPGateway(app.HTTPGateway))
}

func runMigrateOnly(ctx context.Context, c config.Config) error {
	store, err := repo.Open(ctx, c.Database)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer store.Close()
	if err := store.Migrate(ctx); err != nil {
		return fmt.Errorf("migrate db: %w", err)
	}
	return nil
}

func newRestServer(c config.Config, app *bootstrap.App) *rest.Server {
	server := rest.MustNewServer(c.RestConf)
	ctx := svc.NewServiceContext(c, app)
	server.Use(middleware.NewAuditLogMiddleware(app.Audit).Handle)
	handler.RegisterHandlers(server, ctx)
	handler.RegisterExtraHandlers(server)
	return server
}
