package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config sets up the configurations.
type Config struct {
	Port               int    `envconfig:"PORT"`
	Env                string `envconfig:"ENV"`
	AccessSecret       string `envconfig:"ACCESS_SECRET"`
	RefreshSecret      string `envconfig:"REFRESH_SECRET"`
	AccessTokenExpire  int    `envconfig:"ACCESS_TOKEN_EXPIRE"`
	RefreshTokenExpire int    `envconfig:"REFRESH_TOKEN_EXPIRE"`
}

// LoadConfig loads the configuration from .env file in the root directory and environment variables.
func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = envconfig.Process("", &cfg)

	return &cfg, err
}
