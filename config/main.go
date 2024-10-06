package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	// supported databases ["mongo", "postgres"]
	DatabaseEnv    = "DATABASE"
	MongoURLEnv    = "MONGO_URL"
	PostgresURLEnv = "POSTGRES_URL"
)

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
			Url string
		}
		Postgres struct {
			Url string
		}
	}
}

// updateWithEnvs updates the config with the environment variables
// if database sets to mongo, postgres url will be ignored
// vice versa if database sets to postgres, mongo url will be ignored
func updateWithEnvs(config *Config) {
	database, ok := os.LookupEnv(DatabaseEnv)
	if ok {
		config.Database.Db = database
	}

	mongoURL, ok := os.LookupEnv(MongoURLEnv)
	if ok {
		config.Database.MongoDB.Url = mongoURL
	}

	postgresURL, ok := os.LookupEnv(PostgresURLEnv)
	if ok {
		config.Database.Postgres.Url = postgresURL
	}
}

// GetConf reads the config file and returns the Config struct
func LoadConfig() (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to read config, %v", err)
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct, %v", err)
	}

	return &c, nil
}
