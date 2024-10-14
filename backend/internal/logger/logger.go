package logger

import (
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logVal = Logger{}
)

// Log Level
const (
	// LogCritLevel Panic on log
	LogCritLevel LogLevel = iota
	// Ordinary error
	LogErrLevel
	// Warn
	LogWarnLevel
	// Informational message
	LogInfoLevel
	// Messages that are used when debugging programs
	LogDebugLevel
)

type LogLevel int

type Logger struct {
	logger *zap.Logger
}

func (log *Logger) clone() *Logger {
	copy := *log
	return &copy
}

func Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}
func Init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Logger initialized")
}

// Err LogErr by log.Print
func Err(v ...interface{}) {
	logVal.WithOptions(zap.AddCallerSkip(1)).Err(v...)
}

// WithOptions Add options
func (log *Logger) WithOptions(opts ...zap.Option) *Logger {
	copy := log.clone()
	copy.logger = copy.logger.WithOptions(opts...)
	return copy
}

// Err LogErr by log.Print
func (log *Logger) Err(v ...interface{}) {
	log.WithOptions(zap.AddCallerSkip(1)).output(LogErrLevel, fmt.Sprint(v...))
}

// Output Priority
func (log *Logger) output(level LogLevel, msg string, fields ...zap.Field) {
	logIn := log.WithOptions(zap.AddCallerSkip(1)).logger
	defer logIn.Sync()
	if ce := logIn.Check(LogLevelToZapLogLevel(level), msg); ce != nil {
		ce.Write(fields...)
	}
}

func LogLevelToZapLogLevel(level LogLevel) zapcore.Level {
	switch level {
	case LogCritLevel:
		return zapcore.PanicLevel
	case LogErrLevel:
		return zapcore.ErrorLevel
	case LogWarnLevel:
		return zapcore.WarnLevel
	case LogInfoLevel:
		return zapcore.InfoLevel
	case LogDebugLevel:
		return zapcore.DebugLevel
	default:
		return zapcore.WarnLevel
	}
}

// Errf LogErr by log.Printf
func Errf(format string, v ...interface{}) {
	logVal.WithOptions(zap.AddCallerSkip(1)).Errf(format, v...)
}

// Errf LogErr by log.Printf
func (log *Logger) Errf(format string, v ...interface{}) {
	log.WithOptions(zap.AddCallerSkip(1)).outputf(LogErrLevel, format, v...)
}

// Outputf Priority
func (log *Logger) outputf(level LogLevel, format string, v ...interface{}) {
	log.WithOptions(zap.AddCallerSkip(1)).output(level, fmt.Sprintf(format, v...))
}
