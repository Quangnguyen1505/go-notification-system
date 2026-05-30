//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/quangnguyen1505/go-notification-system/global/noti"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/app/router"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/infras/repo"
	NotificationUC "github.com/quangnguyen1505/go-notification-system/internal/notification/usecases/notification"
	"github.com/quangnguyen1505/go-notification-system/pkg/postgres"
	"google.golang.org/grpc"
)

func InitApp(
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

func dbEngineFunc() (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(&noti.Config.Postgres, noti.Logger)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}
