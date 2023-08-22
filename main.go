package main

import (
	"database/sql"
	"log"

	db "example.com/simplebank/db/sqlc"

	"example.com/simplebank/api"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simplebank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	log.Printf("connected to db: %v", conn)

	store := db.NewStore(conn)
	server := api.NewServer(store)

	er := server.Start(serverAddress)
	if er != nil {
		log.Fatal("cannot start server: ", er)
	}

}
