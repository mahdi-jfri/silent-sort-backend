package websocket

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"silent-sort/internal/config"
	"silent-sort/internal/logger"
	"silent-sort/pkg/game"
	"silent-sort/pkg/hub"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	cfg       *config.Config
	hubKeeper *hub.HubKeeper
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebsocketServer(cfg *config.Config) *WebsocketServer {
	return &WebsocketServer{cfg: cfg, hubKeeper: hub.NewHubKeeper()}
}

func (s *WebsocketServer) handleConnections(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("room_id")
	if roomId == "" {
		http.Error(w, "room_id parameter is required", http.StatusBadRequest)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name parameter is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to upgrade connection")
		return
	}

	client := &Client{
		player: hub.NewPlayer(uuid.NewString(), name, conn),
		conn:   conn,
	}

	requestedHub := s.hubKeeper.GetHub(roomId)

	if requestedHub == nil {
		requestedHub = hub.NewHub(roomId, client.player, game.NewSimpleSilentSortGame(100))
		s.hubKeeper.SetHub(roomId, requestedHub)
		client.player.Hub = requestedHub
		go func() {
			requestedHub.Run(context.Background())
			s.hubKeeper.SetHub(roomId, nil)
		}()
	} else {
		client.player.Hub = requestedHub
		requestedHub.Messages <- &hub.MessageEnter{Player: client.player}
	}

	go client.writePump()
	go client.readPump()
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
