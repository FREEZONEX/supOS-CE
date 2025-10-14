package users

import (
	"gitee.com/unitedrhino/share/errors"
	"github.com/golang-jwt/jwt/v5"
)

// 创建一个token
func CreateToken(secretKey string, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseTokenWithFunc(claim jwt.Claims, tokenString string, f jwt.Keyfunc) error {
	token, err := jwt.ParseWithClaims(tokenString, claim, f)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return errors.TokenExpired.WithMsg("登录过期,请退出重新登录")
		case errors.Is(err, jwt.ErrTokenMalformed):
			return errors.TokenMalformed.WithMsg("登录失效,请退出重新登录")
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return errors.TokenNotValidYet.WithMsg("登录失效,请退出重新登录")
		default:
			return errors.TokenInvalid.WithMsg("登录失效,请退出重新登录")
		}
	}
	if token != nil {
		if token.Valid {
			return nil
		}
		return errors.TokenInvalid.WithMsg("登录失效,请退出重新登录")

	} else {
		return errors.TokenInvalid.WithMsg("登录失效,请退出重新登录")
	}
}

// 解析 token
func ParseToken(claim jwt.Claims, tokenString string, secretKey string) error {
	return ParseTokenWithFunc(claim, tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
}
