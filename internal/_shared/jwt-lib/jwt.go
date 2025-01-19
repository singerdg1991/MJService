package jwtlib

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hoitek/Maja-Service/config"
	"math/big"
	"time"
)

type JwtPayload struct {
	ID       int64
	Username string
	Data     interface{}
}

func Encrypt(payload JwtPayload, expSeconds int64) (string, error) {
	var (
		AUDIENCE = config.AppConfig.HostAddress
	)
	claims := jwt.MapClaims{
		"iss": payload.Username,
		"sub": payload.Data,
		"aud": AUDIENCE,
		"nbf": jwt.NewNumericDate(time.Now()),
		"iat": jwt.NewNumericDate(time.Now()),
		"jti": payload.ID,
	}
	if expSeconds > 0 {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Duration(expSeconds) * time.Second))
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(config.AppConfig.JwtSigningKey))
}

func Decrypt(tokenString string) (JwtPayload, error) {
	var (
		AUDIENCE = config.AppConfig.HostAddress
	)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(config.AppConfig.JwtSigningKey), nil
	})
	if err != nil {
		return JwtPayload{}, err
	}
	if !token.Valid {
		return JwtPayload{}, err
	}
	if token.Claims.(jwt.MapClaims)["aud"] != AUDIENCE {
		return JwtPayload{}, err
	}
	return JwtPayload{
		ID:       int64(token.Claims.(jwt.MapClaims)["jti"].(float64)),
		Username: token.Claims.(jwt.MapClaims)["iss"].(string),
		Data:     token.Claims.(jwt.MapClaims)["sub"],
	}, nil
}

func SignToken(token []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.Sum256(token)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return []byte{}, err
	}
	return signature, nil
}

func VerifyToken(token []byte, signature []byte, publicKey *rsa.PublicKey) error {
	hash := sha256.Sum256(token)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
}

func GeneratePrivateAndPublicKey() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

type PublicKey struct {
	N *big.Int
	E int
}

func (pk PublicKey) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"N":"%s","E":%d}`, pk.N.String(), pk.E)), nil
}

func GenerateTokenPair(payload JwtPayload) (string, string, error) {
	var (
		accessTokenExp  = config.AppConfig.JwtTokenExpiration
		refreshTokenExp = config.AppConfig.JwtRefreshTokenExpiration
	)
	accessToken, err := Encrypt(payload, accessTokenExp)
	if err != nil {
		return "", "", err
	}
	payload.Data.(map[string]interface{})["accessToken"] = accessToken
	refreshToken, err := Encrypt(payload, refreshTokenExp)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
