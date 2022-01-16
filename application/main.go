package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	//store := NewInMemoryPlayerStore()
	//server := &PlayerServer{store}
	//log.Fatal(http.ListenAndServe(":5000", server))
	db, err := sql.Open("pq","postgres:secret@tcp(127.0.0.1:5432)")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}
