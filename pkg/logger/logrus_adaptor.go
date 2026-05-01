package logger

// refs:
// https://josephwoodward.co.uk/2022/11/slog-structured-logging-proposal
// https://thedevelopercafe.com/articles/logging-in-go-with-slog-a7bb489755c2

import (
	"context"
	"log/slog"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogrusHandler struct {
	logger *logrus.Logger
}

func NewLogrusHandler(logger *logrus.Logger) *LogrusHandler {
	return &LogrusHandler{
		logger: logger,
	}
}

func ConvertLogLevel(level string) logrus.Level {
	var l logrus.Level

	switch strings.ToLower(level) {
	case "error":
		l = logrus.ErrorLevel
	case "warn", "warm":
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	default:
		l = logrus.InfoLevel
	}

	return l
}

func (h *LogrusHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// support all logging levels
	return true
}

func (h *LogrusHandler) Handle(_ context.Context, rec slog.Record) error {
	fields := make(map[string]interface{}, rec.NumAttrs())

	rec.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	entry := h.logger.WithFields(fields)

	switch rec.Level {
	case slog.LevelDebug:
		entry.Debug(rec.Message)
	case slog.LevelInfo:
		entry.Info(rec.Message)
	case slog.LevelWarn:
		entry.Warn(rec.Message)
	case slog.LevelError:
		entry.Error(rec.Message)
	}

	return nil
}

func (h *LogrusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// not implemented for brevity
	return h
}

func (h *LogrusHandler) WithGroup(name string) slog.Handler {
	// logrus không có concept group → bạn có thể ignore hoặc prefix key
	return h
}
