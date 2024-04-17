package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-shortener/pkg/jwtgen"
)

var jwtCheckerSkipMethods = map[string]struct{}{
	"/proto.URLShortener/GetStats": {},
}

type TokenContextKey string

// JWTCheckerCreator represents JWT checker creator middleware.
type JWTCheckerCreator struct {
	jwtManager *jwtgen.JWTManager
	logger     *zap.Logger
}

// InitJWTCheckerCreator returns a new instance of JWTCheckerCreator.
func InitJWTCheckerCreator(jwtManager *jwtgen.JWTManager, logger *zap.Logger) *JWTCheckerCreator {
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

// JWTCheckOrCreate checks token from gRPC metadata and sets userID in the context.
// Otherwise creates new token and sets it in the context.
func (j *JWTCheckerCreator) GRPCJWTCheckOrCreate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := authMandatoryMethods[info.FullMethod]; ok {
		return handler(ctx, req)
	}

	if _, ok := jwtCheckerSkipMethods[info.FullMethod]; ok {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	c := md.Get(j.jwtManager.TokenName)
	if len(c) < 1 {
		j.logger.Info("token not found, creating new token")
		ctxWithID, err := j.setUserIDAndReturn(ctx, md)
		if err != nil {
			return nil, err
		}
		return handler(ctxWithID, req)
	}

	userID, err := j.jwtManager.GetUserID(c[0])
	if err != nil {
		j.logger.Warn("unable to parse UserID, creating new token", zap.Error(err))
		ctx, err = j.setUserIDAndReturn(ctx, md)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}

	ctx = context.WithValue(ctx, UserIDContextKey("userID"), userID)
	ctx = metadata.NewIncomingContext(ctx, md)
	return handler(ctx, req)
}

func (j *JWTCheckerCreator) setUserIDAndReturn(ctx context.Context, md metadata.MD) (context.Context, error) {
	cookie := j.makeCookie()
	ctx = context.WithValue(ctx, TokenContextKey(j.jwtManager.TokenName), cookie.Value)
	newUserID, err := j.jwtManager.GetUserID(cookie.Value)
	if err != nil {
		j.logger.Error("unable to parse UserID, while creating new token", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "unable to parse UserID, while creating new token")
	}
	ctx = context.WithValue(ctx, UserIDContextKey("userID"), newUserID)
	md.Append(j.jwtManager.TokenName, cookie.Value)
	ctx = metadata.NewIncomingContext(ctx, md)
	return ctx, nil
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

func (j *JWTCheckerCreator) makeCookie() *http.Cookie {
	tokenString, err := j.jwtManager.BuildJWTString()
	if err != nil {
		j.logger.Fatal("unable to create token", zap.Error(err))
	}

	cookie := &http.Cookie{
		Name:  j.jwtManager.TokenName,
		Value: tokenString,
	}
	return cookie
}
