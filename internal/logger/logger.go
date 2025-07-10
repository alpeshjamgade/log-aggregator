package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log-aggregator/config"
	"log-aggregator/internal/constants"
	"os"
	"path/filepath"
	"time"
)

var logger *zap.SugaredLogger
var loggerCount int32

func CreateFileLoggerWithCtx(ctx context.Context) *zap.SugaredLogger {
	if logger == nil {
		ensureLogDirExists()

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

		cfg := zap.Config{
			Level:             GetLevel(),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          nil,
			Encoding:          config.LogEncoding,
			EncoderConfig:     encoderCfg,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stdout"},
		}

		baseLogger := zap.Must(cfg.Build())

		// Add lumberjack writer to core manually
		fileName := fmt.Sprintf("%s_%v.log", constants.ServiceName, time.Now().Format("20060102150405"))
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:  config.LogFilePath + fileName,
			MaxSize:   config.LogFileMaxSize,
			LocalTime: true,
		})

		coreWithFile := zapcore.NewTee(
			baseLogger.Core(), // original core (stdout)
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderCfg),
				fileWriter,
				cfg.Level,
			),
		)

		zapLogger := zap.New(coreWithFile, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		logger = zapLogger.Sugar()
		loggerCount++
		logger.Infof("Logger initialized with zap.Config + file rotation, count: %d", loggerCount)

	}

	if ctx != nil && ctx.Value(constants.TraceID) != nil {
		traceId := ctx.Value(constants.TraceID).(string)
		return logger.WithOptions(zap.Fields(zap.String(constants.TraceID, traceId), zap.String(constants.Service, constants.ServiceName)))
	}

	return logger
}

func CreateLoggerWithCtx(ctx context.Context) *zap.SugaredLogger {
	if logger == nil {
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

		cfg := zap.Config{
			Level:             GetLevel(),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          nil,
			Encoding:          config.LogEncoding,
			EncoderConfig:     encoderCfg,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stdout"},
		}

		logger = zap.Must(cfg.Build()).Sugar()
		loggerCount++
		logger.Debugf("create new loggger, count: %d", loggerCount)

	}

	if ctx != nil && ctx.Value(constants.TraceID) != nil {
		traceId := ctx.Value(constants.TraceID).(string)
		return logger.WithOptions(zap.Fields(zap.String(constants.TraceID, traceId), zap.String(constants.Service, constants.ServiceName)))
	}

	return logger
}

func CreateLogger() *zap.SugaredLogger {
	if logger == nil {

		var fileName, filePath string
		fileName = fmt.Sprintf("%s_%v.log", constants.ServiceName, time.Now().Format("20060102150405"))
		filePath = "/tmp"

		l := &lumberjack.Logger{
			Filename:  filePath + fileName,
			MaxSize:   500,
			LocalTime: true,
		}

		fileWriter := zapcore.AddSync(l)
		consoleWriter := zapcore.AddSync(os.Stderr)

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

		core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), consoleWriter, GetLevel()),
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), fileWriter, GetLevel()),
		)

		zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		logger = zapLogger.Sugar()
		loggerCount++
		logger.Infof("Logger initialized with rotation, count: %d", loggerCount)
	}
	return logger
}

func GetLevel() zap.AtomicLevel {
	switch config.LogLevel {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zap.PanicLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)

	}
	return zap.NewAtomicLevelAt(zap.InfoLevel)
}

func ensureLogDirExists() {
	logDir := filepath.Dir(config.LogFilePath)
	_ = os.MkdirAll(logDir, os.ModePerm)
}
