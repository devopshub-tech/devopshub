// internal/core/infra/logging/writers.go
package logging

import (
	"errors"
	"fmt"
	"strings"

	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func createWriterMap(cfg domain.IConfig) (lfshook.WriterMap, error) {
	levelStr := cfg.GetLogFileLevels()

	writerLevels, err := parseLevels(levelStr)
	if err != nil {
		return nil, err
	}

	writerMap := make(lfshook.WriterMap)

	var _errors string
	for _, level := range writerLevels {
		writer, err := createLogFile(fmt.Sprintf("/%s.log", level), cfg)
		if err != nil {
			_errors = _errors + fmt.Sprintf(" [ file: %s, err: %+v ]", level, err)
			continue
		}
		writerMap[level] = writer
	}

	if _errors != "" {
		return writerMap, errors.New(_errors)
	}
	return writerMap, nil
}

func parseLevels(levelStr string) ([]logrus.Level, error) {
	levelStr = strings.TrimSpace(levelStr)
	levelStr = strings.ToLower(levelStr)
	levelStrList := strings.Split(levelStr, ",")

	var levels []logrus.Level
	for _, level := range levelStrList {
		l, err := logrus.ParseLevel(strings.TrimSpace(level))
		if err != nil {
			return nil, err
		}
		levels = append(levels, l)
	}

	return levels, nil
}
