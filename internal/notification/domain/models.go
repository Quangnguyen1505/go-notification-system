package domain

import "github.com/google/uuid"

type NotificationModel struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Message string
}
