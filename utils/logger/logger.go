package logger

import (
	"context"
	"github.com/j03hanafi/seternak-backend/utils/consts"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime/debug"
	"sync"
)

type ctxKey struct{}

var (
	once   sync.Once
	logger *zap.Logger
)

// Get initializes a zap.Logger instance if it has not been initialized
// already and returns the same instance for subsequent calls.
func Get() *zap.Logger {
	once.Do(func() {
		// Console Log
		stdout := zapcore.AddSync(os.Stdout)

		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

		// File Log
		filename := "logs/app.log"
		if viper.GetString("APP_ENV") == consts.TestMode {
			filename = "../logs/app.log"
		}
		file := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filename,
			MaxSize:    5,
			MaxAge:     30,
			MaxBackups: 15,
			LocalTime:  true,
			Compress:   true,
		})

		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = "timestamp"
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		fileEncoder := zapcore.NewJSONEncoder(productionCfg)

		var gitRevision string
		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		// log to multiple destinations (console and file)
		// extra fields are added to the JSON output alone
		var (
			level = zap.NewAtomicLevelAt(zap.InfoLevel)
			core  zapcore.Core
			opts  []zap.Option
		)

		if viper.GetString("APP_ENV") != consts.ProductionMode {
			level = zap.NewAtomicLevelAt(zap.DebugLevel)
			opts = append(opts, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))
			core = zapcore.NewCore(consoleEncoder, stdout, level)
		} else {
			core = zapcore.NewTee(
				zapcore.NewCore(consoleEncoder, stdout, level),
				zapcore.NewCore(fileEncoder, file, level).
					With([]zapcore.Field{
						zap.String("gitRevision", gitRevision),
						zap.String("goVersion", buildInfo.GoVersion),
					}),
			)
		}

		logger = zap.New(core, opts...)
	})

	return logger
}

// FromCtx returns the Logger associated with the ctx. If no logger
// is associated, the default logger is returned, unless it is nil
// in which case a disabled logger is returned.
func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return zap.NewNop()
}

// WithCtx returns a copy of ctx with the Logger attached.
func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			// Do not store same logger.
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}
