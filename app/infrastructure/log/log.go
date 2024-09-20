package log

import (
	"errors"
	"fmt"
	"os"
	"unsafe"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"jade-mes/app/infrastructure/constant"
	"jade-mes/config"
)

const (
	// LevelDebug ...
	LevelDebug = "debug"
	// LevelInfo ...
	LevelInfo = "info"
	// LevelWarn ...
	LevelWarn = "warn"
	// LevelError ...
	LevelError = "error"
)

// Field ...
type Field = zapcore.Field

// Int64 ...
var Int64 = zap.Int64

// Int ...
var Int = zap.Int

// String ...
var String = zap.String

// Reflect ...
var Reflect = zap.Reflect

// Object ...
var Object = zap.Object

// Bool ...
var Bool = zap.Bool

// Array ...
var Array = zap.Array

// Any ...
var Any = zap.Any

// Err ...
var Err = zap.Error

// LoggerConfiguration ...
type LoggerConfiguration struct {
	EnableConsole  bool
	ConsoleJSON    bool
	ConsoleLevel   string
	EnableFile     bool
	FileJSON       bool
	FileLevel      string
	FileLocation   string
	AccessLocation string
}

// Logger ...
type Logger struct {
	zap          *zap.Logger
	consoleLevel zap.AtomicLevel
	fileLevel    zap.AtomicLevel
}

var logger *Logger
var accessLog *Logger

func getZapLevel(level string) zapcore.Level {
	switch level {
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func makeEncoder(json bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if json {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func loadLogConfig(config *viper.Viper) *LoggerConfiguration {
	var logConfig LoggerConfiguration
	config.UnmarshalKey("log", &logConfig)

	return &logConfig
}

func init() {
	println("initing logger...")
	appConfig := config.GetConfig()
	config := loadLogConfig(appConfig)

	logger = initLogger(config)

	// override file location for access log
	if config.AccessLocation != "" {
		config.FileLocation = config.AccessLocation
	}
	accessLog = initLogger(config)
}

func initLogger(config *LoggerConfiguration) *Logger {
	cores := []zapcore.Core{}
	logger := &Logger{
		consoleLevel: zap.NewAtomicLevelAt(getZapLevel(config.ConsoleLevel)),
		fileLevel:    zap.NewAtomicLevelAt(getZapLevel(config.FileLevel)),
	}

	if config.EnableConsole {
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(makeEncoder(config.ConsoleJSON), writer, logger.consoleLevel)
		cores = append(cores, core)
	}

	if config.EnableFile {
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename: config.FileLocation,
			MaxSize:  500,
			MaxAge:   15, // keep history logs for 15 days
			Compress: true,
		})
		core := zapcore.NewCore(makeEncoder(config.FileJSON), writer, logger.fileLevel)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	logger.zap = zap.New(combinedCore,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
		zap.AddStacktrace(zap.FatalLevel),
	)

	return logger
}

// GetLoggers returns logger
func GetLoggers() (*Logger, *Logger) {
	return logger, accessLog
}

// GetLogger ...
func GetLogger() *Logger {
	return logger
}

// ChangeLevels ...
func (l *Logger) ChangeLevels(config *LoggerConfiguration) {
	l.consoleLevel.SetLevel(getZapLevel(config.ConsoleLevel))
	l.fileLevel.SetLevel(getZapLevel(config.FileLevel))
}

// SetConsoleLevel ...
func (l *Logger) SetConsoleLevel(level string) {
	l.consoleLevel.SetLevel(getZapLevel(level))
}

// With ...
func (l *Logger) With(fields ...Field) *Logger {
	newlogger := *l
	newlogger.zap = newlogger.zap.With(fields...)
	return &newlogger
}

// Debug ...
func (l *Logger) Debug(message string, fields ...Field) {
	for i, field := range fields {
		if (*[2]uintptr)(unsafe.Pointer(&field.Interface))[1] == 0 {
			fields[i].Interface = errors.New("")
		}
	}
	l.With(
		zap.Namespace(constant.NameSpace),
	).zap.Debug(message, fields...)
}

// Info ...
func (l *Logger) Info(message string, fields ...Field) {
	for i, field := range fields {
		if (*[2]uintptr)(unsafe.Pointer(&field.Interface))[1] == 0 {
			fields[i].Interface = errors.New("")
		}
	}
	l.With(
		zap.Namespace(constant.NameSpace),
	).zap.Info(message, fields...)
}

// AccessLog ...
func (l *Logger) AccessLog(message string, fields ...Field) {
	l.With(
		zap.Namespace(constant.NameSpace),
	).zap.Info(message, fields...)
}

// Warn ...
func (l *Logger) Warn(message string, fields ...Field) {
	for i, field := range fields {
		if (*[2]uintptr)(unsafe.Pointer(&field.Interface))[1] == 0 {
			fields[i].Interface = errors.New("")
		}
	}
	l.With(
		zap.Namespace(constant.NameSpace),
	).zap.Warn(message, fields...)
}

// Error ...
func (l *Logger) Error(message string, fields ...Field) {
	for i, field := range fields {
		if (*[2]uintptr)(unsafe.Pointer(&field.Interface))[1] == 0 {
			fields[i].Interface = errors.New("")
		}
	}
	l.With(
		zap.Namespace(constant.NameSpace),
	).zap.Error(message, fields...)
}

// Critical ...
func (l *Logger) Critical(message string, fields ...Field) {
	l.With(
		zap.Namespace(constant.NameSpace),
	).zap.Error(message, fields...)
}

// Printf ...
func (l *Logger) Printf(message string, items ...interface{}) {
	l.With(
		zap.Namespace(constant.NameSpace),
	).zap.Info(fmt.Sprintf(message, items...))
}
