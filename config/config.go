package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
	"kiramishima/m-backend/internal/core/domain"
	"log"
)

func Load() (*domain.Configuration, error) {
	var cfg domain.Configuration
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// NewConfig creates and load config
func NewConfig() *domain.Configuration {
	cfg, err := Load()
	if err != nil {
		log.Printf("Can't load the configuration. Error: %s", err.Error())
	}

	return cfg
}

// Module config
var Module = fx.Options(
	fx.Provide(NewConfig),
)
