package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/surrealdb/surrealdb.go"
)

func newGame(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	numbers := [2]string{generateUnique4DigitNumber(), generateUnique4DigitNumber()}
	// New game
	game := Game{
		Status: "NEW", Numbers: &numbers,
	}

	// Create new game
	data, err := db.Create("games", game)
	log.Println(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("game created")

	createdGame := new(Game)
	err = surrealdb.Unmarshal(data, &createdGame)
	log.Println(createdGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("game unmarsheled")

	createdGame.ID = extractRealID(createdGame.ID)
	log.Println("Number for game " + createdGame.ID + " are: " + createdGame.Numbers[0] + ", " + createdGame.Numbers[1])
	createdGame.Numbers = nil

	jsonGame, err := json.Marshal(createdGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonGame)
}

func extractRealID(id string) string {
	// Split the string at the colon
	parts := strings.Split(id, ":")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func generateUnique4DigitNumber() string {
	digits := []rune("0123456789")
	rand.Shuffle(len(digits), func(i, j int) { digits[i], digits[j] = digits[j], digits[i] })
	uniqueDigits := digits[:4]
	return string(uniqueDigits)
}
