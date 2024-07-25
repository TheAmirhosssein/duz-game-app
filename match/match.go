package match

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/TheAmirhosssein/duz-game-app/client"
)

type Match struct {
	OPlayer *client.Client
	XPlayer *client.Client
	Turn    string
	Moves   map[string]string
}

func New(player client.Client) *Match {
	moves := make(map[string]string)
	for counter := 1; counter <= 9; counter++ {
		moves[fmt.Sprint(counter)] = ""
	}
	match := Match{
		Turn:  "X",
		Moves: moves,
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
	message := "game started and your sign is "
	match.OPlayer.SendMessageToClient([]byte(message + "O"))
	match.XPlayer.SendMessageToClient([]byte(message + "X"))
}

func (match *Match) Move(square string) {
	match.Moves[square] = match.Turn
	match.changeTurn()
	match.showMoves()
}

func (match *Match) RemovePawn(square string) {
	match.Moves[square] = ""
	match.showMoves()
}

func (match *Match) getTurnUser() client.Client {
	if match.Turn == "X" {
		return *match.XPlayer
	} else {
		return *match.OPlayer
	}
}

func (match *Match) changeTurn() {
	if match.Turn == "X" {
		match.Turn = "O"
	} else {
		match.Turn = "X"
	}
}

func (match *Match) showMoves() {
	for counter := 1; counter < 10; counter++ {
		fmt.Print(match.Moves[fmt.Sprint(counter)])
		if match.Moves[fmt.Sprint(counter)] == "" {
			fmt.Print("-")
		}
		if counter%3 == 0 {
			fmt.Print("\n")
		}
	}
}

func (match *Match) CheckUserTurn(player client.Client) error {
	turnUser := match.getTurnUser()
	if turnUser.UserId != player.UserId {
		return errors.New("it is not your turn")
	}
	return nil
}

func (match *Match) CheckValidSquareNumber(square string) error {
	squareNumber, err := strconv.ParseInt(square, 10, 0)
	if err != nil || 1 > squareNumber || squareNumber > 9 {
		return errors.New("invalid square")
	}
	return nil
}

func (match *Match) EmptySquare(square string) bool {
	return match.Moves[square] == ""
}

func (match *Match) CheckValidRemove(square string) bool {
	return match.Moves[square] == match.Turn
}

func (match Match) IsGameOverColumn() bool {
	for counter := 1; counter <= 3; counter++ {
		lastSign := match.Moves[fmt.Sprint(counter)]
		if lastSign != "" {
			gameOver := true
			for squareIterator := counter; squareIterator <= counter+6; squareIterator += 3 {
				if match.Moves[fmt.Sprint(squareIterator)] != lastSign {
					gameOver = false
					break
				}
			}
			if gameOver {
				return true
			}
		}
	}
	return false
}

func (match Match) IsGameOverRow() bool {
	squares := [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}}
	for _, row := range squares {
		gameOver := true
		firstSquareValue := match.Moves[row[0]]
		if firstSquareValue != "" {

			for _, square := range row {
				if match.Moves[square] != firstSquareValue {
					gameOver = false
					break
				}
			}
			if gameOver {
				return true
			}
		}
	}
	return false
}
