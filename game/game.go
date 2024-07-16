package game

import (
	"github.com/TheAmirhosssein/duz-game-app/client"
	"github.com/TheAmirhosssein/duz-game-app/match"
)

var games map[string]*match.Match = make(map[string]*match.Match)

func JoinGame(gameId string, client client.Client) {
	game := games[gameId]
	if game != nil {
		game.SetSecondPlayer(client)
	} else {
		games[gameId] = match.New(client)
	}
}
