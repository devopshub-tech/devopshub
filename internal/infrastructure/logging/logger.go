package logging

import (
	"runtime"
	"sync"

	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/config"
	"github.com/sirupsen/logrus"
)

var (
	onceInitLogger      sync.Once
	globalLogger        *logrus.Logger
	isLoggerInitialized bool
)

type Logger struct {
	logger *logrus.Logger
	config domain.IConfig
}

func NewLogger() *Logger {
	onceInitLogger.Do(initializeGlobalLogger)
	logger := &Logger{logger: globalLogger, config: config.NewConfig()}
	err := setup(logger)
	if err != nil {
		logger.Fatalf("Error setting logger: %v", err)
	}
	return logger
}

func initializeGlobalLogger() {
	globalLogger = logrus.New()
}

func (l *Logger) Debug(msg string) {
	l.log(logrus.DebugLevel, msg)
}
func (l *Logger) Info(msg string) {
	l.log(logrus.InfoLevel, msg)
}
func (l *Logger) Warning(msg string) {
	l.logger.Warning(msg)
}
func (l *Logger) Error(msg string) {
	l.log(logrus.ErrorLevel, msg)
}
func (l *Logger) Fatal(msg string) {
	l.log(logrus.FatalLevel, msg)
	l.logger.Exit(1)
}
func (l *Logger) Panic(msg string) {
	l.log(logrus.PanicLevel, msg)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(logrus.DebugLevel, format, args...)
}
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(logrus.InfoLevel, format, args...)
}
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.log(logrus.WarnLevel, format, args...)
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(logrus.ErrorLevel, format, args...)
}
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(logrus.FatalLevel, format, args...)
	l.logger.Exit(1)
}
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.log(logrus.PanicLevel, format, args...)
}

func (l *Logger) log(level logrus.Level, format string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(2)
	function := runtime.FuncForPC(pc).Name()

	fields := logrus.Fields{
		"file": file,
		"line": line,
		"func": function,
	}
	l.logger.SetReportCaller(false) // disabled to avoid clashes with fiel, line and func fields sent here

	if len(args) > 0 {
		l.logger.WithFields(fields).Logf(level, format, args...)
	} else {
		l.logger.WithFields(fields).Log(level, format)
	}
}
