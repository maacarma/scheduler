package utils

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

// defaultConfig is the default configuration for the application
var defaultConfig = []byte(`
application:
  name: "scheduler"
  port: ":7187"
  env: "development"
database:
  db: "postgres"
  mongodb:
    connString: "mongodb://localhost:27017"
  postgres:
    connString: "postgres://scheduler:scheduler@localhost:5432/scheduler"
`)

// Config struct holds the application configuration
type Config struct {
	Application struct {
		Name string
		Port string
		Env  string
	}
	Database struct {
		Db      string
		MongoDB struct {
			ConnString string
		}
		Postgres struct {
			ConnString string
		}
	}
}

// GetConf reads the config file and returns the Config struct
func LoadConfig() (*Config, error) {
	viper.SetConfigType("YAML")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		return nil, fmt.Errorf("unable to read config, %v", err)
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct, %v", err)
	}

	return &c, nil
}
