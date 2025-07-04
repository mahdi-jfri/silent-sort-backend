package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"silent-sort/internal/config"
	"silent-sort/internal/logger"
)

func mustGetConfig() *config.Config {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("Failed to load .env file: %w", err))
	}
	cfg := &config.Config{}
	err = envconfig.Process("", cfg)
	if err != nil {
		panic(fmt.Errorf("Failed to load config: %w", err))
	}
	return cfg
}

func main() {
	cfg := mustGetConfig()
	logger.Init(cfg)
	logger.Info().Interface("cfg", cfg).Msg("LOG?")
}
