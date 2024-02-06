package logger

import (
	"context"
	"os"

	"github.com/Dostonlv/todo.git/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Module = fx.Provide(NewLogger)

type Params struct {
	fx.In
	fx.Lifecycle
	Config config.IConfig
}

type ILogger interface {
	Info(format string, fields ...interface{})
	Warn(format string, fields ...interface{})
	Error(format string, fields ...interface{})
	Fatal(format string, fields ...interface{})
}

func NewLogger(params Params) ILogger {

	level := getLevel(params.Config)

	// write syncers
	stdoutSyncer := zapcore.Lock(os.Stdout)

	// tee core
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			stdoutSyncer,
			level,
		),
		zapcore.NewCore(getEncoder(), getWriter(params.Config), level),
	)

	// get log core
	//core := zapcore.NewCore(getEncoder(), getWriter(params.Config), level)

	// create log instance with AddCaller option.
	// AddCallerSkip option - skips stack trace where log called
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	params.Lifecycle.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				_ = log.Sync()
				return nil
			},
		},
	)

	return &logger{lg: log.Sugar(), config: params.Config}
}

type logger struct {
	lg     *zap.SugaredLogger
	config config.IConfig
}

// getEncoder returns Encoder
func getEncoder() zapcore.Encoder {

	var encoderConfig = zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getWriter returns WriteSyncer
func getWriter(config config.IConfig) zapcore.WriteSyncer {

	filename := config.GetString("logger.filename")
	if filename == "" {
		filename = "./app.log"
	}

	maxSize := config.GetInt("logger.maxSize")
	if maxSize == 0 {
		maxSize = 200
	}

	var lumberJackLogger = &lumberjack.Logger{
		Filename:   filename, //location of log file
		MaxSize:    maxSize,  //maximum size of log file in MBs, before it is rotated
		MaxBackups: 10,       //maximum no. of old files to retain
		MaxAge:     30,       //maximum number of days it will retain old files
		Compress:   false,    //whether to compress/archive old files
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getLevel(config config.IConfig) zapcore.Level {
	switch config.GetString("logger.level") {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}

func (l *logger) Info(format string, args ...interface{}) {
	l.lg.WithOptions(zap.AddCallerSkip(1)).Infof(format, args...)
}

func (l *logger) Error(format string, args ...interface{}) {
	l.lg.WithOptions(zap.AddCallerSkip(1)).Errorf(format, args...)
}

func (l *logger) Fatal(format string, args ...interface{}) {
	l.lg.WithOptions(zap.AddCallerSkip(1)).Fatalf(format, args...)
}

func (l *logger) Warn(format string, args ...interface{}) {
	l.lg.WithOptions(zap.AddCallerSkip(1)).Warnf(format, args...)
}
