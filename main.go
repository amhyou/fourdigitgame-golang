package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/surrealdb/surrealdb.go"
)

var db *surrealdb.DB

var clients = make(map[string][]*websocket.Conn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

type Game struct {
	ID      string     `json:"id,omitempty"`
	Status  string     `json:"status"`
	Numbers *[2]string `json:"nb"`
	Win     *uint8     `json:"win"`
	Reason  *string    `json:"reason"`
}

type Message struct {
	ID        string  `json:"id,omitempty"`
	Game      string  `json:"game"`
	Action    string  `json:"action"`
	Player    uint8   `json:"player"`
	Guess     *string `json:"guess,omitempty"`
	Exact     *uint8  `json:"exact,omitempty"`
	Misplaced *uint8  `json:"misplaced,omitempty"`
}

func main() {
	initDB()
	http.HandleFunc("/start", startGame)
	http.HandleFunc("/new", newGame)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
