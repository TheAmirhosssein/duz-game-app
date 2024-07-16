package match

import (
	"math/rand"

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
	return &match
}

func (match *Match) SetSecondPlayer(player client.Client) {
	if match.OPlayer == nil {
		match.OPlayer = &player
	} else {
		match.XPlayer = &player
	}
}

func playerIcon() string {
	randNumber := rand.Int()
	if randNumber%2 == 0 {
		return "X"
	} else {
		return "O"
	}
}
