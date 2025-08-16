package config

import (
	"party-games/internal/utils"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host  string `env:"HOST"`
	Port  string `env:"PORT" envDefault:"3000"`
	DATA  string `env:"DATA" envDefault:"data.db"`
	Debug bool   `env:"DEBUG"`

	OAI OpenAI `envPrefix:"OAI_"`
}

type OpenAI struct {
	BaseURL string `env:"BASE_URL"`
	APIKey  string `env:"API_KEY"`
	Model   string `env:"MODEL"`
}

func MustNew() *Config {
	var cfg Config
	utils.Must(env.Parse(&cfg))
	return &cfg
}
