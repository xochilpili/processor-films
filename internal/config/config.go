package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	ENV_PREFFIX = "PF"
)

type Database struct {
	Host     string `default:"" required:"true"`
	Port     string `default:"5432" required:"true"`
	Name     string `required:"true"`
	Username string `required:"true"`
	Password string `required:"true"`
}

type Config struct {
	Host               string   `default:"0.0.0.0" required:"true" split_words:"true"`
	Port               string   `default:"4003" required:"true" split_words:"true"`
	Debug              bool     `default:"false"`
	Database           Database `required:"true" split_words:"true"`
	TransmissionApiUrl string   `required:"true" split_words:"true"`
	TorrentApiUrl      string   `required:"true" split_words:"true"`
	SubtitlerApiUrl    string   `required:"true" split_words:"true"`
}

func New() *Config {
	godotenv.Load()
	cfg, err := Get()
	if err != nil {
		panic(fmt.Errorf("configuration value(s) are not valid for environment: %w", err))
	}
	return cfg
}

func Get() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process(ENV_PREFFIX, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
