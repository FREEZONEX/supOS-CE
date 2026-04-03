package oauth

import (
	"context"
	"strings"
	"time"

	"backend/internal/common/constants"
	iamrepo "backend/internal/repo/iam"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type baseOAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func newBaseOAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) baseOAuthLogic {
	return baseOAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *baseOAuthLogic) repo() (*iamrepo.AuthRepo, error) {
	return iamrepo.NewAuthRepo(l.ctx)
}

func (l *baseOAuthLogic) getSession(sessionID string) (*relationDB.IamSession, error) {
	repo, err := l.repo()
	if err != nil || repo == nil {
		return nil, err
	}
	return repo.GetSession(strings.TrimSpace(sessionID))
}

func (l *baseOAuthLogic) isSessionActive(session *relationDB.IamSession) bool {
	if session == nil || session.RevokedAt != nil {
		return false
	}
	return session.ExpiredAt.After(time.Now())
}

func (l *baseOAuthLogic) touchSession(sessionID string) error {
	repo, err := l.repo()
	if err != nil || repo == nil {
		return err
	}
	now := time.Now()
	return repo.TouchSession(strings.TrimSpace(sessionID), now, now.Add(time.Duration(constants.TokenMaxAge)*time.Second))
}

func (l *baseOAuthLogic) revokeSession(sessionID string) {
	repo, err := l.repo()
	if err != nil || repo == nil {
		return
	}
	if revokeErr := repo.RevokeSession(strings.TrimSpace(sessionID), time.Now()); revokeErr != nil {
		l.Errorf("revoke oauth session failed: %v", revokeErr)
	}
}
