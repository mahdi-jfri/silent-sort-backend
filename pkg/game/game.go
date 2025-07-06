package game

import (
	"math/rand/v2"
)

type SilentSortGame interface {
	CanAnyoneEnter() bool
	CanPlayCard(player string, cardId string) bool
	PlayCard(player string, cardId string)
	CanStartGame() bool
	CanRestartGame() bool
	StartGame(players []string)
	RemovePlayer(player string)
	GetPlayerCards(player string) []Card
	GetPlayedCardsInOrder() []Card
	GetGameState() string
	ShouldShowAllCards() bool
	GetAllCards() []Card
	RestartGame()
}

type Card struct {
	Number int
	Id     string
	Holder string
}

func GenerateCardNumbers(num int, limit int) []int {
	usedNumber := map[int]bool{}
	cardNumbers := []int{}
	for i := 0; i < num; i++ {
		number := 0
		for number == 0 || usedNumber[number] {
			number = rand.N(limit) + 1
		}
		cardNumbers = append(cardNumbers, number)
		usedNumber[number] = true
	}
	return cardNumbers
}
