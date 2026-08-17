package repo

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

const userEmailWriteLockKey int64 = 0x5449455230504949

func lockUserEmailWrites(tx *gorm.DB) error {
	if tx == nil {
		return ErrInvalidArgument
	}
	return tx.Exec("SELECT pg_advisory_xact_lock(?)", userEmailWriteLockKey).Error
}

func ensureUserEmailAvailable(tx *gorm.DB, email string, excludeUserID int64) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil
	}
	type userEmailRow struct {
		UserID int64  `gorm:"column:user_id"`
		Email  string `gorm:"column:email"`
	}
	var rows []userEmailRow
	if err := tx.Table("sys_user_info").
		Select("user_id, email").
		Where("deleted_time = 0 AND user_id <> ? AND email <> ''", excludeUserID).
		Scan(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		plaintext, err := decryptUserContact(userContactFieldEmail, row.Email)
		if err != nil {
			return fmt.Errorf("decrypt email for user %d: %w", row.UserID, err)
		}
		if strings.TrimSpace(plaintext) == email {
			return ErrUserEmailDuplicate
		}
	}
	return nil
}

func decryptUserContacts(userID int64, email, phone string) (string, string, error) {
	plaintextEmail, plaintextPhone, err := decryptUserContactPair(email, phone)
	if err != nil {
		return "", "", fmt.Errorf("decrypt contacts for user %d: %w", userID, err)
	}
	return plaintextEmail, plaintextPhone, nil
}
