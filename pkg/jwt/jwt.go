// Package jwt 提供 JWT 签发与解析。
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义 JWT 声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Sign 签发 JWT，token 有效期为 ttl。
func Sign(secret string, userID uint, userName, role string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "fing",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Parse 解析并校验 JWT。
func Parse(secret, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
