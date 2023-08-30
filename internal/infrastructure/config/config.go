package config

import (
	"log"
	"sync"

	"github.com/devopshub-tech/devopshub/internal/domain"
	configPkg "github.com/devopshub-tech/devopshub/pkg/config"
)

var onceLoadEnv sync.Once

type Config struct {
	cfg *configPkg.Config
}

func NewConfig() domain.IConfig {
	c := &Config{cfg: configPkg.NewConfig()}
	c.loadEnvOnce()
	return c
}

func (ca *Config) loadEnv() {
	if err := ca.cfg.Load(); err != nil {
		log.Printf("Could not load .env file: %v", err)
	}
	log.Println("Configuration loaded from .env file.")
}

func (ca *Config) loadEnvOnce() {
	onceLoadEnv.Do(ca.loadEnv)
}

func (ca *Config) Reload() {
	ca.loadEnv()
}

func (ca *Config) GetApiPort() int {
	return ca.cfg.GetEnvAsInt("API_PORT", &configPkg.GetEnvOptions{Fallback: "8080"})
}

func (ca *Config) GetApiHost() string {
	return ca.cfg.GetEnv("API_HOST", &configPkg.GetEnvOptions{Fallback: "0.0.0.0"})
}

func (ca *Config) GetEnvMode() string {
	return ca.cfg.GetEnv("ENV_MODE", &configPkg.GetEnvOptions{Fallback: "development"})
}

func (ca *Config) GetDbMongoUri() string {
	return ca.cfg.GetEnv("DB_MONGO_URI", &configPkg.GetEnvOptions{Required: true})
}

func (ca *Config) GetDbName() string {
	return ca.cfg.GetEnv("DB_NAME", &configPkg.GetEnvOptions{Fallback: "devopshubdb"})
}

func (ca *Config) GetDbConnectionTimeout() int {
	return ca.cfg.GetEnvAsInt("DB_CONNECTION_TIMEOUT", &configPkg.GetEnvOptions{Fallback: "10"})
}

func (ca *Config) GetDbServerSelectionTimeout() int {
	return ca.cfg.GetEnvAsInt("DB_SERVER_SELECTION_TIMEOUT", &configPkg.GetEnvOptions{Fallback: "30"})
}

func (ca *Config) GetDbMaxPoolSize() int {
	return ca.cfg.GetEnvAsInt("DB_MAX_POOL_SIZE", &configPkg.GetEnvOptions{Fallback: "100"})
}

func (ca *Config) GetDbMinPoolSize() int {
	return ca.cfg.GetEnvAsInt("DB_MIN_POOL_SIZE", &configPkg.GetEnvOptions{Fallback: "1"})
}

func (ca *Config) GetDbReadPreference() string {
	return ca.cfg.GetEnv("DB_READ_PREFERENCE", &configPkg.GetEnvOptions{Fallback: "primary"})
}

func (ca *Config) GetDbWriteConcern() int {
	return ca.cfg.GetEnvAsInt("DB_WRITE_CONCERN", &configPkg.GetEnvOptions{Fallback: "1"})
}

func (ca *Config) GetDbRetryWrites() bool {
	return ca.cfg.GetEnvAsBool("DB_RETRY_WRITES", &configPkg.GetEnvOptions{Fallback: "true"})
}

func (ca *Config) GetDbReadConcern() string {
	return ca.cfg.GetEnv("DB_READ_CONCERN", &configPkg.GetEnvOptions{Fallback: "majority"})
}

func (ca *Config) GetDbWriteConcernTimeout() int {
	return ca.cfg.GetEnvAsInt("DB_WRITE_CONCERN_TIMEOUT", &configPkg.GetEnvOptions{Fallback: "10000"})
}

func (ca *Config) GetDbHeartbeatInterval() int {
	return ca.cfg.GetEnvAsInt("DB_HEARTBEAT_INTERVAL", &configPkg.GetEnvOptions{Fallback: "10"})
}

func (ca *Config) GetDbHealthTimeout() int {
	return ca.cfg.GetEnvAsInt("DB_HEALTH_TIMEOUT", &configPkg.GetEnvOptions{Fallback: "5"})
}

func (ca *Config) GetQueueUri() string {
	return ca.cfg.GetEnv("QUEUE_URI", &configPkg.GetEnvOptions{Required: true})
}

func (ca *Config) GetQueueHealthTimeout() int {
	return ca.cfg.GetEnvAsInt("QUEUE_HEALTH_TIMEOUT", &configPkg.GetEnvOptions{Fallback: "10"})
}

func (ca *Config) GetLogLevel() string {
	return ca.cfg.GetEnv("LOG_LEVEL", &configPkg.GetEnvOptions{Fallback: "info"})
}

func (ca *Config) GetLogDir() string {
	return ca.cfg.GetEnv("LOG_DIR", &configPkg.GetEnvOptions{Fallback: "logs"})
}

func (ca *Config) GetLogFormatConsole() string {
	return ca.cfg.GetEnv("LOG_FORMAT_CONSOLE", &configPkg.GetEnvOptions{Fallback: "text"})
}

func (ca *Config) GetLogFormatFile() string {
	return ca.cfg.GetEnv("LOG_FORMAT_FILE", &configPkg.GetEnvOptions{Fallback: "json"})
}

func (ca *Config) GetLogFileLevels() string {
	return ca.cfg.GetEnv("LOG_FILE_LEVELS", &configPkg.GetEnvOptions{Fallback: "error,warning"})
}

func (ca *Config) GetLogWriteFile() bool {
	return ca.cfg.GetEnvAsBool("LOG_WRITE_FILE", &configPkg.GetEnvOptions{Fallback: "true"})
}
func (ca *Config) GetLogMaxSize() int {
	return ca.cfg.GetEnvAsInt("LOG_MAX_SIZE", &configPkg.GetEnvOptions{Fallback: "100"})
}
func (ca *Config) GetLogMaxBackups() int {
	return ca.cfg.GetEnvAsInt("LOG_MAX_BACKUPS", &configPkg.GetEnvOptions{Fallback: "6"})
}
func (ca *Config) GetLogMaxAge() int {
	return ca.cfg.GetEnvAsInt("LOG_MAX_AGE", &configPkg.GetEnvOptions{Fallback: "30"})
}
func (ca *Config) GetLogCompress() bool {
	return ca.cfg.GetEnvAsBool("LOG_COMPRESS", &configPkg.GetEnvOptions{Fallback: "true"})
}
