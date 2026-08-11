// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package cliauth

import (
	"context"
	"strings"

	respx "backend/internal/httpx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CliAuthStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCliAuthStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CliAuthStatusLogic {
	return &CliAuthStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CliAuthStatusLogic) CliAuthStatus(req *types.CliAuthStatusReq) (resp *types.Envelope, err error) {
	if strings.TrimSpace(req.SetupCode) == "" {
		return nil, respx.NewHTTPError(400, "setupCode is required")
	}
	return &types.Envelope{Code: 200, Msg: "", Data: loadCliAuthStatus(req.SetupCode)}, nil
}
