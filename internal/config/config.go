package config

import (
	"log"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Host        string `env:"HOST" envDefault:"localhost"`
	HttpPort    string `env:"SERVER_PORT" envDefault:":8081"`
	AllowOrigin string `env:"ALLOW_ORIGIN" envDefault:"*"`
	Env         string `env:"ENV" envDefault:"dev"`

	GrpcURL       string        `env:"GRPC_URL" envDefault:"localhost:8090"`
	Secret        string        `env:"SECRET" envDefault:"SecretSecretSecret"`
	TokenLifespan time.Duration `env:"TOKEN_LIFESPAN" envDefault:"24h"`
}

func NewConfig() (Config, error) {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Print(".env file could not be imported")
	}

	err = env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
