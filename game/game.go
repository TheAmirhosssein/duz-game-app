package game

import (
	"errors"

	"github.com/TheAmirhosssein/duz-game-app/client"
	"github.com/TheAmirhosssein/duz-game-app/match"
)

var games map[string]*match.Match = make(map[string]*match.Match)
var users map[string]*client.Client = make(map[string]*client.Client)

func RegisterUser(user client.Client) {
	users[user.UserId] = &user
}

func GetUser(userId string) (*client.Client, error) {
	user := users[userId]
	if user == nil {
		return nil, errors.New("user not found")
	} else {
		return user, nil
	}
}

func JoinGame(gameId string, client client.Client) {
	game := games[gameId]
	if game != nil {
		game.SetSecondPlayer(client)
	} else {
		games[gameId] = match.New(client)
	}
}

func GetMatch(gameId string) (*match.Match, error) {
	game := games[gameId]
	if game == nil {
		return nil, errors.New("match not found")
	} else {
		return game, nil
	}
}
