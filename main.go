package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TheAmirhosssein/duz-game-app/client"
	"github.com/TheAmirhosssein/duz-game-app/game"
	"github.com/TheAmirhosssein/duz-game-app/messages"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read:", err)
			break
		}
		messageTypeJson, err := messages.GetMessageType(&message)
		if err != nil {
			err = conn.WriteMessage(messageType, []byte(err.Error()))
			if err != nil {
				log.Println("Write:", err)
				break
			}
		}
		if messageTypeJson == "join_match" {
			validKeys := []string{"game_id", "user_id"}
			matchData, err := messages.ParseMessage(validKeys, message)
			if err != nil {
				conn.WriteMessage(messageType, []byte(err.Error()))
			} else {
				id := r.URL.Query().Get("id")
				client := client.New(id, matchData["userId"], conn)
				game.RegisterUser(*client)
				game.JoinGame(matchData["gameId"], *client)
			}
		}
		if messageTypeJson == "move" || messageTypeJson == "remove" {
			validKeys := []string{"game_id", "user_id", "square"}
			matchData, err := messages.ParseMessage(validKeys, message)
			if err != nil {
				conn.WriteMessage(messageType, []byte(err.Error()))
				continue
			}
			match, err := game.GetMatch(matchData["gameId"])
			turn := match.Turn
			if err != nil {
				conn.WriteMessage(messageType, []byte(err.Error()))
				break
			}
			user, err := game.GetUser(matchData["userId"])
			if err != nil {
				conn.WriteMessage(messageType, []byte(err.Error()))
				break
			}
			turnError := match.CheckUserTurn(*user)
			if turnError != nil {
				conn.WriteMessage(messageType, []byte(turnError.Error()))
				continue
			}
			squareNumberError := match.CheckValidSquareNumber(matchData["square"])
			if squareNumberError != nil {
				conn.WriteMessage(messageType, []byte(squareNumberError.Error()))
				continue
			}
			squareNumber := match.EmptySquare(matchData["square"])
			if squareNumber != nil {
				conn.WriteMessage(messageType, []byte(squareNumber.Error()))
				continue
			}
			if messageTypeJson == "move" {
				if user.MaxMove() {
					user.SendMessageToClient([]byte("you can not move any pawn"))
					continue
				}
				err = match.Move(*user, matchData["square"])
				if err != nil {
					conn.WriteMessage(messageType, []byte(err.Error()))
				} else {
					message := fmt.Sprintf("%s selected %v square", turn, matchData["square"])
					match.XPlayer.SendMessageToClient([]byte(message))
					match.OPlayer.SendMessageToClient([]byte(message))
					user.MovedPawn()
				}

			} else {
				if !user.MaxMove() {
					user.SendMessageToClient([]byte("you can not remove any pawn"))
					continue
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/", echo)
	fmt.Println("listening on localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
