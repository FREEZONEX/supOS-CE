package repo

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"strings"

	"backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed migrations/schema.sql migrations/sink_schema.sql migrations/versions/common/*.sql migrations/versions/sink/*.sql
var migrationFS embed.FS

// ErrNotFound is returned when a single-row lookup or a conditional
// update/delete matches no row. It aliases gorm.ErrRecordNotFound so callers
// may use errors.Is(err, repo.ErrNotFound) regardless of the driver.
var ErrNotFound = gorm.ErrRecordNotFound

var ErrInvalidArgument = errors.New("invalid argument")
var ErrConflict = errors.New("conflict")

var ErrDuplicate = errors.New("common.duplicateResource")
var ErrUserAccountDuplicate = errors.New("UserManagement.accountDuplicate")
var ErrUserEmailDuplicate = errors.New("UserManagement.emailDuplicate")

var (
	commonConn *gorm.DB
	unsConn    *gorm.DB
)

// Store owns the gorm connection lifecycle and schema bootstrap only.
// Business data access is done through explicit NewXxxRepo(ctx) instances.
type Store struct {
	commonDB                 *gorm.DB
	unsDB                    *gorm.DB
	timeseriesRetentionYears int
}

func Open(ctx context.Context, c config.DatabaseConf) (*Store, error) {
	commonDB, err := openPostgres(c.UnsDbUrl)
	if err != nil {
		return nil, err
	}
	unsDSN := strings.TrimSpace(c.SinkDbUrl)
	if unsDSN == "" {
		unsDSN = c.UnsDbUrl
	}
	unsDB := commonDB
	if strings.TrimSpace(unsDSN) != strings.TrimSpace(c.UnsDbUrl) {
		unsDB, err = openPostgres(unsDSN)
		if err != nil {
			closeDB(commonDB)
			return nil, err
		}
	}
	commonConn = commonDB
	unsConn = unsDB
	return &Store{
		commonDB:                 commonDB,
		unsDB:                    unsDB,
		timeseriesRetentionYears: config.NormalizeTimeseriesRetentionYears(c.TimeseriesRetentionYears),
	}, nil
}

func openPostgres(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	return db, nil
}

func GetCommonConn(in any) *gorm.DB {
	if store, ok := in.(*Store); ok {
		return store.CommonDB()
	}
	return getConn(in, commonConn)
}

func GetUnsConn(in any) *gorm.DB {
	if store, ok := in.(*Store); ok {
		return store.UnsDB()
	}
	return getConn(in, unsConn)
}

func getConn(in any, fallback *gorm.DB) *gorm.DB {
	if db, ok := in.(*gorm.DB); ok {
		return db
	}
	ctx, ok := in.(context.Context)
	if !ok {
		if in == nil {
			ctx = context.Background()
		} else if fallback != nil {
			conn := fallback.Session(&gorm.Session{})
			conn.Error = fmt.Errorf("%w: repository input must be context.Context or *gorm.DB", ErrInvalidArgument)
			return conn
		} else {
			return nil
		}
	}
	if fallback == nil {
		return nil
	}
	return fallback.WithContext(ctx)
}

func (s *Store) DB() *gorm.DB {
	return s.CommonDB()
}

func (s *Store) CommonDB() *gorm.DB {
	if s == nil {
		return nil
	}
	return s.commonDB
}

func (s *Store) UnsDB() *gorm.DB {
	if s == nil {
		return nil
	}
	if s.unsDB != nil {
		return s.unsDB
	}
	return s.commonDB
}

func (s *Store) Close() {
	if s == nil {
		return
	}
	if s.unsDB != nil && s.unsDB != s.commonDB {
		closeDB(s.unsDB)
	}
	if s.commonDB != nil {
		closeDB(s.commonDB)
	}
	if commonConn == s.commonDB {
		commonConn = nil
	}
	if unsConn == s.unsDB {
		unsConn = nil
	}
}

func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.commonDB == nil {
		return errors.New("database not initialized")
	}
	if err := pingDB(ctx, s.commonDB); err != nil {
		return err
	}
	if s.unsDB != nil && s.unsDB != s.commonDB {
		if err := pingDB(ctx, s.unsDB); err != nil {
			return fmt.Errorf("uns sink database: %w", err)
		}
	}
	return nil
}

func (s *Store) Migrate(ctx context.Context) error {
	sqlBytes, err := migrationFS.ReadFile("migrations/schema.sql")
	if err != nil {
		return err
	}
	if err := s.commonDB.WithContext(ctx).Exec(string(sqlBytes)).Error; err != nil {
		return err
	}
	sinkSQLBytes, err := migrationFS.ReadFile("migrations/sink_schema.sql")
	if err != nil {
		return err
	}
	if err := s.UnsDB().WithContext(ctx).Exec(string(sinkSQLBytes)).Error; err != nil {
		return err
	}
	if err := s.runVersionedMigrations(ctx, "common", s.commonDB, "migrations/versions/common"); err != nil {
		return err
	}
	if err := s.runVersionedMigrations(ctx, "sink", s.UnsDB(), "migrations/versions/sink"); err != nil {
		return err
	}
	return s.syncTimeseriesRetentionPolicy(ctx)
}

func pingDB(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func closeDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}
