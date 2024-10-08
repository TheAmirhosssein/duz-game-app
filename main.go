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
				match, err := game.GetMatch(matchData["gameId"])
				if err != nil {
					conn.WriteMessage(messageType, []byte(err.Error()))
					break
				}
				jsonMessage := map[string]any{"user_sign": match.GetUserSign(client), "is_game_ready": match.IsGameReady()}
				userMessage := messages.GenerateMessage("join_game", matchData["userId"], matchData["gameId"], jsonMessage)
				client.SendMessageToClient(userMessage)
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
				info := map[string]any{"error": turnError.Error()}
				message := string(messages.GenerateMessage("error", matchData["userId"], matchData["gameId"], info))
				conn.WriteMessage(messageType, []byte(message))
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
					info := map[string]any{"error": "square is not empty"}
					message = string(messages.GenerateMessage("error", matchData["userId"], matchData["gameId"], info))
					conn.WriteMessage(messageType, []byte(message))
					continue
				}
				if user.MaxMove() {
					info := map[string]any{"error": "you can not move any pawn"}
					message = string(messages.GenerateMessage("error", matchData["userId"], matchData["gameId"], info))
					conn.WriteMessage(messageType, []byte(message))
					continue
				}
				match.Move(matchData["square"])
				info := map[string]any{"square": matchData["square"], "sign": turn}
				message = string(messages.GenerateMessage("move", matchData["userId"], matchData["gameId"], info))
				user.MovedPawn()
			} else {
				if !user.MaxMove() {
					info := map[string]any{"error": "you can not remove any pawn"}
					message = string(messages.GenerateMessage("error", matchData["userId"], matchData["gameId"], info))
					user.SendMessageToClient([]byte(message))
					continue
				}
				if match.EmptySquare(matchData["square"]) {
					info := map[string]any{"error": "there is no pawn in this square"}
					message = string(messages.GenerateMessage("error", matchData["userId"], matchData["gameId"], info))
					conn.WriteMessage(messageType, []byte(message))
					continue
				}
				if !match.CheckValidRemove(matchData["square"]) {
					info := map[string]any{"error": "this square is for your opponent"}
					message = string(messages.GenerateMessage("error", matchData["userId"], matchData["gameId"], info))
					conn.WriteMessage(messageType, []byte(message))
					continue
				}
				match.RemovePawn(matchData["square"])
				user.RemovedPawn()
				info := map[string]any{"square": matchData["square"], "sign": turn}
				message = string(messages.GenerateMessage("remove", matchData["userId"], matchData["gameId"], info))
			}
			match.XPlayer.SendMessageToClient([]byte(message))
			match.OPlayer.SendMessageToClient([]byte(message))
			if match.IsGameOverColumn() || match.IsGameOverRow() || match.IsGameOverDiagonal() {
				info := map[string]any{"winner": turn}
				message = string(messages.GenerateMessage("game_over", matchData["userId"], matchData["gameId"], info))
				match.XPlayer.SendMessageToClient([]byte(message))
				match.OPlayer.SendMessageToClient([]byte(message))
				break
			}
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func main() {
	http.HandleFunc("/ws/", echo)
	fmt.Println("listening on localhost:8080")
	staticDir := "./static"
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveHome)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
