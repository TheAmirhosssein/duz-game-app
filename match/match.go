package match

import (
	"github.com/TheAmirhosssein/duz-game-app/client"
)

type Match struct {
	OPlayer *client.Client
	XPlayer *client.Client
	Turn    string
	Moves   map[string]string
}

func New(player client.Client) *Match {
	match := Match{
		Turn: "X",
	}
	if playerIcon() == "X" {
		match.XPlayer = &player
	} else {
		match.OPlayer = &player
	}
	(&player).SendMessageToClient([]byte("waiting for your opponent to become ready"))
	return &match
}

func (match *Match) SetSecondPlayer(player client.Client) {
	if match.OPlayer == nil {
		match.OPlayer = &player
	} else {
		match.XPlayer = &player
	}
	message := "game started"
	match.OPlayer.SendMessageToClient([]byte(message))
	match.XPlayer.SendMessageToClient([]byte(message))
}
