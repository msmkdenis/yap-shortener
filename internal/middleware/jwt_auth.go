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

var authMandatoryMethods = map[string]struct{}{
	"/proto.URLShortener/GetURLsByUserID":    {},
	"/proto.URLShortener/DeleteURLsByUserID": {},
}

type UserIDContextKey string

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

// GRPCJWTAuth checks token from gRPC metadata and sets userID in the context. otherwise returns 401.
func (j *JWTAuth) GRPCJWTAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := authMandatoryMethods[info.FullMethod]; !ok {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	c := md.Get(j.jwtManager.TokenName)
	if len(c) < 1 {
		j.logger.Info("authentification failed")
		return nil, status.Errorf(codes.Unauthenticated, "no token found")
	}

	userID, err := j.jwtManager.GetUserID(c[0])
	if err != nil {
		j.logger.Info("authentification failed", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "authentification by UserID failed")
	}

	ctx = context.WithValue(ctx, UserIDContextKey("userID"), userID)
	return handler(ctx, req)
}
