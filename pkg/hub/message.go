package hub

type MessageEnter struct {
	Player *Player
}

type MessageExit struct {
	Player *Player
}

type MessageStartGame struct {
	Player *Player
}

type MessagePlayCard struct {
	Player *Player
	CardId string
}

type MessageRestartGame struct {
	Player *Player
}
