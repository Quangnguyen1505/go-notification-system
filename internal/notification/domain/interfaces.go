package domain

import "context"

type (
	NotificationRepo interface {
		Create(context.Context, *NotificationModel) error
		GetByID(context.Context, *NotificationModel) error
		GetAllByUserID(context.Context, *NotificationModel) error
	}
)
