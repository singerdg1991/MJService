package jwtlib

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/hoitek/Maja-Service/config"
	"time"
)

type JwtPayload struct {
	ID       int
	Username string
	Data     interface{}
}

func Encrypt(payload JwtPayload) (string, error) {
	var (
		AUDIENCE         = config.AppConfig.HostAddress
		TOKEN_EXPIRATION = config.AppConfig.TokenExpiration
	)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": payload.Username,
		"sub": payload.Data,
		"aud": AUDIENCE,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Duration(TOKEN_EXPIRATION) * time.Second)),
		"nbf": jwt.NewNumericDate(time.Now()),
		"iat": jwt.NewNumericDate(time.Now()),
		"jti": payload.ID,
	})
	return jwtToken.SignedString([]byte(config.AppConfig.SigningKey))
}
