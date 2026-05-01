package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/quangnguyen1505/go-notification-system/cmd/notification/config"
	notificationapp "github.com/quangnguyen1505/go-notification-system/internal/notification/app"
	"github.com/quangnguyen1505/go-notification-system/pkg/logger"
	"github.com/sirupsen/logrus"
	"go.uber.org/automaxprocs/maxprocs"
	"google.golang.org/grpc"
)

func main() {
	// set GOMAXPROCS
	_, err := maxprocs.Set()
	if err != nil {
		slog.Error("failed set max procs", "err", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed get config", "err", err)
		return
	}

	slog.Info("⚡ init app", "name", cfg.Name, "version", cfg.Version)

	// set up logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	// integrate Logrus with the slog logger
	slog.New(logger.NewLogrusHandler(logrus.StandardLogger()))

	server := grpc.NewServer()

	if _, err := notificationapp.InitApp(cfg, server); err != nil {
		slog.Error("failed init app", "err", err)
		return
	}

	// GRPC server
	address := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	network := "tcp"

	l, err := net.Listen(network, address)
	if err != nil {
		slog.Error("failed to listen to address", "err", err, "network", network, "address", address)
		cancel()
		return
	}

	slog.Info("🌏 start server...", "address", address)

	defer func() {
		if err1 := l.Close(); err != nil {
			slog.Error("failed to close", "err", err1, "network", network, "address", address)
			<-ctx.Done()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.Serve(l)
	}()

	select {
	case v := <-quit:
		slog.Info("signal.Notify", "signal", v)
		cancel()
	case err := <-serveErr:
		if err != nil {
			slog.Error("gRPC server exited", "err", err)
		}
		cancel()
	case done := <-ctx.Done():
		slog.Info("ctx.Done", "app done", done)
	}

	server.GracefulStop()

}
