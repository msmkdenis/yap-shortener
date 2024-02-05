package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/pkg/jwtgen"
)

// JWTAuth represents JWT authentication middleware.
type JWTAuth struct {
	jwtManager *jwtgen.JWTManager
	logger     *zap.Logger
}

// InitJWTAuth returns a new instance of JWTAuth.
func InitJWTAuth(jwtManager *jwtgen.JWTManager, logger *zap.Logger) *JWTAuth {
	j := &JWTAuth{
		jwtManager: jwtManager,
		logger:     logger,
	}
	return j
}

// JWTAuth checks token and sets userID in the context. otherwise returns 401.
func (j *JWTAuth) JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Request().Cookie(j.jwtManager.TokenName)
			if err != nil {
				j.logger.Info("authentification failed", zap.Error(err))
				return c.NoContent(http.StatusUnauthorized)
			}
			userID, err := j.jwtManager.GetUserID(cookie.Value)
			if err != nil {
				j.logger.Info("authentification failed", zap.Error(err))
				return c.NoContent(http.StatusUnauthorized)
			}
			c.Set("userID", userID)
			j.logger.Info("authenticated", zap.String("userID", userID))
			return next(c)
		}
	}
}
