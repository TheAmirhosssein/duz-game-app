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
		if messageTypeJson == "move" {
			validKeys := []string{"game_id", "user_id", "square"}
			matchData, err := messages.ParseMessage(validKeys, message)
			if err != nil {
				conn.WriteMessage(messageType, []byte(err.Error()))
			} else {
				match, err := game.GetMatch(matchData["gameId"])
				if err != nil {
					conn.WriteMessage(messageType, []byte(err.Error()))
				}
				user, err := game.GetUser(matchData["userId"])
				if err != nil {
					conn.WriteMessage(messageType, []byte(err.Error()))
					break
				}
				err = match.Move(*user, matchData["square"])
				if err != nil {
					conn.WriteMessage(messageType, []byte(err.Error()))
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
