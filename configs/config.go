package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		Ssl      string
	}
	Auth struct {
		AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
		RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	}
}

var AppConfig Config

func LoadConfig(path string) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
	}

	if err := viper.UnmarshalKey(env, &AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %s\n", err)
	}
}
