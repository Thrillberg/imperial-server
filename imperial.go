package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type RegisterPlayerData struct {
	Type    string `json:"type"`
	Payload struct {
		Name string `json:"name"`
	} `json:"payload"`
}

type RegisteredPlayersOutput struct {
	Type    string `json:"type"`
	Payload struct {
		Players []string `json:"players"`
	} `json:"payload"`
}

type ActionData struct {
	Type    string `json:"type"`
	Payload struct {
		Action string `json:"action"`
	} `json:"payload"`
}

type GameLogOutput struct {
	Type    string `json:"type"`
	Payload struct {
		GameLog []string `json:"gameLog"`
	} `json:"payload"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(*http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var output map[string]interface{}

var payload map[string][]string

var connections = map[int]*websocket.Conn{}
var players = map[int]string{}
var gameLog = []string{}
var nextId = 1

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	myId := nextId
	nextId++
	connections[myId] = conn

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err, myId)
			return
		}

		log.Println("incoming message: ", string(p))

		var output interface{}
		var rawMessage map[string]interface{}
		json.Unmarshal(p, &rawMessage)

		if rawMessage["type"] == "action" {
			var message ActionData
			json.Unmarshal(p, &message)
			log.Println("action full: ", message.Payload.Action)
			if message.Payload.Action != "" {
				gameLog = append(gameLog, message.Payload.Action)
				log.Println("after adding action: ", gameLog)
			}

			output = GameLogOutput{
				Type: "gameLog",
				Payload: struct {
					GameLog []string `json:"gameLog"`
				}{GameLog: gameLog},
			}
		} else {
			var message RegisterPlayerData
			json.Unmarshal(p, &message)

			if rawMessage["type"] == "registerPlayer" {
				players[myId] = message.Payload.Name
			} else if rawMessage["type"] == "unregisterPlayer" {
				for unregisterKey, _ := range players {
					if players[unregisterKey] == message.Payload.Name {
						delete(players, unregisterKey)
					}
				}
			}

			playersSlice := []string{}
			for _, value := range players {
				log.Println(value)
				playersSlice = append(playersSlice, value)
			}

			output = RegisteredPlayersOutput{
				Type: "registeredPlayers",
				Payload: struct {
					Players []string `json:"players"`
				}{Players: playersSlice},
			}
		}

		out, _ := json.Marshal(output)

		for _, connValue := range connections {
			if err := connValue.WriteMessage(messageType, out); err != nil {
				log.Println(err, myId)
				return
			}
			log.Println("outgoing message:", string(out))
		}
	}
}

func main() {
	http.HandleFunc("/ws", handler)
	http.ListenAndServe(":8080", nil)
}
