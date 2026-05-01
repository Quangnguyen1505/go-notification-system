package repo

import (
	"context"

	"github.com/google/wire"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/domain"
)

var _ domain.NotificationRepo = (*notificationInMemRepo)(nil)

var RepositorySet = wire.NewSet(
	NewNotificationInMemRepo,
	wire.Bind(new(domain.NotificationRepo), new(*notificationInMemRepo)),
)

type notificationInMemRepo struct {
	helloworld string
}

func NewNotificationInMemRepo() *notificationInMemRepo {
	return &notificationInMemRepo{
		helloworld: "Hello World",
	}
}

func (r *notificationInMemRepo) Create(ctx context.Context, _ *domain.NotificationModel) error {
	// Implementation for creating a notification in memory
	return nil
}

func (r *notificationInMemRepo) GetByID(ctx context.Context, _ *domain.NotificationModel) error {
	// Implementation for getting a notification by ID in memory
	return nil
}

func (r *notificationInMemRepo) GetAllByUserID(ctx context.Context, _ *domain.NotificationModel) error {
	// Implementation for getting all notifications by user ID in memory
	return nil
}
