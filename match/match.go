package match

import "math/rand"

type Match struct {
	OPlayer string
	XPlayer string
	Turn    string
	Moves   map[string]string
}

func New(gameId, firstPlayer string) *Match {
	match := Match{
		Turn: "X",
	}
	if playerIcon() == "X" {
		match.XPlayer = firstPlayer
	} else {
		match.OPlayer = firstPlayer
	}
	return &match
}

func playerIcon() string {
	randNumber := rand.Int()
	if randNumber%2 == 0 {
		return "X"
	} else {
		return "O"
	}
}
