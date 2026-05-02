package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	configs "github.com/quangnguyen1505/go-notification-system/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger(config configs.Log) *LoggerZap {
	loglevel := config.Loglevel
	// debug -> info -> warn -> error -> fatal -> panic
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	encoder := GetEncoderLog()
	hook := lumberjack.Logger{
		Filename:   config.File_name,
		MaxSize:    config.Max_size, // megabytes
		MaxBackups: config.Max_backups,
		MaxAge:     config.Max_age,  //days
		Compress:   config.Compress, // disabled by default
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(&hook)),
		level,
	)
	return &LoggerZap{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

// format logger
func GetEncoderLog() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	//conver 1724381960.2213519 -> 2024-08-23T09:59:20.219+0700
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	//change title ts -> time
	encoderConfig.TimeKey = "time"

	// change level from info -> INFO
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	//"caller":"cli/main.logs.go:21"
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
