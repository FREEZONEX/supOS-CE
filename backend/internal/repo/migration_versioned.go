package repo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const migrationLedgerTable = "sys_schema_migration"

type versionedMigration struct {
	Scope    string
	Version  string
	Sequence int
	Name     string
	Path     string
	SQL      string
	Checksum string
}

type migrationLedgerRow struct {
	Checksum string
	Success  bool
}

func (s *Store) runVersionedMigrations(ctx context.Context, scope string, db *gorm.DB, dir string) error {
	if db == nil {
		return fmt.Errorf("run %s migrations: database is nil", scope)
	}
	if err := ensureMigrationLedger(ctx, db); err != nil {
		return fmt.Errorf("ensure %s migration ledger: %w", scope, err)
	}
	migrations, err := loadVersionedMigrations(scope, dir)
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		if err := applyVersionedMigration(ctx, db, migration); err != nil {
			return err
		}
	}
	return nil
}

func ensureMigrationLedger(ctx context.Context, db *gorm.DB) error {
	const ledgerSQL = `
CREATE TABLE IF NOT EXISTS sys_schema_migration (
    id BIGSERIAL PRIMARY KEY,
    scope VARCHAR(32) NOT NULL,
    version VARCHAR(64) NOT NULL,
    sequence INTEGER NOT NULL,
    name VARCHAR(200) NOT NULL,
    checksum VARCHAR(64) NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    duration_ms BIGINT NOT NULL DEFAULT 0,
    success BOOLEAN NOT NULL DEFAULT false,
    error_message TEXT NOT NULL DEFAULT '',
    UNIQUE(scope, version, sequence, name)
);
CREATE INDEX IF NOT EXISTS idx_sys_schema_migration_scope_version
    ON sys_schema_migration(scope, version, sequence);
`
	return db.WithContext(ctx).Exec(ledgerSQL).Error
}

func loadVersionedMigrations(scope string, dir string) ([]versionedMigration, error) {
	entries, err := fs.Glob(migrationFS, dir+"/*.sql")
	if err != nil {
		return nil, fmt.Errorf("list %s migrations: %w", scope, err)
	}
	migrations := make([]versionedMigration, 0, len(entries))
	for _, entry := range entries {
		content, err := migrationFS.ReadFile(entry)
		if err != nil {
			return nil, fmt.Errorf("read migration %s: %w", entry, err)
		}
		migration, err := parseVersionedMigration(scope, entry, string(content))
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}
	sort.Slice(migrations, func(i, j int) bool {
		return compareVersionedMigration(migrations[i], migrations[j]) < 0
	})
	return migrations, nil
}

func parseVersionedMigration(scope string, filePath string, sql string) (versionedMigration, error) {
	base := path.Base(filePath)
	if !strings.HasPrefix(base, "V") || !strings.HasSuffix(base, ".sql") {
		return versionedMigration{}, fmt.Errorf("invalid migration file name %q: expected V<version>_<sequence>__<name>.sql", base)
	}
	nameWithoutExt := strings.TrimSuffix(strings.TrimPrefix(base, "V"), ".sql")
	versionPart, rest, ok := strings.Cut(nameWithoutExt, "_")
	if !ok {
		return versionedMigration{}, fmt.Errorf("invalid migration file name %q: missing sequence", base)
	}
	sequencePart, namePart, ok := strings.Cut(rest, "__")
	if !ok {
		return versionedMigration{}, fmt.Errorf("invalid migration file name %q: missing migration name", base)
	}
	sequence, err := strconv.Atoi(sequencePart)
	if err != nil || sequence <= 0 {
		return versionedMigration{}, fmt.Errorf("invalid migration sequence in %q", base)
	}
	namePart = strings.TrimSpace(namePart)
	if namePart == "" {
		return versionedMigration{}, fmt.Errorf("invalid migration file name %q: empty migration name", base)
	}
	sum := sha256.Sum256([]byte(sql))
	return versionedMigration{
		Scope:    scope,
		Version:  strings.TrimSpace(versionPart),
		Sequence: sequence,
		Name:     namePart,
		Path:     filePath,
		SQL:      sql,
		Checksum: hex.EncodeToString(sum[:]),
	}, nil
}

func compareVersionedMigration(a, b versionedMigration) int {
	if c := compareSemver(a.Version, b.Version); c != 0 {
		return c
	}
	if a.Sequence != b.Sequence {
		if a.Sequence < b.Sequence {
			return -1
		}
		return 1
	}
	return strings.Compare(a.Name, b.Name)
}

func compareSemver(a, b string) int {
	aCore, aSuffix := splitVersionSuffix(a)
	bCore, bSuffix := splitVersionSuffix(b)
	for i := 0; i < 3; i++ {
		av := versionComponent(aCore, i)
		bv := versionComponent(bCore, i)
		if av < bv {
			return -1
		}
		if av > bv {
			return 1
		}
	}
	if aSuffix == bSuffix {
		return 0
	}
	if strings.HasPrefix(aSuffix, "-") && bSuffix == "" {
		return -1
	}
	if aSuffix == "" && strings.HasPrefix(bSuffix, "-") {
		return 1
	}
	return strings.Compare(aSuffix, bSuffix)
}

func splitVersionSuffix(version string) (string, string) {
	version = strings.TrimSpace(version)
	for _, sep := range []string{"-", "+"} {
		if idx := strings.Index(version, sep); idx >= 0 {
			return version[:idx], version[idx:]
		}
	}
	return version, ""
}

func versionComponent(version string, index int) int {
	parts := strings.Split(version, ".")
	if index >= len(parts) {
		return 0
	}
	value, _ := strconv.Atoi(parts[index])
	return value
}

func applyVersionedMigration(ctx context.Context, db *gorm.DB, migration versionedMigration) error {
	var existing migrationLedgerRow
	err := db.WithContext(ctx).
		Table(migrationLedgerTable).
		Select("checksum", "success").
		Where("scope = ? AND version = ? AND sequence = ? AND name = ?", migration.Scope, migration.Version, migration.Sequence, migration.Name).
		Take(&existing).Error
	if err == nil {
		if existing.Checksum != migration.Checksum {
			return fmt.Errorf("migration checksum mismatch for %s", migration.Path)
		}
		if existing.Success {
			return nil
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("read migration ledger for %s: %w", migration.Path, err)
	}

	started := time.Now()
	runErr := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(migration.SQL).Error; err != nil {
			return err
		}
		return recordMigrationResult(ctx, tx, migration, time.Since(started).Milliseconds(), true, "")
	})
	if runErr != nil {
		recordErr := recordMigrationResult(ctx, db, migration, time.Since(started).Milliseconds(), false, runErr.Error())
		if recordErr != nil {
			return fmt.Errorf("run migration %s: %w; record failure: %v", migration.Path, runErr, recordErr)
		}
		return fmt.Errorf("run migration %s: %w", migration.Path, runErr)
	}
	return nil
}

func recordMigrationResult(ctx context.Context, db *gorm.DB, migration versionedMigration, durationMS int64, success bool, message string) error {
	const sql = `
INSERT INTO sys_schema_migration(scope, version, sequence, name, checksum, duration_ms, success, error_message, applied_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(scope, version, sequence, name) DO UPDATE SET
    checksum = EXCLUDED.checksum,
    duration_ms = EXCLUDED.duration_ms,
    success = EXCLUDED.success,
    error_message = EXCLUDED.error_message,
    applied_at = EXCLUDED.applied_at
`
	return db.WithContext(ctx).Exec(sql,
		migration.Scope,
		migration.Version,
		migration.Sequence,
		migration.Name,
		migration.Checksum,
		durationMS,
		success,
		message,
	).Error
}
