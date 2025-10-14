package users

import (
	"github.com/golang-jwt/jwt/v5"
)

// Custom claims structure
type OpenClaims struct {
	UserID     int64  `json:"userID,string"` //账号
	TenantCode string `json:"tenantCode"`
	Code       string `json:"code"`
	jwt.RegisteredClaims
}
