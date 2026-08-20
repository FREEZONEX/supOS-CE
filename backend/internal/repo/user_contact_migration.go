package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func (s *Store) migrateUserContacts(ctx context.Context) error {
	if s == nil || s.commonDB == nil {
		return fmt.Errorf("migrate user contacts: database is not initialized")
	}
	return s.commonDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("LOCK TABLE sys_user_info IN ACCESS EXCLUSIVE MODE").Error; err != nil {
			return fmt.Errorf("lock sys_user_info for contact migration: %w", err)
		}
		type contactRow struct {
			UserID int64  `gorm:"column:user_id"`
			Email  string `gorm:"column:email"`
			Phone  string `gorm:"column:phone"`
		}
		var rows []contactRow
		if err := tx.Table("sys_user_info").Select("user_id, email, phone").Order("user_id").Scan(&rows).Error; err != nil {
			return fmt.Errorf("load user contacts for migration: %w", err)
		}
		for _, row := range rows {
			email, emailChanged, err := migrateStoredUserContact(userContactFieldEmail, row.Email)
			if err != nil {
				return fmt.Errorf("migrate user %d email: %w", row.UserID, err)
			}
			phone, phoneChanged, err := migrateStoredUserContact(userContactFieldPhone, row.Phone)
			if err != nil {
				return fmt.Errorf("migrate user %d phone: %w", row.UserID, err)
			}
			if !emailChanged && !phoneChanged {
				continue
			}
			updates := map[string]any{}
			if emailChanged {
				updates["email"] = email
			}
			if phoneChanged {
				updates["phone"] = phone
			}
			if err := tx.Table("sys_user_info").Where("user_id = ?", row.UserID).Updates(updates).Error; err != nil {
				return fmt.Errorf("update encrypted contacts for user %d: %w", row.UserID, err)
			}
		}
		return nil
	})
}

func migrateStoredUserContact(field, stored string) (string, bool, error) {
	if stored == "" {
		return "", false, nil
	}
	if userContactIsEncrypted(stored) {
		if _, err := decryptUserContact(field, stored); err != nil {
			return "", false, err
		}
		return stored, false, nil
	}
	encrypted, err := encryptUserContact(field, stored)
	if err != nil {
		return "", false, err
	}
	return encrypted, true, nil
}
