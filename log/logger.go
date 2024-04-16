package log

import (
	"go.uber.org/zap"
	"log"
)

// Logger Interface to support
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
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

func (zl *ZapLogger) Debugf(template string, args ...interface{}) {
	zl.sugar.Debugf(template, args)
}

func (zl *ZapLogger) Infof(template string, args ...interface{}) {
	zl.sugar.Infof(template, args)
}

func (zl *ZapLogger) Errorf(template string, args ...interface{}) {
	zl.sugar.Errorf(template, args)
}

func (zl *ZapLogger) Fatalf(template string, args ...interface{}) {
	zl.sugar.Fatalf(template, args)
}
