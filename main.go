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
			var message string
			if messageTypeJson == "move" {
				if !match.EmptySquare(matchData["square"]) {
					conn.WriteMessage(messageType, []byte("square is not empty"))
					continue
				}
				if user.MaxMove() {
					user.SendMessageToClient([]byte("you can not move any pawn"))
					continue
				}
				match.Move(matchData["square"])
				message = fmt.Sprintf("%s selected %v square", turn, matchData["square"])
				user.MovedPawn()

			} else {
				if !user.MaxMove() {
					user.SendMessageToClient([]byte("you can not remove any pawn"))
					continue
				}
				if match.EmptySquare(matchData["square"]) {
					conn.WriteMessage(messageType, []byte("there is no pawn in this square"))
					continue
				}
				if !match.CheckValidRemove(matchData["square"]) {
					conn.WriteMessage(messageType, []byte("this square is for your opponent"))
					continue
				}
				match.RemovePawn(matchData["square"])
				user.RemovedPawn()
				message = fmt.Sprintf("%s removed %v square", turn, matchData["square"])
			}
			fmt.Println(match.IsGameOverColumn())
			match.XPlayer.SendMessageToClient([]byte(message))
			match.OPlayer.SendMessageToClient([]byte(message))
		}
	}
}

func main() {
	http.HandleFunc("/", echo)
	fmt.Println("listening on localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
