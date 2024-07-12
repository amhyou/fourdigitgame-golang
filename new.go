package main

import (
	"encoding/json"
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdGames := make([]Game, 1)
	err = surrealdb.Unmarshal(data, &createdGames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	createdGame := createdGames[0]

	createdGame.ID = extractRealID(createdGame.ID)
	// createdGame.Numbers = nil

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
