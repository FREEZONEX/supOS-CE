package repo

import "strings"

const (
	APIKeyPermissionRead  = "read_only"
	APIKeyPermissionWrite = "data_writer"
	APIKeyPermissionFull  = "full_access"
)

func NormalizeAPIKeyPermission(permission string) string {
	switch strings.TrimSpace(strings.ToLower(permission)) {
	case "", "read", "read_only", "readonly":
		return APIKeyPermissionRead
	case "write", "data_writer", "writer":
		return APIKeyPermissionWrite
	case "full", "full_access", "fullaccess", "manage":
		return APIKeyPermissionFull
	default:
		return APIKeyPermissionRead
	}
}

func APIKeyResourceKeysForPermission(permission string) []string {
	read := []string{"uns.read", "flow.read", "openapi.base"}
	write := append(append([]string{}, read...), "uns.write", "flow.manage")
	full := append(append([]string{}, write...), "uns.manage", "uns.delete")
	switch NormalizeAPIKeyPermission(permission) {
	case APIKeyPermissionWrite:
		return uniqueStrings(write)
	case APIKeyPermissionFull:
		return uniqueStrings(full)
	default:
		return uniqueStrings(read)
	}
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
