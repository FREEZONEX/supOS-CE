package iam

import (
	"errors"
	"strings"
	"time"

	"backend/internal/repo/relationDB"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *AuthRepo) GetOAuthClientByClientID(clientID string) (*relationDB.IamOAuthClient, error) {
	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		return nil, nil
	}

	var client relationDB.IamOAuthClient
	err := r.db.Where("LOWER(client_id) = LOWER(?) AND enabled = ?", clientID, true).First(&client).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &client, nil
}

func (r *AuthRepo) CreateAuthorizationCode(code *relationDB.IamOAuthAuthorizationCode) error {
	if code == nil {
		return nil
	}
	return stores.ErrFmt(r.db.Create(code).Error)
}

func (r *AuthRepo) ConsumeAuthorizationCode(codeValue, clientID string, now time.Time) (*relationDB.IamOAuthAuthorizationCode, error) {
	codeValue = strings.TrimSpace(codeValue)
	clientID = strings.TrimSpace(clientID)
	if codeValue == "" || clientID == "" {
		return nil, nil
	}

	var code relationDB.IamOAuthAuthorizationCode
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("code = ? AND LOWER(client_id) = LOWER(?)", codeValue, clientID).
			Take(&code).Error; err != nil {
			return err
		}
		if code.UsedAt != nil || !code.ExpiredAt.After(now) {
			return gorm.ErrRecordNotFound
		}
		return tx.Model(&relationDB.IamOAuthAuthorizationCode{}).
			Where("code = ? AND used_at IS NULL", codeValue).
			Update("used_at", now).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	code.UsedAt = &now
	return &code, nil
}

func (r *AuthRepo) CreateAccessToken(token *relationDB.IamOAuthAccessToken) error {
	if token == nil {
		return nil
	}
	return stores.ErrFmt(r.db.Create(token).Error)
}

func (r *AuthRepo) GetAccessToken(accessToken string) (*relationDB.IamOAuthAccessToken, error) {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return nil, nil
	}

	var token relationDB.IamOAuthAccessToken
	err := r.db.Where("access_token = ?", accessToken).First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &token, nil
}

func (r *AuthRepo) RevokeAccessToken(accessToken string, when time.Time) error {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return nil
	}
	return stores.ErrFmt(r.db.Model(&relationDB.IamOAuthAccessToken{}).
		Where("access_token = ?", accessToken).
		Update("revoked_at", when).Error)
}
