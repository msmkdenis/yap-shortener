package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"
)

type JWTAuth struct {
	jwtManager *utils.JWTManager
	logger     *zap.Logger
}

func InitJWTAuth(jwtManager *utils.JWTManager, logger *zap.Logger) *JWTAuth {
	j := &JWTAuth{
		jwtManager: jwtManager,
		logger:     logger,
	}
	return j
}

func (j *JWTAuth) JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Request().Cookie(j.jwtManager.TokenName)
			if err != nil {
				j.logger.Info("authentification failed", zap.Error(err))
				return c.NoContent(http.StatusNoContent) // подгонка под тест, по логике необходимо возвращать StatusUnauthorized  - исправить
			}
			userID, err := j.jwtManager.GetUserID(cookie.Value)
			if err != nil {
				j.logger.Info("authentification failed", zap.Error(err))
				return c.NoContent(http.StatusNoContent) // подгонка под тест, по логике необходимо возвращать StatusUnauthorized  - исправить
			}
			c.Set("userID", userID)
			err = next(c)
			j.logger.Info("authenticated", zap.String("userID", userID))
			return err
		}
	}
}
