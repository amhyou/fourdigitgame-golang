package main

import (
	"log"
	"net/http"

	"github.com/surrealdb/surrealdb.go"
)

func startGame(w http.ResponseWriter, r *http.Request) {

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
			}
			log.Println("guess: ", message)

			// Get game
			gameid := "games:" + message.Game
			data, _ := db.Select(gameid)
			game := new(Game)
			surrealdb.Unmarshal(data, &game)

			exact, misplaced := compareStrings(game.Numbers[message.Player], *message.Guess)
			message.Exact = &exact
			message.Misplaced = &misplaced

			for i := 0; i < 2; i++ {
				clients[message.Game][i].WriteJSON(&message)
			}

			message.Game = gameid
			db.Create("messages", message)

		case "join":
			wbs := clients[message.Game]
			length := len(wbs)
			if length > 1 {
				message.Action = "notfound"
				c.WriteJSON(&message)
				log.Println("game is already has 2 players")
				return
			}
			// Get game
			gameid := "games:" + message.Game
			data, _ := db.Select(gameid)
			game := new(Game)
			surrealdb.Unmarshal(data, &game)
			if game.Status != "NEW" {
				message.Action = "notfound"
				c.WriteJSON(&message)
				log.Println("game is already ended")
				return
			}

			message.Player = uint8(length)
			clients[message.Game] = append(clients[message.Game], c)
			log.Printf("player %v joined to game %v", message.Player, message.Game)
			clients[message.Game][length].WriteJSON(&message)
			if length+1 == 2 {
				startMessage := Message{Game: message.Game, Action: "start", Player: 2}
				for _, conn := range clients[message.Game] {
					conn.WriteJSON(&startMessage)
				}
			}

		// case "close":
		// 	log.Println(clients[message.Game])
		// 	for _, conn := range clients[message.Game] {
		// 		conn.Close()
		// 	}
		// 	clients[message.Game] = make([]*websocket.Conn, 2)
		// 	return

		case "stop":
			log.Println("game " + message.Game + " stopped")
			// Get game
			gameid := "games:" + message.Game
			data, _ := db.Select(gameid)
			game := new(Game)
			surrealdb.Unmarshal(data, &game)

			game.Status = "Ended"
			game.Reason = message.Guess
			game.Win = &message.Player

			db.Change(gameid, game)

			for i := 0; i < 2; i++ {
				clients[message.Game][i].WriteJSON(&message)
				clients[message.Game][i].Close()
			}
			delete(clients, message.Game)
			return
		}

	}
}

func compareStrings(str1, str2 string) (exact, displaced uint8) {
	exact = 0
	displaced = 0

	// Maps to keep track of unmatched digits
	unmatchedStr1 := make(map[rune]uint8)
	unmatchedStr2 := make(map[rune]uint8)

	// First pass: Count exact matches
	for i := 0; i < len(str1); i++ {
		if str1[i] == str2[i] {
			exact++
		} else {
			unmatchedStr1[rune(str1[i])]++
			unmatchedStr2[rune(str2[i])]++
		}
	}

	// Second pass: Count displaced matches
	for digit, count1 := range unmatchedStr1 {
		if count2, found := unmatchedStr2[digit]; found {
			displaced += min(count1, count2)
		}
	}
	return exact, displaced
}
