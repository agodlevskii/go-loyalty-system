package configs

import (
	"flag"
	"github.com/caarlos0/env"
)

type Config struct {
	Addr       string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DBURL      string `env:"DATABASE_URI"`
	AccrualURL string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := initFromEnv(cfg); err != nil {
		return nil, err
	}

	initFromFlags(cfg)
	return cfg, nil
}

func initFromEnv(cfg *Config) error {
	return env.Parse(cfg)
}

func initFromFlags(cfg *Config) {
	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "The application server address")
	flag.StringVar(&cfg.DBURL, "d", cfg.DBURL, "The database connection URI")
	flag.StringVar(&cfg.AccrualURL, "r", cfg.AccrualURL, "The accrual system URL")
	flag.Parse()
}
