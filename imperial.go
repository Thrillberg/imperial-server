package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(*http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var output map[string]interface{}
var payload map[string][]string

var connections = map[int]*websocket.Conn{}
var players = map[int]string{}
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
		for _, value := range connections {
			players[myId] = string(p)

			var playersSlice []string
			for _, value := range players {
				playersSlice = append(playersSlice, value)
			}

			output = make(map[string]interface{})
			output["type"] = "registerPlayers"
			payload = make(map[string][]string)
			payload["players"] = playersSlice
			output["payload"] = payload

			out, _ := json.Marshal(output)
			if err := value.WriteMessage(messageType, out); err != nil {
				log.Println(err, myId)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handler)
	http.ListenAndServe(":8080", nil)
}
