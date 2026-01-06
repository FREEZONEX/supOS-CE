package authsvc

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"gitee.com/unitedrhino/share/errors"
)

// DecodeJWTClaims parses the payload section of a JWT access token.
func DecodeJWTClaims(accessToken string) (map[string]any, error) {
	parts := strings.Split(accessToken, ".")
	if len(parts) < 2 {
		return nil, errors.Parameter.WithMsg("invalid access token format")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.Parameter.WithMsgf("decode token payload failed: %v", err)
	}
	claims := make(map[string]any)
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.Parameter.WithMsgf("unmarshal token payload failed: %v", err)
	}
	return claims, nil
}

// ClaimString extracts a string value from the JWT claims map.
func ClaimString(claims map[string]any, key string) string {
	if claims == nil {
		return ""
	}
	if v, ok := claims[key]; ok {
		switch val := v.(type) {
		case string:
			return strings.TrimSpace(val)
		case []any:
			if len(val) > 0 {
				if s, ok := val[0].(string); ok {
					return strings.TrimSpace(s)
				}
			}
		}
	}
	return ""
}
