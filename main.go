package main

import (
	"database/sql"
	"log"

	db "example.com/simplebank/db/sqlc"
	"example.com/simplebank/util"

	"example.com/simplebank/api"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	log.Printf("connected to db: %v", conn)

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	er := server.Start(config.ServerAddress)
	if er != nil {
		log.Fatal("cannot start server: ", er)
	}

}
