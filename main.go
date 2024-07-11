package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients map[string][]*websocket.Conn

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
} // use default options

type Message struct {
	Game    string `json:"game"`
	Player  uint8  `json:"player"`
	Action  string `json:"action"`
	Content string `json:"content"`
}

func game(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		var message Message
		err := c.ReadJSON(&message)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Println(message)
		switch message.Action {

		case "guess":
			if length := len(clients[message.Game]); length < 2 {
				log.Println("game does not contain 2 players yet")
				return
			} else {
				log.Println("guess: ", message)
				for i := 0; i < length; i++ {
					clients[message.Game][i].WriteJSON(&message)
				}
			}

		case "join":
			if length := len(clients[message.Game]); length > 1 {
				log.Println("game is already has 2 players")
				return
			} else {
				message.Player = uint8(length)
				clients[message.Game] = append(clients[message.Game], c)
				log.Printf("player %v joined to game %v", message.Player, message.Game)
				clients[message.Game][length].WriteJSON(&message)
				if length+1 == 2 {
					startMessage := Message{Game: message.Game, Action: "start", Player: 2, Content: ""}
					for _, conn := range clients[message.Game] {
						conn.WriteJSON(&startMessage)
					}
				}
			}

		case "close":
			log.Println(clients[message.Game])
			for _, conn := range clients[message.Game] {
				conn.Close()
			}
			clients[message.Game] = make([]*websocket.Conn, 2)
			return
		}

	}
}

type Response struct {
	Message string `json:"message"`
}

func home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Create a response struct
	response := Response{
		Message: "Hello, World!",
	}

	// Marshal response struct into JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {

	clients = make(map[string][]*websocket.Conn)
	http.HandleFunc("/ws", game)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
