package auth

import (
	"backend/internal/common/constants"
	"backend/internal/common/vo"
	"backend/internal/repo/relationDB"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (l *baseAuthLogic) getIAMSession(sessionKey string) (*relationDB.IamSession, error) {
	repo, err := l.iamRepo()
	if err != nil || repo == nil {
		return nil, err
	}
	return repo.GetSession(strings.TrimSpace(sessionKey))
}

func (l *baseAuthLogic) isIAMSessionActive(session *relationDB.IamSession) bool {
	if session == nil {
		return false
	}
	if session.RevokedAt != nil {
		return false
	}
	return session.ExpiredAt.After(time.Now())
}

func (l *baseAuthLogic) buildIAMUserInfo(userID string) (*vo.UserInfoVo, error) {
	repo, err := l.iamRepo()
	if err != nil || repo == nil {
		return nil, err
	}
	return repo.BuildUserInfo(l.ctx, strings.TrimSpace(userID), constants.DefaultHomepage)
}

func (l *baseAuthLogic) createIAMSession(userID string) (string, error) {
	repo, err := l.iamRepo()
	if err != nil || repo == nil {
		return "", err
	}
	now := time.Now()
	sessionID := uuid.NewString()
	err = repo.CreateSession(&relationDB.IamSession{
		ID:            sessionID,
		UserID:        strings.TrimSpace(userID),
		ExpiredAt:     now.Add(time.Duration(constants.TokenMaxAge) * time.Second),
		LastAccessAt:  now,
		PasswordReady: true,
		CreatedAt:     now,
	})
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (l *baseAuthLogic) touchIAMSession(sessionKey string) error {
	repo, err := l.iamRepo()
	if err != nil || repo == nil {
		return err
	}
	now := time.Now()
	return repo.TouchSession(strings.TrimSpace(sessionKey), now, now.Add(time.Duration(constants.TokenMaxAge)*time.Second))
}

func (l *baseAuthLogic) removeIAMSession(sessionKey string) error {
	repo, err := l.iamRepo()
	if err != nil || repo == nil {
		return err
	}
	return repo.RevokeSession(strings.TrimSpace(sessionKey), time.Now())
}
