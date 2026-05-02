package app

import (
	"github.com/quangnguyen1505/go-notification-system/cmd/notification/config"
	"github.com/quangnguyen1505/go-notification-system/pkg/logger"
	"github.com/quangnguyen1505/go-notification-system/pkg/postgres"
	"github.com/quangnguyen1505/go-notification-system/proto/gen"
)

type App struct {
	Cfg                    *config.Config
	Logger                 *logger.LoggerZap
	DB                     postgres.DBEngine
	NotificationGRPCServer gen.NotificationServiceServer
}

func New(
	cfg *config.Config,
	logger *logger.LoggerZap,
	db postgres.DBEngine,
	notificationGRPCServer gen.NotificationServiceServer,
) *App {
	return &App{
		Cfg:                    cfg,
		Logger:                 logger,
		DB:                     db,
		NotificationGRPCServer: notificationGRPCServer,
	}
}
