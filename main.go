package main

import (
	"fmt"
	"log"
	"net/http"

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
			matchData, err := messages.ParseMessageForMatch(message)
			if err != nil {
				conn.WriteMessage(messageType, []byte(err.Error()))
			} else {
				fmt.Println(matchData)
			}
		}
	}
}

func main() {
	http.HandleFunc("/", echo)
	fmt.Println("listening on localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
