package gapi

import (
	"fmt"

	db "example.com/simplebank/db/sqlc"
	"example.com/simplebank/pb"
	"example.com/simplebank/token"
	"example.com/simplebank/util"
	"example.com/simplebank/worker"
)

// this gRPC server will be used for banking service
type Server struct {
	pb.UnimplementedFinGoServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// Creates new gRPC instance
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}
	server := &Server{
		store: store, 
		tokenMaker: maker, 
		config: config,
		taskDistributor: taskDistributor,
	}

	return server, nil

}
