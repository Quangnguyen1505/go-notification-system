package app

import (
	"github.com/quangnguyen1505/go-notification-system/pkg/postgres"
	"github.com/quangnguyen1505/go-notification-system/proto/gen"
)

type App struct {
	DB                     postgres.DBEngine
	NotificationGRPCServer gen.NotificationServiceServer
}

func New(
	db postgres.DBEngine,
	notificationGRPCServer gen.NotificationServiceServer,
) *App {
	return &App{
		DB:                     db,
		NotificationGRPCServer: notificationGRPCServer,
	}
}
