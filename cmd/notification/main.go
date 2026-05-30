package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/quangnguyen1505/go-notification-system/cmd/notification/config"
	"github.com/quangnguyen1505/go-notification-system/global/noti"
	notificationapp "github.com/quangnguyen1505/go-notification-system/internal/notification/app"
	"github.com/quangnguyen1505/go-notification-system/pkg/logger"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// set GOMAXPROCS -> get default from cpu core -> ex: 4 cores -> GOMAXPROCS=8
	_, err := maxprocs.Set()
	if err != nil {
		log.Printf("failed set max procs: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	noti.Config, err = config.NewConfig()
	if err != nil {
		log.Printf("failed get config: %v", err)
		return
	}

	noti.Logger = logger.NewLogger(noti.Config.Log)
	noti.Logger.Info("⚡ init app", zap.String("name", noti.Config.App.Name), zap.String("version", noti.Config.App.Version))
	// set up logrus
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	// logrus.SetOutput(os.Stdout)
	// logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	server := grpc.NewServer()

	if _, cleanup, err := notificationapp.InitApp(server); err != nil {
		noti.Logger.Error("failed init app", zap.Error(err))
		return
	} else if cleanup != nil {
		defer cleanup()
	}

	// GRPC server
	address := fmt.Sprintf("%s:%d", noti.Config.HTTP.Host, noti.Config.HTTP.Port)
	network := "tcp"

	l, err := net.Listen(network, address)
	if err != nil {
		noti.Logger.Error(
			"failed to listen to address",
			zap.Error(err),
			zap.String("network", network),
			zap.String("address", address),
		)
		return
	}

	noti.Logger.Info("🌏 start server...", zap.String("address", address))

	defer func() {
		if err1 := l.Close(); err1 != nil {
			noti.Logger.Error(
				"failed to close listener",
				zap.Error(err1),
				zap.String("network", network),
				zap.String("address", address),
			)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.Serve(l) // init gRPC server && block until the server is stopped
	}()

	select {
	case v := <-quit:
		noti.Logger.Info("signal.Notify", zap.String("signal", v.String()))
		cancel()
	case err := <-serveErr:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			noti.Logger.Error("gRPC server exited", zap.Error(err))
		}
		cancel()
	case <-ctx.Done():
		noti.Logger.Info("ctx.Done", zap.Error(ctx.Err()))
	}

	server.GracefulStop()

}
