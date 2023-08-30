// pkg/config/config.go
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type GetEnvOptions struct {
	Required bool
	Fallback string
}

type Config struct{}

func NewConfig(filenames ...string) *Config {
	return &Config{}
}

func (c *Config) Load(filenames ...string) error {
	if err := godotenv.Load(filenames...); err != nil {
		return err
	}
	return nil
}

func (c *Config) GetEnv(key string, options ...*GetEnvOptions) string {
	opt := mergeGetEnvOptions(options...)

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if opt.Required {
		panic("'" + key + "' environment variable is required")
	}

	return opt.Fallback
}

func (c *Config) GetEnvAsBool(key string, options ...*GetEnvOptions) bool {
	return getEnvAsType(key, parseBool, options...).(bool)
}

func (c *Config) GetEnvAsInt(key string, options ...*GetEnvOptions) int {
	return getEnvAsType(key, parseInt, options...).(int)
}

func mergeGetEnvOptions(opts ...*GetEnvOptions) *GetEnvOptions {
	g := GetEnvOptions{}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if len(opt.Fallback) > 0 {
			g.Fallback = opt.Fallback
		}
		if opt.Required {
			g.Required = opt.Required
		}
	}

	return &g
}

func parseBool(value string) (interface{}, error) {
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return nil, err
	}
	return boolValue, nil
}

func parseInt(value string) (interface{}, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	return intValue, nil
}

func getEnvAsType(key string, converter func(string) (interface{}, error), options ...*GetEnvOptions) interface{} {
	opt := mergeGetEnvOptions(options...)

	value, exists := os.LookupEnv(key)
	if exists {
		convertedValue, err := converter(value)
		if err == nil {
			return convertedValue
		}
		log.Fatalf("Environment variable '%s' cannot be converted: %v", key, err)
	}

	if opt.Required {
		log.Fatalf("Environment variable '%s' is required", key)
	}

	convertedFallback, err := converter(opt.Fallback)
	if err != nil {
		log.Fatalf("The fallback value for environment variable '%s' cannot be converted: %v", key, err)
	}
	return convertedFallback
}
