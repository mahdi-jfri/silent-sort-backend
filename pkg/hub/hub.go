package hub

import (
	"context"
	"fmt"
	"silent-sort/internal/logger"
	"silent-sort/pkg/game"
)

type Hub struct {
	id       string
	Messages chan interface{}
	owner    *Player
	players  []*Player
	game     game.SilentSortGame
}

func NewHub(id string, owner *Player, game game.SilentSortGame) *Hub {
	return &Hub{
		id:       id,
		game:     game,
		owner:    owner,
		Messages: make(chan interface{}),
		players:  []*Player{owner},
	}
}

func (h *Hub) Run(ctx context.Context) {
	for _, player := range h.players {
		data := h.getGameDataForPlayer(player)
		player.OutMessages <- data
	}
	for {
		select {
		case <-ctx.Done():
			break
		case message := <-h.Messages:
			switch message.(type) {
			case *MessageEnter:
				h.enter(message.(*MessageEnter))
			case *MessageExit:
				h.exit(message.(*MessageExit))
			case *MessageStartGame:
				h.startGame(message.(*MessageStartGame))
			case *MessagePlayCard:
				h.playCard(message.(*MessagePlayCard))
			default:
				panic(fmt.Errorf("invalid message: %+v", message))
			}
		}
		for _, player := range h.players {
			data := h.getGameDataForPlayer(player)
			player.OutMessages <- data
		}
		if len(h.players) == 0 {
			return
		}
	}
}

func (h *Hub) hasPlayer(player *Player) bool {
	for _, playerInside := range h.players {
		if playerInside == player {
			return true
		}
	}
	return false
}

func (h *Hub) enter(message *MessageEnter) {
	if h.hasPlayer(message.Player) {
		// handle already inside
		return
	}
	if !h.game.CanAnyoneEnter() {
		// handle cannot enter
		return
	}
	h.players = append(h.players, message.Player)
	if h.owner == nil {
		h.owner = message.Player
	}
	return
}

func (h *Hub) exit(message *MessageExit) {
	if !h.hasPlayer(message.Player) {
		// handle no such Player
		return
	}
	newPlayers := []*Player{}
	for _, player := range h.players {
		if player != message.Player {
			newPlayers = append(newPlayers, player)
		}
	}
	h.players = newPlayers
	h.game.RemovePlayer(message.Player.id)
	if h.owner == message.Player {
		if len(h.players) > 0 {
			h.owner = h.players[0]
		} else {
			h.owner = nil
		}
	}
}

func (h *Hub) getPlayerIds() []string {
	playerIds := []string{}
	for _, player := range h.players {
		playerIds = append(playerIds, player.id)
	}
	return playerIds
}

func (h *Hub) startGame(message *MessageStartGame) {
	logger.Info().Msgf("Starting game, owner: %p, player: %p", h.owner, message.Player)
	if h.owner != message.Player {
		return
	}
	h.game.StartGame(h.getPlayerIds())
}

func (h *Hub) playCard(message *MessagePlayCard) {
	if !h.hasPlayer(message.Player) {
		// handle no such Player
		return
	}
	if !h.game.CanPlayCard(message.Player.id, message.CardId) {
		// handle cannot play
		return
	}
	h.game.PlayCard(message.Player.id, message.CardId)
}

func (h *Hub) getGameDataForPlayer(player *Player) map[string]any {
	gameData := h.getGeneralGameData()
	playerData := h.getPlayerData(player)

	return map[string]any{
		"game_data":     gameData,
		"personal_data": playerData,
	}
}

func (h *Hub) getGeneralGameData() interface{} {
	allPlayers := []interface{}{}
	for _, player := range h.players {
		allPlayers = append(allPlayers, map[string]any{
			"id":   player.id,
			"name": player.name,
		})
	}
	playedCards := []interface{}{}
	for _, card := range h.game.GetPlayedCardsInOrder() {
		playedCards = append(playedCards, map[string]any{
			"id":     card.Id,
			"number": card.Number,
			"holder": card.Holder,
		})
	}

	result := map[string]any{
		"state":        h.game.GetGameState(),
		"players":      allPlayers,
		"played_cards": playedCards,
	}

	if h.game.ShouldShowAllCards() {
		allCards := []interface{}{}
		for _, card := range h.game.GetAllCards() {
			allCards = append(allCards, map[string]any{
				"id":     card.Id,
				"number": card.Number,
				"holder": card.Holder,
			})
		}
		result["all_cards"] = allCards
	}
	return result
}

func (h *Hub) getPlayerData(player *Player) interface{} {
	playerCards := []interface{}{}
	for _, card := range h.game.GetPlayerCards(player.id) {
		playerCards = append(playerCards, map[string]any{
			"id":     card.Id,
			"number": card.Number,
		})
	}

	return map[string]any{
		"id":       player.id,
		"cards":    playerCards,
		"is_owner": h.owner == player,
	}
}
