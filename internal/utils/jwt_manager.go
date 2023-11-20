package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"go.uber.org/zap"
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
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},

		// собственное утверждение
		UserID: uuid.New().String(),
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
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
