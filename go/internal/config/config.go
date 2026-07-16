// Package config loads runtime settings from the environment.
//
// The Java service reads these from application.yml; the Go service is
// 12-factor and takes them purely from env vars, which is also what
// docker-compose and Coolify supply.
package config

import (
	"fmt"
	"os"

	"github.com/shopspring/decimal"
)

type Config struct {
	Port           string
	DatabaseURL    string
	InitialBalance decimal.Decimal
}

const (
	defaultPort           = "8081"
	defaultDatabaseURL    = "postgres://threedice:threedice@localhost:5432/threedice_go?sslmode=disable"
	defaultInitialBalance = "1000.00"
)

func Load() (Config, error) {
	balance, err := decimal.NewFromString(env("INITIAL_BALANCE", defaultInitialBalance))
	if err != nil {
		return Config{}, fmt.Errorf("INITIAL_BALANCE must be a decimal: %w", err)
	}

	return Config{
		Port:           env("PORT", defaultPort),
		DatabaseURL:    env("DATABASE_URL", defaultDatabaseURL),
		InitialBalance: balance,
	}, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
