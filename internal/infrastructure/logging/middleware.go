// pkg/logging/http_middleware.go
package logging

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HttpLogger creates a Gin middleware for logging requests.
// It logs various details of the incoming request and response.
// Configuring the logger at 'debug' level so that detailed fields are added to the log.
func (l *Logger) GinLogger() gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := float64(stop.Nanoseconds()) / 1e6
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		var level logrus.Level
		var msg string
		l.logger.SetReportCaller(true)
		if len(c.Errors) > 0 {
			level = logrus.ErrorLevel
			msg = c.Errors.ByType(gin.ErrorTypePrivate).String()
		} else {
			if statusCode >= http.StatusInternalServerError {
				level = logrus.ErrorLevel
			} else if statusCode >= http.StatusBadRequest {
				level = logrus.WarnLevel
			} else {
				level = logrus.InfoLevel
			}

			msg = fmt.Sprintf("%s - %s '%s %s' %d (%.4f ms)",
				clientIP, hostname,
				c.Request.Method, path, statusCode, latency)
		}

		// Only add additional fields and debug info if the log level is Debug
		if l.logger.GetLevel() == logrus.DebugLevel {
			entry := l.logger.WithFields(logrus.Fields{
				"hostname":   hostname,
				"statusCode": statusCode,
				"latency":    latency,
				"clientIP":   clientIP,
				"method":     c.Request.Method,
				"path":       path,
			})

			entry.Log(level, msg, c.Request.Referer(), c.Request.UserAgent(), fmt.Sprintf("size: %d", c.Writer.Size()))
		} else {
			l.logger.Log(level, msg)
		}
		l.logger.SetReportCaller(false)
	}
}
