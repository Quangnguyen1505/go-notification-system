package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/quangnguyen1505/go-notification-system/cmd/proxy/config"
	"github.com/quangnguyen1505/go-notification-system/pkg/logger"
	"github.com/quangnguyen1505/go-notification-system/proto/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGateway(
	ctx context.Context,
	cfg *config.Config,
	opts []gwruntime.ServeMuxOption,
) (http.Handler, error) {
	notificationHost := cfg.NotificationHost
	if notificationHost == "" || notificationHost == "0.0.0.0" {
		notificationHost = "127.0.0.1"
	}
	notificationEndPoint := fmt.Sprintf("%s:%d", notificationHost, cfg.NotificationPort)
	// userEndPoint := fmt.Sprintf("%s:%d", cfg.UserHost, cfg.UserPort)

	mux := gwruntime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := gen.RegisterNotificationServiceHandlerFromEndpoint(ctx, mux, notificationEndPoint, dialOpts)
	if err != nil {
		return nil, err
	}

	// err = gen.RegisterUserServiceHandlerFromEndpoint(ctx, mux, userEndPoint, dialOpts)
	// if err != nil {
	// 	return nil, err
	// }

	return mux, nil
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)

				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	headers := []string{"*"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))

	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))

	slog.Info("preflight request", "http_path", r.URL.Path)
}

func withLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Run request", "http_method", r.Method, "http_url", r.URL)

		h.ServeHTTP(w, r)
	})
}

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// set up logrus logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	// intergrate logrus with the slog logger
	slog.New(logger.NewLogrusHandler(logrus.StandardLogger()))

	mux := http.NewServeMux()

	gw, err := NewGateway(ctx, cfg, nil)
	if err != nil {
		slog.Error("failed to create a new gateway", "err", err)
		return
	}

	mux.Handle("/", gw)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: withLogger(allowCORS(mux)),
	}

	go func() {
		<-ctx.Done()
		slog.Info("shutting down the server")

		if err := s.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown the server", "err", err)
		}
	}()

	slog.Info("start listening...", "address", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to listen and serve", "err", err)
	}
}
