package websocket

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"silent-sort/internal/config"
	"silent-sort/internal/logger"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	cfg *config.Config
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebsocketServer(cfg *config.Config) *WebsocketServer {
	return &WebsocketServer{cfg: cfg}
}

func (s *WebsocketServer) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to upgrade connection")
		return
	}
	defer func(ws *websocket.Conn) {
		if err := ws.Close(); err != nil {
			logger.Error().Err(err).Msg("Websocket close failed")
		}
	}(ws)

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				logger.Error().Err(err).Msg("Error reading message")
			}
			break
		}

		if err := ws.WriteMessage(messageType, p); err != nil {
			logger.Error().Err(err).Msg("Error writing message")
			break
		}
	}
}

func (s *WebsocketServer) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleConnections)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.WebsocketPort),
		Handler: mux,
	}

	serverDied := make(chan error)

	go func() {
		logger.Info().Msg("Starting websocket server")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverDied <- err
		}
	}()

	<-ctx.Done()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("Server shutdown failed")
		}
	case err := <-serverDied:
		return err
	}

	return nil
}
