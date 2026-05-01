package notification

import (
	"context"

	"github.com/quangnguyen1505/go-notification-system/internal/notification/domain"
)

type (
	UseCase interface {
		Create(context.Context, *domain.NotificationModel) error
		GetByID(context.Context, *domain.NotificationModel) error
		GetAllByUserID(context.Context, *domain.NotificationModel) error
	}
)
