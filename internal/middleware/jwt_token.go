package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/pkg/auth"
)

// JWTCheckerCreator represents JWT checker creator middleware.
type JWTCheckerCreator struct {
	jwtManager *auth.JWTManager
	logger     *zap.Logger
}

// InitJWTCheckerCreator returns a new instance of JWTCheckerCreator.
func InitJWTCheckerCreator(jwtManager *auth.JWTManager, logger *zap.Logger) *JWTCheckerCreator {
	j := &JWTCheckerCreator{
		jwtManager: jwtManager,
		logger:     logger,
	}
	return j
}

// JWTCheckOrCreate checks token and sets userID in the context.
// Otherwise creates new token and sets it in the context.
func (j *JWTCheckerCreator) JWTCheckOrCreate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, cookieErr := c.Request().Cookie(j.jwtManager.TokenName)
			if cookieErr != nil {
				j.logger.Info("token not found, creating new token", zap.Error(cookieErr))
				token := j.setCookieAndReturn(c)
				newUserID, err := j.jwtManager.GetUserID(token)
				if err != nil {
					j.logger.Error("unable to parse UserID, while creating new token", zap.Error(err))
					return c.NoContent(http.StatusInternalServerError)
				}
				c.Set("userID", newUserID)
				j.logger.Info("token created", zap.String("userID", newUserID))
				return next(c)
			}

			userID, err := j.jwtManager.GetUserID(cookie.Value)
			if err != nil {
				j.logger.Warn("unable to parse UserID, creating new token", zap.Error(err))
				token := j.setCookieAndReturn(c)
				newUserID, err := j.jwtManager.GetUserID(token)
				if err != nil {
					j.logger.Error("unable to parse UserID, while creating new token", zap.Error(err))
					return c.NoContent(http.StatusInternalServerError)
				}
				c.Set("userID", newUserID)
				j.logger.Info("token created", zap.String("userID", newUserID))
				return next(c)
			}

			c.Set("userID", userID)
			return next(c)
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
