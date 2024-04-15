package log

import (
	"go.uber.org/zap"
	"log"
)

// Logger Interface to support
type Logger interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
}

type LogLevel int

const (
	DebugLogLevel LogLevel = iota
	InfoLogLevel           = iota
	ErrorLogLevel          = iota
)

type ZapLogger struct {
	sugar *zap.SugaredLogger
}

// Logger implementation
func NewZapLogger(logPath string, level LogLevel) Logger {
	config := zap.NewDevelopmentConfig()
	if level == ErrorLogLevel {
		config = zap.NewProductionConfig()
	}
	config.OutputPaths = []string{logPath}

	sugar, err := config.Build()
	if err != nil {
		log.Fatalf("could not start logger: %v", err)
	}

	zl := &ZapLogger{
		sugar: sugar.Sugar(),
	}

	return zl
}

func (zl *ZapLogger) Debug(template string, args ...interface{}) {
	zl.sugar.Debugf(template, args)
}

func (zl *ZapLogger) Info(template string, args ...interface{}) {
	zl.sugar.Infof(template, args)
}

func (zl *ZapLogger) Error(template string, args ...interface{}) {
	zl.sugar.Errorf(template, args)
}

func (zl *ZapLogger) Fatal(template string, args ...interface{}) {
	zl.sugar.Fatalf(template, args)
}
