package logging

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type ShortFilePathHook struct{}

func (hook *ShortFilePathHook) Fire(entry *logrus.Entry) error {
	entry.Data["shortFilePath"] = filepath.Base(entry.Caller.File)
	return nil
}

func (hook *ShortFilePathHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
