package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"
)

type JWTCheckerCreator struct {
	jwtManager *utils.JWTManager
	logger     *zap.Logger
}

func InitJWTCheckerCreator(jwtManager *utils.JWTManager, logger *zap.Logger) *JWTCheckerCreator {
	j := &JWTCheckerCreator{
		jwtManager: jwtManager,
		logger:     logger,
	}
	return j
}

func (j *JWTCheckerCreator) JWTCheckOrCreate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, er := c.Request().Cookie(j.jwtManager.TokenName)
			if er != nil {
				j.logger.Warn("token not found, creating new token", zap.Error(er))
				token := j.setCookieAndReturn(c)
				userID, _ := j.jwtManager.GetUserID(token)
				c.Set("userID", userID)
				err := next(c)
				j.logger.Info("token created", zap.String("userID", userID))
				return err
			}

			userID, err := j.jwtManager.GetUserID(cookie.Value)
			if err != nil {
				j.logger.Warn("unable to parse UserID, creating new token", zap.Error(err))
				token := j.setCookieAndReturn(c)
				userID, _ := j.jwtManager.GetUserID(token)
				c.Set("userID", userID)
				err := next(c)
				j.logger.Info("token created", zap.String("userID", userID))
				return err
			}

			j.logger.Info("token checked", zap.String("userID", userID))
			c.Set("userID", userID)
			err = next(c)
			j.logger.Info("token created", zap.String("userID", userID))
			return err
		}
	}
}

func (j *JWTCheckerCreator) setCookieAndReturn(c echo.Context) string {
	tokenString, err := j.jwtManager.BuildJWTString()
	if err != nil {
		j.logger.Fatal("unable to create token", zap.Error(err))
	}

	cookie := &http.Cookie{
		Name:  j.jwtManager.TokenName,
		Value: tokenString,
	}
	c.SetCookie(cookie)
	return cookie.Value
}
