package config

import (
	"log"

	"github.com/spf13/viper"
)

func fallbackConfigs() {
	viper.SetDefault("DB_CONNECTION_FORMAT", "%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true")
	viper.SetDefault("DB", "user")
	viper.SetDefault("DB_USER", "alex")
	viper.SetDefault("DB_PASSWORD", "alexsecret")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("MIGRATIONS_DIR", "/internal/migrations")
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("MAX_CONNECTIONS", 100)
	viper.SetDefault("MAX_IDLE_CONNECTIONS", 10)
	viper.SetDefault("MAX_LIFETIME", 1)
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_EXPIRY", 60)
	viper.SetDefault("NATS_URL", "nats://127.0.0.1:4222")
	viper.SetDefault("JWT_SECRET", "testsecret")
}

// Configurations app configs from env file, env params or fallback configs
type Configurations struct {
	DBConnectionFormat string `mapstructure:"DB_CONNECTION_FORMAT"`
	DBName             string `mapstructure:"DB"`
	DBUser             string `mapstructure:"DB_USER"`
	DBPassword         string `mapstructure:"DB_PASSWORD"`
	Host               string `mapstructure:"DB_HOST"`
	MigrationsDir      string `mapstructure:"MIGRATIONS_DIR"`
	Port               string `mapstructure:"PORT"`
	MaxConnections     int    `mapstructure:"MAX_CONNECTIONS"`
	MaxIdleConnections int    `mapstructure:"MAX_IDLE_CONNECTIONS"`
	MaxLifetime        int    `mapstructure:"MAX_LIFETIME"`
	RedisAddress       string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword      string `mapstructure:"REDIS_PASSWORD"`
	RedisDB            int    `mapstructure:"REDIS_DB"`
	NatsURL            string `mapstructure:"NATS_URL"`
	JWTSecret          string `mapstructure:"JWT_SECRET"`
}

var CurrentConfigs Configurations

// initConfig reads in config file and ENV variables if set.
func init() {

	var err error

	CurrentConfigs, err = LoadConfig("./")
	if err != nil {
		log.Fatalf("something went wrong setting up configs: %+v", err)
	}
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Configurations, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fallbackConfigs()
		} else {
			return
		}
	}

	err = viper.Unmarshal(&config)

	return
}
