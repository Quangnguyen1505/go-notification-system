package notification

import (
	"context"

	"github.com/google/wire"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/domain"
	"github.com/quangnguyen1505/go-notification-system/pkg/logger"
	"go.uber.org/zap"
)

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(
	NewService,
	wire.Bind(new(UseCase), new(*service)),
)

type service struct {
	repo   domain.NotificationRepo
	logger *logger.LoggerZap
}

func NewService(repo domain.NotificationRepo, logger *logger.LoggerZap) *service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Create(ctx context.Context, model *domain.NotificationModel) error {
	s.logger.Info("Creating notification", zap.Any("model", model))
	return s.repo.Create(ctx, model)
}

func (s *service) GetByID(ctx context.Context, model *domain.NotificationModel) (string, error) {
	s.logger.Info("Getting notification by ID", zap.Any("model", model))
	return s.repo.GetByID(ctx, model)
}

func (s *service) GetAllByUserID(ctx context.Context, model *domain.NotificationModel) error {
	s.logger.Info("Getting all notifications by user ID", zap.Any("model", model))
	return s.repo.GetAllByUserID(ctx, model)
}
