// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package cliauth

import (
	"context"
	"strings"

	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/repo"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CliAuthBindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCliAuthBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CliAuthBindLogic {
	return &CliAuthBindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CliAuthBindLogic) CliAuthBind(req *types.CliAuthBindReq) (resp *types.Envelope, err error) {
	if strings.TrimSpace(req.SetupCode) == "" {
		return nil, respx.NewHTTPError(400, "setupCode is required")
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = "cli-" + strings.ToLower(strings.TrimSpace(req.SetupCode))
	}
	permissionInput := strings.TrimSpace(req.Permission)
	if permissionInput == "" {
		permissionInput = repo.APIKeyPermissionFull
	}
	permission := repo.NormalizeAPIKeyPermission(permissionInput)
	data, err := l.svcCtx.App.APIKey.Create(l.ctx, name, logicx.UserID(l.ctx), permission, "cli", "personal", nil)
	if err != nil {
		return nil, logicx.Error(err)
	}
	result := types.CliAuthBindResp{
		ApiKey:        mapString(data, "apiKey"),
		Name:          mapString(data, "name"),
		Permission:    mapString(data, "permission"),
		WorkspaceName: "default",
	}
	storeCliAuthCompleted(req.SetupCode, result)
	return respx.Envelope(result), nil
}

func mapString(data map[string]any, key string) string {
	if data == nil {
		return ""
	}
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}
