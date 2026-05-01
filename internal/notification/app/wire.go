//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/quangnguyen1505/go-notification-system/cmd/notification/config"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/app/router"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/infras/repo"
	NotificationUC "github.com/quangnguyen1505/go-notification-system/internal/notification/usecases/notification"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	grpcServer *grpc.Server,
) (*App, error) {
	panic(wire.Build(
		New,
		router.NotificationGRPCServerSet,
		repo.RepositorySet,
		NotificationUC.UseCaseSet,
	))
}
