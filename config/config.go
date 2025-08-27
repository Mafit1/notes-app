package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      App      `yaml:"app"`
		HTTP     HTTP     `yaml:"http"`
		Postgres Postgres `yaml:"postgres"`
		Log      Log      `yaml:"logger"`
		Auth     Auth     `yaml:"auth"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"SERVER_PORT"`
	}

	Postgres struct {
		URL            string        `env-required:"true" yaml:"url" env:"POSTGRES_URL"`
		ConnectionTime time.Duration `env-required:"true" yaml:"connect_timeout" env:"POSTGRES_CONNECT_TIMEOUT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	Auth struct {
		JWTAccessSecretKey  string        `env-required:"true" yaml:"jwt_secret_key" env:"JWT_ACCESS_SECRET_KEY"`
		JWTRefreshSecretKey string        `env-required:"true" yaml:"jwt_secret_key" env:"JWT_REFRESH_SECRET_KEY"`
		AccessTokenTTL      time.Duration `env_required:"true" yaml:"access_token_ttl" env:"JWT_ACCESS_SECRET_KEY"`
		RefreshTokenTTL     time.Duration `env_required:"true" yaml:"access_token_ttl" env:"JWT_REFRESH_SECRET_KEY"`
	}
)

func New(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config reading error: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config updating error: %w", err)
	}

	return cfg, nil
}
