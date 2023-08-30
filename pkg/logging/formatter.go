// pkg/logging/formatter.go
package logging

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var timeFormat = "2006-01-02 15:04:05"

type CustomFormatter struct{}

// Format formats the log entry in a custom format.
// It includes the timestamp, caller's short file name and line number,
// log level, and log message.
// Example: 2023-08-22 14:47:13,678 - http_server.go - [line:58] - INFO - Server exited
func (mf *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	levelStr := strings.ToUpper(entry.Level.String())
	filename := "FIELD_LOG_NOT_FOUND"
	line := -1

	if entry.Logger.ReportCaller {
		filename = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
	} else {
		if entry.Data != nil {
			if fileVal, ok := entry.Data["file"]; ok {
				filename, _ = fileVal.(string)
				filename = filepath.Base(filename)
			}
			if lineVal, ok := entry.Data["line"]; ok {
				if lineNum, err := strconv.Atoi(fmt.Sprintf("%v", lineVal)); err == nil {
					line = lineNum
				}
			}
		}
	}

	b.WriteString(fmt.Sprintf("%s - %s - [line:%d] - %s - %s\n",
		entry.Time.Format(timeFormat), filename,
		line, levelStr, entry.Message))
	return b.Bytes(), nil
}
