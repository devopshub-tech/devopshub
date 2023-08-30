// internal/domain/logger.go
package domain

// Logger is the interface that defines the logging methods.
type ILogger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
	Fatal(msg string)
	Panic(msg string)

	Debugf(msg string, args ...interface{})
	Warningf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}
