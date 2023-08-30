package logging

import (
	"path/filepath"

	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/pkg/filesystem"
	loggingpkg "github.com/devopshub-tech/devopshub/pkg/logging"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func setup(l *Logger) error {
	if isLoggerInitialized {
		return nil
	}

	err := filesystem.NewDirectory(l.config.GetLogDir())
	if err != nil {
		return err
	}

	level, err := logrus.ParseLevel(l.config.GetLogLevel())
	if err != nil {
		level = logrus.InfoLevel
	}
	l.logger.SetLevel(level)

	if l.config.GetLogFormatConsole() == "json" {
		l.logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		l.logger.SetFormatter(&loggingpkg.CustomFormatter{})
	}

	if l.config.GetLogWriteFile() {
		err := configureLogFile(l)
		if err != nil {
			return err
		}
	}

	l.logger.SetReportCaller(true)
	l.logger.Info("Logging setup completed.")
	isLoggerInitialized = true
	return nil
}

func configureLogFile(l *Logger) error {
	var fileFormatter logrus.Formatter
	if l.config.GetLogFormatFile() == "json" {
		fileFormatter = &logrus.JSONFormatter{}
	} else {
		fileFormatter = &logrus.TextFormatter{}
	}

	// Configuration main log file
	mainLogFile, err := createLogFile("main.log", l.config)
	if err != nil {
		return err
	}
	l.logger.Hooks.Add(lfshook.NewHook(
		mainLogFile,
		fileFormatter,
	))

	// Configuration log file by level
	pathMap, err := createWriterMap(l.config)
	if err != nil {
		return err
	}

	l.logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		fileFormatter,
	))

	return nil
}

func createLogFile(filename string, config domain.IConfig) (*lumberjack.Logger, error) {
	filePath := filepath.Join(config.GetLogDir(), filename)
	logFile := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    config.GetLogMaxSize(),
		MaxBackups: config.GetLogMaxBackups(),
		MaxAge:     config.GetLogMaxAge(),
		Compress:   config.GetLogCompress(),
	}
	return logFile, nil
}
