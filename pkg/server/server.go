package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"silent-sort/internal/config"
	"silent-sort/internal/logger"
	"time"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Config
}

func NewServer(cfg *config.Config) *Server {
	gin.SetMode(string(cfg.HTTPMode))
	ginEngine := gin.New()
	server := &Server{
		engine: ginEngine,
		cfg:    cfg,
	}
	return server
}

func (s *Server) addRoutes() {
	s.engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}

func (s *Server) Run(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.HTTPPort),
		Handler: s.engine,
	}

	serverDied := make(chan error)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverDied <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("Server shutdown failed")
		}
	case err := <-serverDied:
		return err
	}

	return nil
}
