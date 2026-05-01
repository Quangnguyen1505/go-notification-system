package app

import (
	"github.com/quangnguyen1505/go-notification-system/cmd/notification/config"
	"github.com/quangnguyen1505/go-notification-system/proto/gen"
)

type App struct {
	Cfg                    *config.Config
	NotificationGRPCServer gen.NotificationServiceServer
}

func New(
	cfg *config.Config,
	notificationGRPCServer gen.NotificationServiceServer,
) *App {
	return &App{
		Cfg:                    cfg,
		NotificationGRPCServer: notificationGRPCServer,
	}
}
