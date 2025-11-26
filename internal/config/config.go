package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Web `yaml:"web"`
	DB  `yaml:"db"`
}

type Web struct {
	Address         string        `yaml:"address" env-default:"0.0.0.0:8080"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env-default:"10s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"120s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"5s"`
}

type DB struct {
	User         string `yaml:"user" env-default:"postgres"`
	Password     string `yaml:"password" env-default:"postgres"`
	Host         string `yaml:"host" env-default:"localhost"`
	Name         string `yaml:"name" env-default:"postgres"`
	MaxIdleConns int    `yaml:"max_idle_conns" env-default:"0"`
	MaxOpenConns int    `yaml:"max_open_conns" env-default:"0"`
	DisableTLS   bool   `yaml:"disable_tls" env-default:"true"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		return nil, fmt.Errorf("empty config path value")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %s", err)
	}

	return &cfg, nil
}
