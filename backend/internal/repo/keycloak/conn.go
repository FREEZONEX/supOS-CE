package keycloak

import (
	"context"
	"fmt"

	"gitee.com/unitedrhino/share/stores"
)

var (
	keycloakDB   *stores.DB
	initErr      error
	isConfigured bool
)

// InitWithDB registers an existing gorm.DB instance (used by service context).
func InitWithDB(db *stores.DB) error {
	if db == nil {
		return fmt.Errorf("nil keycloak db connection")
	}
	keycloakDB = db
	isConfigured = true
	initErr = nil
	return nil
}

// Enabled reports whether the Keycloak database connection is ready.
func Enabled() bool {
	return keycloakDB != nil && initErr == nil
}

// GetConn returns the initialized Keycloak *gorm.DB, optionally binding the provided context.
func GetConn(ctx context.Context) *stores.DB {
	if !Enabled() {
		return nil
	}
	if ctx != nil {
		return keycloakDB.WithContext(ctx)
	}
	return keycloakDB
}
