package game

import (
	"github.com/TheAmirhosssein/duz-game-app/match"
)

var games map[string]*match.Match = make(map[string]*match.Match)

func JoinGame(gameInfo map[string]string) {
	game := games[gameInfo["gameId"]]
	if game != nil {
		game.SetSecondPlayer(gameInfo["userId"])
	} else {
		games[gameInfo["gameId"]] = match.New(gameInfo["userId"])
	}
}
