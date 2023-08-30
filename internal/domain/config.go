// internal/domain/config.go
package domain

type IConfig interface {
	// Applications configs
	GetApiPort() int
	GetApiHost() string

	// Enviroment configs
	GetEnvMode() string

	// Database configs
	GetDbMongoUri() string
	GetDbName() string
	GetDbConnectionTimeout() int
	GetDbServerSelectionTimeout() int
	GetDbMaxPoolSize() int
	GetDbMinPoolSize() int
	GetDbReadPreference() string
	GetDbWriteConcern() int
	GetDbRetryWrites() bool
	GetDbReadConcern() string
	GetDbWriteConcernTimeout() int
	GetDbHeartbeatInterval() int
	GetDbHealthTimeout() int

	// Queue configs
	GetQueueUri() string
	GetQueueHealthTimeout() int

	// Logging configs
	GetLogLevel() string
	GetLogDir() string
	GetLogFormatConsole() string
	GetLogFormatFile() string
	GetLogFileLevels() string
	GetLogWriteFile() bool
	GetLogMaxSize() int
	GetLogMaxBackups() int
	GetLogMaxAge() int
	GetLogCompress() bool
}
