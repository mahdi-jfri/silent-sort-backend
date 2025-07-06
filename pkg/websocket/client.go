package websocket

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"silent-sort/internal/logger"
	"silent-sort/pkg/hub"
	"time"
)

type Client struct {
	player *hub.Player
	conn   *websocket.Conn
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (c *Client) readPump() {
	defer func() {
		c.player.Hub.Messages <- &hub.MessageExit{Player: c.player}
		if err := c.conn.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close connection")
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error().Err(err).Msg("Error while reading message from ws")
			}
			break
		}
		type MessageJson struct {
			MessageType int            `json:"type"`
			Data        map[string]any `json:"data"`
		}
		messageJson := MessageJson{}
		if err := json.NewDecoder(bytes.NewReader(message)).Decode(&messageJson); err != nil {
			logger.Error().Err(err).Msg("Failed to decode json message")
			break
		}
		if messageJson.MessageType == 2 {
			c.player.Hub.Messages <- &hub.MessageStartGame{Player: c.player}
		} else if messageJson.MessageType == 3 {
			cardIdAny, exists := messageJson.Data["card_id"]
			if !exists {
				continue
			}
			cardId := cardIdAny.(string)
			logger.Info().Msg("Sending play card message")
			c.player.Hub.Messages <- &hub.MessagePlayCard{Player: c.player, CardId: cardId}
			logger.Info().Msg("Sent")
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.player.OutMessages:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
