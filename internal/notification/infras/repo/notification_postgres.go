package repo

import (
	"context"

	"github.com/google/wire"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/domain"
	"github.com/quangnguyen1505/go-notification-system/pkg/logger"
	"github.com/quangnguyen1505/go-notification-system/pkg/postgres"
)

var _ domain.NotificationRepo = (*notificationRepo)(nil)

var RepositorySet = wire.NewSet(
	NewnotificationRepo,
	wire.Bind(new(domain.NotificationRepo), new(*notificationRepo)),
)

type notificationRepo struct {
	pg     postgres.DBEngine
	logger *logger.LoggerZap
}

func NewnotificationRepo(pg postgres.DBEngine, logger *logger.LoggerZap) *notificationRepo {
	return &notificationRepo{
		pg:     pg,
		logger: logger,
	}
}

func (r *notificationRepo) Create(ctx context.Context, _ *domain.NotificationModel) error {
	r.logger.Info("Creating notification in memory")
	// Implementation for creating a notification in memory
	return nil
}

func (r *notificationRepo) GetByID(ctx context.Context, _ *domain.NotificationModel) (string, error) {
	r.logger.Info("Getting notification by ID in memory")
	// Implementation for getting a notification by ID in memory
	return "hehe", nil
}

func (r *notificationRepo) GetAllByUserID(ctx context.Context, _ *domain.NotificationModel) error {
	r.logger.Info("Getting all notifications by user ID in memory")
	// Implementation for getting all notifications by user ID in memory
	return nil
}
