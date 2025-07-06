package hub

import "github.com/gorilla/websocket"

type Player struct {
	id          string
	name        string
	conn        *websocket.Conn
	OutMessages chan interface{}
	Hub         *Hub
}

func NewPlayer(id string, name string, conn *websocket.Conn) *Player {
	return &Player{
		id:          id,
		name:        name,
		conn:        conn,
		OutMessages: make(chan interface{}),
	}
}
