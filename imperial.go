package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(*http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connections = map[int]*websocket.Conn{}
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
			if err := value.WriteMessage(messageType, p); err != nil {
				log.Println(err, myId)
				return
			}
			log.Println(messageType, myId, string(p))
		}
	}
}

func main() {
	http.HandleFunc("/ws", handler)
	http.ListenAndServe(":8080", nil)
}
