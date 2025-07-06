package game

import (
	"fmt"
	"github.com/google/uuid"
)

type GameState int8

const (
	GameStateInLobby = 1
	GameStateStarted = 2
	GameStateWon     = 3
	GameStateLost    = 4
)

type SimpleSilentSortGame struct {
	gameState       GameState
	players         []string
	cardsById       map[string]Card
	playedCardsById map[string]bool
	playedCardIds   []string
	limit           int
}

func NewSimpleSilentSortGame(limit int) *SimpleSilentSortGame {
	return &SimpleSilentSortGame{
		gameState: GameStateInLobby,
		players:   nil,
		limit:     limit,
	}
}

func (s *SimpleSilentSortGame) CanAnyoneEnter() bool {
	return s.gameState == GameStateInLobby
}

func (s *SimpleSilentSortGame) HasBeenPlayed(cardId string) bool {
	return s.playedCardsById[cardId]
}

func (s *SimpleSilentSortGame) CanPlayCard(player string, cardId string) bool {
	if s.gameState != GameStateStarted {
		return false
	}
	card, exists := s.cardsById[cardId]
	if !exists {
		return false
	}
	if s.HasBeenPlayed(cardId) {
		return false
	}
	return card.Holder == player
}

func (s *SimpleSilentSortGame) PlayCard(player string, cardId string) {
	cardToPlay := s.cardsById[cardId]
	s.playedCardsById[cardToPlay.Id] = true
	s.playedCardIds = append(s.playedCardIds, cardToPlay.Id)
	if !s.HasRemainingCard() {
		s.gameState = GameStateWon
		return
	}
	for _, card := range s.cardsById {
		if s.HasBeenPlayed(card.Id) {
			continue
		}
		if card.Number < cardToPlay.Number {
			s.gameState = GameStateLost
			break
		}
	}
}

func (s *SimpleSilentSortGame) CanStartGame() bool {
	return s.gameState == GameStateInLobby
}

func (s *SimpleSilentSortGame) CanRestartGame() bool {
	return s.gameState == GameStateWon || s.gameState == GameStateLost
}

func (s *SimpleSilentSortGame) StartGame(players []string) {
	s.players = players
	s.cardsById = map[string]Card{}
	s.playedCardsById = map[string]bool{}
	s.playedCardIds = nil

	s.gameState = GameStateStarted
	cardNumbers := GenerateCardNumbers(len(s.players), s.limit)
	for i, cardNumber := range cardNumbers {
		card := Card{
			Number: cardNumber,
			Id:     uuid.NewString(),
			Holder: players[i],
		}
		s.cardsById[card.Id] = card
	}
}

func (s *SimpleSilentSortGame) RemovePlayer(player string) {
	if s.gameState != GameStateStarted {
		return
	}
	newPlayers := []string{}
	for _, inGamePlayer := range s.players {
		if inGamePlayer != player {
			newPlayers = append(newPlayers, inGamePlayer)
		}
	}
	s.players = newPlayers
	newCardsById := map[string]Card{}
	for _, card := range s.cardsById {
		if card.Holder != player || s.playedCardsById[card.Id] {
			newCardsById[card.Id] = card
		}
	}
	s.cardsById = newCardsById
	if !s.HasRemainingCard() {
		s.gameState = GameStateWon
	}
}

func (s *SimpleSilentSortGame) HasRemainingCard() bool {
	for _, card := range s.cardsById {
		if !s.playedCardsById[card.Id] {
			return true
		}
	}
	return false
}

func (s *SimpleSilentSortGame) GetPlayerCards(player string) []Card {
	cards := []Card{}
	for _, card := range s.cardsById {
		if card.Holder == player {
			cards = append(cards, card)
		}
	}
	return cards
}

func (s *SimpleSilentSortGame) GetPlayedCardsInOrder() []Card {
	cards := []Card{}
	for _, cardId := range s.playedCardIds {
		cards = append(cards, s.cardsById[cardId])
	}
	return cards
}

func (s *SimpleSilentSortGame) GetGameState() string {
	if s.gameState == GameStateInLobby {
		return "in_lobby"
	}
	if s.gameState == GameStateStarted {
		return "started"
	}
	if s.gameState == GameStateWon {
		return "won"
	}
	if s.gameState == GameStateLost {
		return "lost"
	}
	panic(fmt.Errorf("invalid game state: %v", s.gameState))
}

func (s *SimpleSilentSortGame) ShouldShowAllCards() bool {
	return s.gameState == GameStateLost
}

func (s *SimpleSilentSortGame) GetAllCards() []Card {
	cards := []Card{}
	for _, card := range s.cardsById {
		cards = append(cards, card)
	}
	return cards
}

func (s *SimpleSilentSortGame) RestartGame() {
	s.gameState = GameStateInLobby
	s.cardsById = map[string]Card{}
	s.playedCardsById = map[string]bool{}
	s.playedCardIds = []string{}
}
