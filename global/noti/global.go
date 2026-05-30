package noti

import (
	"github.com/quangnguyen1505/go-notification-system/cmd/notification/config"
	"github.com/quangnguyen1505/go-notification-system/pkg/logger"
)

var (
	Logger *logger.LoggerZap
	Config *config.Config
)
