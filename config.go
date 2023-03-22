package main

import (
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
)

const envPrefix = "DEMO_"

type Config struct {
	Port            int           `env:"PORT,required"             validate:"gte=0,lte=65536"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,required" validate:"gte=0"`
	Logger          string        `env:"LOGGER,required"           validate:"oneof=dev prod"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg, env.Options{
		Prefix: envPrefix,
	})
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
