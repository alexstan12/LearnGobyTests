package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "secret"
	dbname   = "postgres"
	sslmode  = "disable"
)

func main() {
	//store := NewInMemoryPlayerStore()
	//server := &PlayerServer{store}
	//log.Fatal(http.ListenAndServe(":5000", server))
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s port=%d", user, password, dbname, sslmode, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the db")
}
