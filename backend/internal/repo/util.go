package repo

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"backend/internal/infra/idgen"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// returningAll asks postgres to RETURNING * so an UPDATE can populate the
// destination model in a single round trip, mirroring the old UPDATE ... RETURNING.
func returningAll() clause.Returning { return clause.Returning{} }

func normalizeDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicate
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrDuplicate
	}
	return err
}

func hashSecret(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func secretEdges(raw string) (string, string) {
	if len(raw) <= 16 {
		return raw, raw
	}
	return raw[:16], raw[len(raw)-4:]
}

func randomHex(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	}
	return hex.EncodeToString(buf)
}

func RandomHexForApp(n int) string {
	return randomHex(n)
}

func randomReadableSecret() string {
	return fmt.Sprintf("admin-%d", time.Now().UTC().UnixNano())
}

func ensureID(id *int64) {
	if id != nil && *id == 0 {
		*id = idgen.NextID()
	}
}

func splitCSV(value string) []string {
	var out []string
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func strconvInt64(value int64) string {
	if value <= 0 {
		return "0"
	}
	return fmt.Sprintf("%d", value)
}

func methodMatches(methods, method string) bool {
	if methods == "" || methods == "*" {
		return true
	}
	method = strings.ToUpper(method)
	for _, m := range strings.Split(methods, ",") {
		if strings.TrimSpace(strings.ToUpper(m)) == method {
			return true
		}
	}
	return false
}

func pathMatches(pattern, path string) bool {
	if pattern == path {
		return true
	}
	if strings.HasSuffix(pattern, "/**") {
		prefix := strings.TrimSuffix(pattern, "/**")
		return path == prefix || strings.HasPrefix(path, prefix+"/")
	}
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(patternParts) != len(pathParts) {
		return false
	}
	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], ":") && pathParts[i] != "" {
			continue
		}
		if patternParts[i] != pathParts[i] {
			return false
		}
	}
	return len(patternParts) > 0
}
