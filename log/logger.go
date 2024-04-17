package log

import (
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// Logger Interface to support
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Named(string) Logger
}

type LogLevel int

const (
	DebugLogLevel LogLevel = iota
	InfoLogLevel           = iota
	ErrorLogLevel          = iota
)

type ZapLogger struct {
	*zap.SugaredLogger
}

// Logger implementation
func NewZapLogger(logPath string, level LogLevel) Logger {
	config := zap.NewDevelopmentConfig()
	if level == ErrorLogLevel {
		config = zap.NewProductionConfig()
	}

	err := os.MkdirAll(filepath.Dir(logPath), os.ModePerm)
	if err != nil {
		log.Fatalf("could not open logfile: %v", err)
	}
	config.OutputPaths = []string{logPath}

	sugar, err := config.Build()
	if err != nil {
		log.Fatalf("could not start logger: %v", err)
	}

	sugar = sugar.Named("main")

	zl := &ZapLogger{sugar.Sugar()}

	return zl
}

// Returns a new logger with the segment 'name' appeneded
func (zl *ZapLogger) Named(name string) Logger {
	return &ZapLogger{
		zl.Desugar().Named(name).Sugar(),
	}
}
