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
func GetLogger(cfg *config.Config) (log Logger, err error) {
	if cfg.Env.Environment == "test" {
		logger := zap.NewNop()
		defer logger.Sync()
		sugar := logger.Sugar()
		return sugar, nil
	} else {
		zapCfg := zap.NewProductionConfig()
		zapCfg.OutputPaths = []string{
			cfg.Logger.Path,
			"stdout",
		}

		zapCfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		zapCfg.Level.SetLevel(getLevel(cfg.Logger.Level))

		logger, err := zapCfg.Build()
		if err != nil {
			return nil, err
		}

		defer logger.Sync()
		sugar := logger.Sugar()
		return sugar, nil
	}
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
