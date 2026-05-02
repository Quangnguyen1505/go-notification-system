package repo

import (
	"context"

	"github.com/google/wire"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/domain"
	"github.com/quangnguyen1505/go-notification-system/pkg/logger"
)

var _ domain.NotificationRepo = (*notificationInMemRepo)(nil)

var RepositorySet = wire.NewSet(
	NewNotificationInMemRepo,
	wire.Bind(new(domain.NotificationRepo), new(*notificationInMemRepo)),
)

type notificationInMemRepo struct {
	logger *logger.LoggerZap
}

func NewNotificationInMemRepo(logger *logger.LoggerZap) *notificationInMemRepo {
	return &notificationInMemRepo{
		logger: logger,
	}
}

func (r *notificationInMemRepo) Create(ctx context.Context, _ *domain.NotificationModel) error {
	r.logger.Info("Creating notification in memory")
	// Implementation for creating a notification in memory
	return nil
}

func (r *notificationInMemRepo) GetByID(ctx context.Context, _ *domain.NotificationModel) (string, error) {
	r.logger.Info("Getting notification by ID in memory")
	// Implementation for getting a notification by ID in memory
	return "hehe", nil
}

func (r *notificationInMemRepo) GetAllByUserID(ctx context.Context, _ *domain.NotificationModel) error {
	r.logger.Info("Getting all notifications by user ID in memory")
	// Implementation for getting all notifications by user ID in memory
	return nil
}
