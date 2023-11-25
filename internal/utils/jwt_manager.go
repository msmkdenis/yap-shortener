package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/apperrors"
)

type JWTManager struct {
	logger    *zap.Logger
	TokenName string
}

const tokenExp = time.Hour * 24
const secretKey = "supersecretkey"
const tokenName = "token"

type claims struct {
	jwt.RegisteredClaims
	UserID string
}

func InitJWTManager(logger *zap.Logger) *JWTManager {
	j := &JWTManager{
		logger:    logger,
		TokenName: tokenName,
	}
	return j
}

func (j *JWTManager) BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: uuid.New().String(),
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTManager) GetUserID(tokenString string) (string, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperrors.NewValueError(fmt.Sprintf("unexpected signing method: %v", t.Header["alg"]), Caller(), errors.New("unexpected signing method"))
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return "", apperrors.NewValueError("token is not valid", Caller(), errors.New("token is not valid"))
	}

	return claims.UserID, nil
}
