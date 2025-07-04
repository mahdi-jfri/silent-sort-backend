package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
	"os/signal"
	"silent-sort/internal/config"
	"silent-sort/internal/logger"
	"silent-sort/pkg/server"
	"silent-sort/pkg/websocket"
	"sync"
	"syscall"
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
	wsServer := websocket.NewWebsocketServer(cfg)
	srv := server.NewServer(cfg)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)

	childrenErrors := make(chan interface{}, 100)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err := srv.Run(ctx)
		wg.Done()
		if err != nil {
			childrenErrors <- err
		}
	}()

	wg.Add(1)
	go func() {
		err := wsServer.Run(ctx)
		wg.Done()
		if err != nil {
			childrenErrors <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info().Msg("Interrupt received")
	case <-childrenErrors:
		logger.Info().Msg("Some child process failed. Shutting down the rest...")
		cancel()
	}
	wg.Wait()
	logger.Info().Msg("All child processes done")

	return
}
