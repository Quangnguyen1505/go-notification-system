//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/quangnguyen1505/go-notification-system/cmd/notification/config"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/app/router"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/infras/repo"
	NotificationUC "github.com/quangnguyen1505/go-notification-system/internal/notification/usecases/notification"
	"github.com/quangnguyen1505/go-notification-system/pkg/logger"
	"github.com/quangnguyen1505/go-notification-system/pkg/postgres"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	logger *logger.LoggerZap,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		router.NotificationGRPCServerSet,
		NotificationUC.UseCaseSet,
		repo.RepositorySet,
	))
}

func dbEngineFunc(config *config.Config, logger *logger.LoggerZap) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(&config.Postgres, logger)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}
