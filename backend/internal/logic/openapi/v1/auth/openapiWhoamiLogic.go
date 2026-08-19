// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package auth

import (
	"context"
	"sort"
	"strings"

	"backend/internal/contextx"
	domaincommon "backend/internal/domain/common"
	respx "backend/internal/httpx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiWhoamiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiWhoamiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiWhoamiLogic {
	return &OpenapiWhoamiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiWhoamiLogic) OpenapiWhoami() (resp *types.Envelope, err error) {
	subject, _ := contextx.SubjectFrom(l.ctx)
	permissions := make([]string, 0, len(subject.ResourceKeys))
	for key := range subject.ResourceKeys {
		permissions = append(permissions, key)
	}
	sort.Strings(permissions)
	roles := compactStrings(subject.AuthType, subject.KeyType)
	return respx.Envelope(types.AuthWhoamiResp{
		UserID:        subject.UserID,
		UserName:      subject.UserName,
		Email:         subject.Email,
		WorkspaceID:   domaincommon.PlatformWorkspaceID,
		WorkspaceName: "",
		ApiKeyName:    subject.APIKeyName,
		KeyPrefix:     subject.APIKeyPrefix,
		Permissions:   permissions,
		Roles:         roles,
		KeyType:       subject.KeyType,
	}), nil
}

func compactStrings(values ...string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
