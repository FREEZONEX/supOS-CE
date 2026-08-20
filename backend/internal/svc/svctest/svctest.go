// Package svctest provides the retained open-source ServiceContext test fixture.
package svctest

import (
	"context"
	"testing"
	"time"

	"backend/internal/bootstrap"
	"backend/internal/config"
	"backend/internal/contextx"
	"backend/internal/domain/apikey"
	"backend/internal/domain/asset"
	"backend/internal/domain/audit"
	domainauth "backend/internal/domain/auth"
	"backend/internal/domain/dataingest"
	"backend/internal/domain/flow"
	"backend/internal/domain/iam"
	"backend/internal/domain/uns"
	"backend/internal/gateway"
	"backend/internal/permission"
	"backend/internal/svc"
	"backend/internal/testkit"
)

// NewTestServiceContext assembles only the services retained by Tier0 Edge.
func NewTestServiceContext(t *testing.T) *svc.ServiceContext {
	t.Helper()
	store, dsn := testkit.NewTestDB(t)
	cacheClient := testkit.NewTestCache(t)
	ctx := context.Background()

	cfg := testConfig(t, dsn)
	authSvc := domainauth.New(ctx, cfg.Security.JwtSecret)
	auditSvc := audit.New(ctx)
	apiKeySvc := apikey.New(ctx)
	iamSvc := iam.New(ctx)
	permissionEvaluator := permission.New(ctx)
	assetSvc, err := asset.New(ctx, asset.Config{FileRoot: t.TempDir()})
	if err != nil {
		t.Fatalf("svctest: asset.New: %v", err)
	}
	dataIngestSvc, err := dataingest.New(ctx, cfg.DataIngest, "", cacheClient)
	if err != nil {
		t.Fatalf("svctest: dataingest.New: %v", err)
	}
	flowSvc := flow.New(ctx, "", "")
	unsSvc := uns.New(ctx, flowSvc, dataIngestSvc)

	app := &bootstrap.App{
		Config:      cfg,
		Store:       store,
		Cache:       cacheClient,
		Auth:        authSvc,
		Audit:       auditSvc,
		APIKey:      apiKeySvc,
		Asset:       assetSvc,
		DataIngest:  dataIngestSvc,
		Flow:        flowSvc,
		IAM:         iamSvc,
		UNS:         unsSvc,
		Permission:  permissionEvaluator,
		HTTPGateway: gateway.NewHTTPGateway(ctx, cfg.Gateway, authSvc, permissionEvaluator, flowSvc),
		StartedAt:   time.Now().UTC(),
	}
	return svc.NewServiceContext(cfg, app)
}

// CtxWithUser returns a context carrying the retained local-session subject.
func CtxWithUser(userID int64) context.Context {
	return contextx.WithSubject(context.Background(), contextx.Subject{
		UserID:   userID,
		UserName: "tier0",
		AuthType: "session",
	})
}

func testConfig(t *testing.T, dsn string) config.Config {
	t.Helper()
	return config.Config{
		Database: config.DatabaseConf{UnsDbUrl: dsn},
		Security: config.SecurityConf{
			JwtSecret:            "test-jwt-secret",
			InitialAdminPassword: "tier0",
		},
		Storage:    config.StorageConf{FileRoot: t.TempDir()},
		DataIngest: config.DataIngestConf{},
		Gateway:    config.GatewayConf{},
	}
}
