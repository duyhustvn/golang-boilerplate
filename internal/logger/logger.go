package logger

import (
	"boilerplate/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is interface for log
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Panic(...interface{})
	Fatal(...interface{})

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})
	Fatalf(string, ...interface{})
}

// GetLogger return logger
func GetLogger(env *config.Config) (log Logger, err error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		env.Logger.Path,
		"stdout",
	}

	cfg.Level.SetLevel(getLevel(env.Logger.Level))

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	defer logger.Sync()
	sugar := logger.Sugar()
	return sugar, nil
}

func getLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "PANIC":
		return zapcore.PanicLevel
	case "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
