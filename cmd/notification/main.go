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
	notificationapp "github.com/quangnguyen1505/go-notification-system/internal/notification/app"
	"github.com/quangnguyen1505/go-notification-system/pkg/logger"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var Logger *logger.LoggerZap

func main() {
	// set GOMAXPROCS -> get default from cpu core -> ex: 4 cores -> GOMAXPROCS=8
	_, err := maxprocs.Set()
	if err != nil {
		log.Printf("failed set max procs: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Printf("failed get config: %v", err)
		return
	}

	Logger = logger.NewLogger(cfg.Log)
	Logger.Info("⚡ init app", zap.String("name", cfg.Name), zap.String("version", cfg.Version))

	// set up logrus
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	// logrus.SetOutput(os.Stdout)
	// logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	server := grpc.NewServer()

	if _, cleanup, err := notificationapp.InitApp(cfg, Logger, server); err != nil {
		Logger.Error("failed init app", zap.Error(err))
		return
	} else if cleanup != nil {
		defer cleanup()
	}

	// GRPC server
	address := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	network := "tcp"

	l, err := net.Listen(network, address)
	if err != nil {
		Logger.Error(
			"failed to listen to address",
			zap.Error(err),
			zap.String("network", network),
			zap.String("address", address),
		)
		return
	}

	Logger.Info("🌏 start server...", zap.String("address", address))

	defer func() {
		if err1 := l.Close(); err1 != nil {
			Logger.Error(
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
		Logger.Info("signal.Notify", zap.String("signal", v.String()))
		cancel()
	case err := <-serveErr:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			Logger.Error("gRPC server exited", zap.Error(err))
		}
		cancel()
	case <-ctx.Done():
		Logger.Info("ctx.Done", zap.Error(ctx.Err()))
	}

	server.GracefulStop()

}
