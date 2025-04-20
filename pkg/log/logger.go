package log

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

type ctxKeyType string

var LoggerCtxKey ctxKeyType = "logger"

type Logger struct {
	zapLogger *zap.Logger
	level     *zap.AtomicLevel
}

type config struct {
	zap.Config
	callerSkip int
}

type option func(*config) error

var defaultLogger *Logger

func init() {
	var err error
	defaultLogger, err = NewLogger()
	if err != nil {
		panic(err)
	}
}

func InitializeLogger(options ...option) (err error) {
	defaultLogger, err = NewLogger(options...)
	if err != nil {
		return err
	}

	return nil
}

func NewLogger(opts ...option) (*Logger, error) {
	conf := &config{zap.NewProductionConfig(), 1}
	opts = append(opts, timestamp())
	for _, opt := range opts {
		err := opt(conf)

		if err != nil {
			return nil, err
		}
	}

	zapLogger, err := conf.Build(
		zap.AddStacktrace(zapcore.DPanicLevel),
		zap.AddCallerSkip(conf.callerSkip),
	)
	if err != nil {
		return nil, err
	}

	l := &Logger{
		zapLogger: zapLogger,
		level:     &conf.Level,
	}

	return l, nil
}
func DefaultLogger() *Logger {
	return defaultLogger
}

// Sync flushes any buffered log entries
func (logger *Logger) Sync() error {
	return logger.zapLogger.Sync()
}

func parseLevelString(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

// timestamp sets the time key name and format
func timestamp() option {
	return func(conf *config) error {
		conf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)
		conf.EncoderConfig.TimeKey = "@timestamp"
		return nil
	}
}

// ContextWithLogger creates a new context derived from the provided context with a Logger instance
func ContextWithLogger(ctx context.Context, l *Logger) context.Context {
	switch c := ctx.(type) {
	case *gin.Context:
		c.Set(string(LoggerCtxKey), l)
		return c
	default:
		return context.WithValue(ctx, string(LoggerCtxKey), l)
	}
}

func Level(level string) option {
	return func(conf *config) error {
		lvl := parseLevelString(level)
		conf.Level.SetLevel(lvl)

		return nil
	}
}

func Formatter(format string) option {
	return func(conf *config) error {
		switch format {
		case "text", "ascii", "terminal":
			conf.Encoding = "console"
		default:
			conf.Encoding = "json"

		}

		return nil
	}
}
