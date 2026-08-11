// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package platform

import (
	"context"

	domainiam "backend/internal/domain/iam"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiPlatformGetMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiPlatformGetMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiPlatformGetMembersLogic {
	return &OpenapiPlatformGetMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiPlatformGetMembersLogic) OpenapiPlatformGetMembers(req *types.PlatformGetMembersReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.IAM.PlatformMembers(l.ctx, domainiam.PlatformMembersQuery{
		Keyword:        req.Keyword,
		RoleKey:        req.RoleKey,
		Roles:          req.Roles,
		Statuses:       req.Statuses,
		UpdatedAtStart: req.UpdatedAtStart,
		UpdatedAtEnd:   req.UpdatedAtEnd,
		Page:           req.Page,
		Size:           req.Size,
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
