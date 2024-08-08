package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const ASSETS_PATH = "assets"

func GetDB() *sql.DB {
	db, err := sql.Open("sqlite3", ASSETS_PATH+"/lyrics.db")

	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	log.Printf("Connected to database: %s", ASSETS_PATH+"/lyrics.db")

	return db
}
