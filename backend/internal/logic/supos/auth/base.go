package auth

import (
	"context"
	"os"
	"strings"
	"time"

	iamrepo "backend/internal/repo/iam"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type baseAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func newBaseAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) baseAuthLogic {
	return baseAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *baseAuthLogic) removeSession(sessionKey string) {
	if sessionKey == "" {
		return
	}
	if repo, err := l.iamRepo(); err == nil && repo != nil {
		if revokeErr := repo.RevokeSession(sessionKey, time.Now()); revokeErr != nil {
			l.Errorf("revoke iam session failed: %v", revokeErr)
		}
	}
}

func (l *baseAuthLogic) authDisabled() bool {
	return strings.EqualFold(os.Getenv("SYS_OS_AUTH_ENABLE"), "false")
}

func (l *baseAuthLogic) iamRepo() (*iamrepo.AuthRepo, error) {
	return iamrepo.NewAuthRepo(l.ctx)
}
