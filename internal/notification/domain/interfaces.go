package domain

import "context"

type (
	NotificationRepo interface {
		Create(context.Context, *NotificationModel) error
		GetByID(context.Context, *NotificationModel) (string, error)
		GetAllByUserID(context.Context, *NotificationModel) error
	}
)
