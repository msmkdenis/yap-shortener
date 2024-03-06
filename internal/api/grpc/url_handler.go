package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-shortener/internal/dto"
	"github.com/msmkdenis/yap-shortener/internal/model"
	pb "github.com/msmkdenis/yap-shortener/internal/proto"
	urlErr "github.com/msmkdenis/yap-shortener/internal/urlerr"
	"github.com/msmkdenis/yap-shortener/pkg/apperr"
	"github.com/msmkdenis/yap-shortener/pkg/jwtgen"
	"github.com/msmkdenis/yap-shortener/pkg/workerpool"
)

type URLHandler struct {
	urlService    URLService
	urlPrefix     string
	trustedSubnet string
	jwtManager    *jwtgen.JWTManager
	logger        *zap.Logger
	wg            *sync.WaitGroup
	pb.UnimplementedURLShortenerServer
}

// URLService represents URL service interface.
type URLService interface {
	Add(ctx context.Context, s string, host string, userID string) (*model.URL, error)
	AddAll(ctx context.Context, urls []dto.URLBatchRequest, host string, userID string) ([]dto.URLBatchResponse, error)
	GetAll(ctx context.Context) ([]string, error)
	GetAllByUserID(ctx context.Context, userID string) ([]dto.URLBatchResponseByUserID, error)
	DeleteAll(ctx context.Context) error
	DeleteURLByUserID(ctx context.Context, userID string, shortURLs string) error
	GetByyID(ctx context.Context, key string) (string, error)
	GetStats(ctx context.Context) (*dto.URLStats, error)
	Ping(ctx context.Context) error
}

// NewURLHandler creates a new gRPC URLHandler instance
func NewURLHandler(service URLService, urlPrefix string, trustedSubnet string, jwtManager *jwtgen.JWTManager, logger *zap.Logger, wg *sync.WaitGroup) *URLHandler {
	handler := &URLHandler{
		urlService:    service,
		urlPrefix:     urlPrefix,
		trustedSubnet: trustedSubnet,
		jwtManager:    jwtManager,
		logger:        logger,
		wg:            wg,
	}

	return handler
}

// GetListURLs handles gRPC GetListURLs request
func (h *URLHandler) GetListURLs(ctx context.Context, _ *pb.GetListURLsRequest) (*pb.GetListURLsResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	urls, err := h.urlService.GetAll(ctx)
	if err != nil {
		h.logger.Info("GetUrls", zap.Error(err))
		return nil, err
	}

	response := &pb.GetListURLsResponse{
		Urls: urls,
	}

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		h.logger.Error("GRPCInternalServerError: internal error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return response, nil
}

// PostURL handles gRPC PostURL request
func (h *URLHandler) PostURL(ctx context.Context, in *pb.PostURLRequest) (*pb.PostURLResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	if in.Url == "" {
		h.logger.Info("GRPCBadRequest", zap.Error(status.Error(codes.InvalidArgument, "empty url")))
		return nil, status.Error(codes.InvalidArgument, "empty url")
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		h.logger.Error("Internal server error", zap.Error(urlErr.ErrUnableToGetUserIDFromContext))
		return nil, status.Error(codes.Internal, "internal error")
	}

	url, err := h.urlService.Add(ctx, in.Url, h.urlPrefix, userID)
	if err != nil && !errors.Is(err, urlErr.ErrURLAlreadyExists) {
		h.logger.Error("GRPCInternalServerError: internal error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Internal, "internal error")
	}

	if errors.Is(err, urlErr.ErrURLAlreadyExists) {
		h.logger.Warn("GRPCConflict: url already exists", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return &pb.PostURLResponse{ShortUrl: url.Shortened}, status.Error(codes.AlreadyExists, "url already exists")
	}

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		h.logger.Error("GRPCInternalServerError: internal error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostURLResponse{ShortUrl: url.Shortened}, nil
}

// PostBatchURLs handles gRPC PostBatchURLs request
func (h *URLHandler) PostBatchURLs(ctx context.Context, in *pb.PostBatchURLRequest) (*pb.PostBatchURLResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	if len(in.BatchUrls) == 0 {
		h.logger.Info("GRPCBadRequest", zap.Error(status.Error(codes.InvalidArgument, "empty batch urls")))
		return nil, status.Error(codes.InvalidArgument, "empty batch urls")
	}

	urls := make([]dto.URLBatchRequest, 0, len(in.BatchUrls))
	for _, v := range in.BatchUrls {
		urls = append(urls, dto.URLBatchRequest{
			CorrelationID: v.CorrelationId,
			OriginalURL:   v.OriginalUrl,
		})
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		h.logger.Error("Internal server error", zap.Error(urlErr.ErrUnableToGetUserIDFromContext))
		return nil, status.Error(codes.Internal, "internal error")
	}

	savedURLs, err := h.urlService.AddAll(ctx, urls, h.urlPrefix, userID)
	if err != nil {
		h.logger.Error("GRPCInternalServerError: internal error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Internal, "internal error")
	}

	batchURLs := make([]*pb.BatchURLResponse, 0, len(savedURLs))
	for _, v := range savedURLs {
		batchURLs = append(batchURLs, &pb.BatchURLResponse{
			CorrelationId: v.CorrelationID,
			ShortenedURL:  v.ShortenedURL,
		})
	}

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		h.logger.Error("GRPCInternalServerError: internal error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostBatchURLResponse{BatchUrls: batchURLs}, nil
}

// GetURL handles gRPC GetURL request
func (h *URLHandler) GetURL(ctx context.Context, in *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	if in.ShortUrl == "" {
		h.logger.Info("GRPCBadRequest", zap.Error(status.Error(codes.InvalidArgument, "empty url")))
		return nil, status.Error(codes.InvalidArgument, "empty url")
	}

	originalURL, err := h.urlService.GetByyID(ctx, in.ShortUrl)

	switch {
	case errors.Is(err, urlErr.ErrURLNotFound):
		h.logger.Info("StatusBadRequest: url not found", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.NotFound, fmt.Sprintf("URL with id %s not found", in.ShortUrl))

	case errors.Is(err, urlErr.ErrURLDeleted):
		h.logger.Info("StatusBadRequest: url not found", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Unavailable, fmt.Sprintf("URL with id %s has been deleted", in.ShortUrl))

	case err != nil:
		h.logger.Error("GRPCInternalServerError: internal error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Internal, "internal error")
	}

	md.Append("Location", originalURL)
	err = grpc.SendHeader(ctx, md)
	if err != nil {
		h.logger.Error("GRPCInternalServerError: internal error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.GetURLResponse{Url: originalURL}, nil
}

// Ping handles gRPC Ping request
func (h *URLHandler) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	err := h.urlService.Ping(ctx)
	if err != nil {
		h.logger.Error("GRPCInternalServerError: internal error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Internal, "internal database error")
	}

	return &pb.PingResponse{}, nil
}

// DeleteAllURLs handles gRPC DeleteAllURLs request
func (h *URLHandler) DeleteAllURLs(ctx context.Context, _ *pb.DeleteAllURLsRequest) (*pb.DeleteAllURLsResponse, error) {
	if err := h.urlService.DeleteAll(ctx); err != nil {
		h.logger.Error("StatusInternalServerError: Unknown error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.DeleteAllURLsResponse{}, nil
}

// GetURLsByUserID handles gRPC GetURLsByUserID request
func (h *URLHandler) GetURLsByUserID(ctx context.Context, _ *pb.GetURLsByUserIDRequest) (*pb.GetURLsByUserIDResponse, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		h.logger.Error("Internal server error", zap.Error(urlErr.ErrUnableToGetUserIDFromContext))
		return nil, status.Error(codes.Internal, "internal error")
	}

	savedURLs, err := h.urlService.GetAllByUserID(ctx, userID)
	if err != nil && !errors.Is(err, urlErr.ErrURLNotFound) {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Internal, "internal error")
	}

	if errors.Is(err, urlErr.ErrURLNotFound) {
		h.logger.Warn("StatusNoContent: urls not found", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.NotFound, "urls not found")
	}

	urls := make([]*pb.URLByUserID, 0, len(savedURLs))
	for _, url := range savedURLs {
		urls = append(urls, &pb.URLByUserID{
			ShortUrl:    url.ShortURL,
			OriginalUrl: url.OriginalURL,
			DeletedFlag: url.DeletedFlag,
		})
	}
	return &pb.GetURLsByUserIDResponse{Urls: urls}, nil
}

// DeleteURLsByUserID handles gRPC DeleteURLsByUserID request
func (h *URLHandler) DeleteURLsByUserID(ctx context.Context, in *pb.DeleteURLsByUserIDRequest) (*pb.DeleteURLsByUserIDResponse, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		h.logger.Error("Internal server error", zap.Error(urlErr.ErrUnableToGetUserIDFromContext))
		return nil, status.Error(codes.Internal, "internal error")
	}

	if len(in.ShortUrls) == 0 {
		h.logger.Info("GRPCBadRequest", zap.Error(status.Error(codes.InvalidArgument, "empty batch urls")))
		return nil, status.Error(codes.InvalidArgument, "empty batch urls")
	}

	workerPool := workerpool.NewWorkerPool(100, h.logger)
	workerPool.Start()
	defer workerPool.Stop()

	h.wg.Add(len(in.ShortUrls))
	for _, shortURL := range in.ShortUrls {
		log.Info("Submitting task", zap.String("delete shortURL", shortURL))
		url := shortURL
		userID := userID
		workerPool.Submit(func() error {
			defer h.wg.Done()
			return h.urlService.DeleteURLByUserID(context.WithoutCancel(ctx), userID, url)
		})
	}

	return &pb.DeleteURLsByUserIDResponse{}, nil
}

// GetStats handles gRPC GetStats request
func (h *URLHandler) GetStats(ctx context.Context, _ *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	if h.trustedSubnet == "" {
		return nil, status.Error(codes.PermissionDenied, "stats not available without trusted subnet")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	ip := md.Get("X-Real-IP")
	if len(ip) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing X-Real-IP")
	}

	_, ipNet, err := net.ParseCIDR(h.trustedSubnet)
	if err != nil {
		h.logger.Error("Internal server error", zap.Error(urlErr.ErrUnableToGetUserIDFromContext))
		return nil, status.Error(codes.Internal, "internal error")
	}

	if !ipNet.Contains(net.ParseIP(ip[0])) {
		return nil, status.Error(codes.PermissionDenied, "internal error")
	}

	stats, err := h.urlService.GetStats(ctx)
	if err != nil {
		h.logger.Error("StatusInternalServerError: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.GetStatsResponse{
		Urls:  uint32(stats.Urls),
		Users: uint32(stats.Users),
	}, nil
}
