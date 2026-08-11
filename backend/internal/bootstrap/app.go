package bootstrap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/domain/apikey"
	"backend/internal/domain/asset"
	"backend/internal/domain/audit"
	"backend/internal/domain/auth"
	"backend/internal/domain/dataingest"
	"backend/internal/domain/flow"
	"backend/internal/domain/iam"
	"backend/internal/domain/uns"
	"backend/internal/gateway"
	"backend/internal/infra/cache"
	infraoutbox "backend/internal/infra/outbox"
	"backend/internal/permission"
	"backend/internal/repo"
)

type App struct {
	Config       config.Config
	Store        *repo.Store
	Cache        *cache.Client
	Auth         *auth.Service
	Audit        *audit.Service
	APIKey       *apikey.Service
	Asset        *asset.Service
	DataIngest   *dataingest.Service
	Flow         *flow.Service
	IAM          *iam.Service
	UNS          *uns.Service
	Permission   *permission.Evaluator
	HTTPGateway  *gateway.HTTPGateway
	OutboxWorker *infraoutbox.Worker
	StartedAt    time.Time
}

func New(ctx context.Context, c config.Config) (*App, error) {
	if strings.TrimSpace(c.Security.JwtSecret) == "" {
		return nil, fmt.Errorf("JWT secret is required")
	}
	if strings.TrimSpace(c.Security.InitialAdminPassword) == "" {
		return nil, fmt.Errorf("initial admin password is required")
	}
	store, err := repo.Open(ctx, c.Database)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := store.Migrate(ctx); err != nil {
		store.Close()
		return nil, fmt.Errorf("migrate db: %w", err)
	}
	if err := store.Seed(ctx, c.Security, c.Gateway); err != nil {
		store.Close()
		return nil, fmt.Errorf("seed db: %w", err)
	}
	cacheClient := cache.Open(c.Redis)
	app := &App{
		Config:     c,
		Store:      store,
		Cache:      cacheClient,
		Permission: permission.New(ctx),
		StartedAt:  time.Now().UTC(),
	}
	app.Auth = auth.New(ctx, c.Security.JwtSecret)
	app.Audit = audit.New(ctx)
	app.APIKey = apikey.New(ctx)
	app.Asset, err = asset.New(ctx, asset.Config{FileRoot: c.Storage.FileRoot})
	if err != nil {
		_ = cacheClient.Close()
		store.Close()
		return nil, fmt.Errorf("open filestore: %w", err)
	}
	app.DataIngest, err = dataingest.New(ctx, c.DataIngest, c.Database.SinkDbUrl, cacheClient)
	if err != nil {
		_ = cacheClient.Close()
		store.Close()
		return nil, fmt.Errorf("init dataingest: %w", err)
	}
	app.Flow = flow.New(ctx, c.Gateway.SourceFlowUrl, c.Gateway.EventFlowUrl)
	app.IAM = iam.New(ctx)
	app.UNS = uns.New(ctx, app.Flow, app.DataIngest)
	app.HTTPGateway = gateway.NewHTTPGateway(ctx, c.Gateway, app.Auth, app.Permission, app.Flow)
	registerEventHandlers(app)
	app.OutboxWorker = newOutboxWorker(app)
	if err := app.DataIngest.Start(ctx); err != nil {
		app.Close()
		return nil, fmt.Errorf("start dataingest: %w", err)
	}
	app.StartOutboxWorkerAsync()
	return app, nil
}

func (a *App) StartOutboxWorkerAsync() {
	if a == nil || a.OutboxWorker == nil {
		return
	}
	a.OutboxWorker.Start(context.Background())
}

func (a *App) Close() {
	if a == nil {
		return
	}
	if a.OutboxWorker != nil {
		a.OutboxWorker.Stop()
	}
	if a.Cache != nil {
		_ = a.Cache.Close()
	}
	if a.DataIngest != nil {
		a.DataIngest.Close()
	}
	if a.Store != nil {
		a.Store.Close()
	}
}

func (a *App) Readiness(ctx context.Context) map[string]string {
	out := map[string]string{"db": "ok", "redis": "ok"}
	if err := a.Store.Ping(ctx); err != nil {
		out["db"] = err.Error()
	}
	if err := a.Cache.Ping(ctx); err != nil {
		out["redis"] = err.Error()
	}
	return out
}

func (a *App) Ready(ctx context.Context) error {
	if err := a.Store.Ping(ctx); err != nil {
		return err
	}
	return a.Cache.Ping(ctx)
}
