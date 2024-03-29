// Package utils provides some utilities.
package jwtgen

import (
	"errors"
	"fmt"
	"time"

	"github.com/msmkdenis/yap-shortener/pkg/apperr"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// JWTManager represents the JWT manager.
type JWTManager struct {
	logger    *zap.Logger
	TokenName string
	secretKey string
}

const (
	tokenExp = time.Hour * 24
)

type claims struct {
	jwt.RegisteredClaims
	UserID string
}

// InitJWTManager returns a new instance of JWTManager.
func InitJWTManager(tokenName string, secretKey string, logger *zap.Logger) *JWTManager {
	j := &JWTManager{
		logger:    logger,
		TokenName: tokenName,
		secretKey: secretKey,
	}
	return j
}

// BuildJWTString creates JWT token with userID.
func (j *JWTManager) BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: uuid.New().String(),
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserID returns userID from JWT token.
func (j *JWTManager) GetUserID(tokenString string) (string, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperr.NewValueError(fmt.Sprintf("unexpected signing method: %v", t.Header["alg"]), apperr.Caller(), errors.New("unexpected signing method"))
			}
			return []byte(j.secretKey), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		j.logger.Warn("token is not valid", zap.Error(err))
		return "", apperr.NewValueError("token is not valid", apperr.Caller(), errors.New("token is not valid"))
	}

	return claims.UserID, nil
}
