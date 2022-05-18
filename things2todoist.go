package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile string = "main.sqlite"

func main() {

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	thingsDb := ThingsDB{db}
	tags, err := thingsDb.getTags()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tags)
}
