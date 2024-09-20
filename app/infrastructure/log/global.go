package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *Logger
var accessLogger *Logger

// InitGlobalLogger ...
func InitGlobalLogger(logger *Logger) {
	globalLogger = logger
	Debug = globalLogger.Debug
	Info = globalLogger.Info
	Warn = globalLogger.Warn
	Error = globalLogger.Error
	Critical = globalLogger.Critical
	With = globalLogger.With
	Printf = globalLogger.Printf
}

// InitAccessLogger ...
func InitAccessLogger(logger *Logger) {
	accessLogger = logger
	AccessLog = accessLogger.AccessLog
}

// RedirectStdLog ...
func RedirectStdLog(logger *Logger) {
	zap.RedirectStdLogAt(logger.zap.With(zap.String("source", "stdlog")), zapcore.InfoLevel)
}

// logFunc ...
type logFunc func(string, ...Field)

type printFunc func(string, ...interface{})

type withFunc func(...Field) *Logger

// GloballyDisableDebugLogForTest ...
func GloballyDisableDebugLogForTest() {
	globalLogger.consoleLevel.SetLevel(zapcore.ErrorLevel)
}

// GloballyEnableDebugLogForTest ...
func GloballyEnableDebugLogForTest() {
	globalLogger.consoleLevel.SetLevel(zapcore.DebugLevel)
}

// With ...
var With withFunc = defaultWithLog

// Debug ...
var Debug logFunc = defaultDebugLog

// Info ...
var Info logFunc = defaultInfoLog

// AccessLog ...
var AccessLog logFunc = defaultAccessLog

// Warn ...
var Warn logFunc = defaultWarnLog

// Error ...
var Error logFunc = defaultErrorLog

// Critical ...
var Critical logFunc = defaultCriticalLog

// Printf ...
var Printf printFunc = defaultPrintLog
