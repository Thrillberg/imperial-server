package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var gameLog []Action

type Action struct {
	Type    string
	Payload map[string]interface{}
}

func main() {
	http.HandleFunc("/tick", startHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	if r.Method == "OPTIONS" {
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var action Action
	err := json.Unmarshal(body, &action)
	if err != nil {
		fmt.Println("error:", err)
	}
	gameLog = append(gameLog, action)
	fmt.Println(gameLog)

	out, _ := json.Marshal(gameLog)
	io.WriteString(w, string(out))
}
