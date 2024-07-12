package main

import (
	"github.com/surrealdb/surrealdb.go"
)

func initDB() {
	var err error
	db, err = surrealdb.New("ws://localhost:8000/rpc")
	if err != nil {
		panic(err)
	}

	if _, err = db.Use("test", "test"); err != nil {
		panic(err)
	}
}
