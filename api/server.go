package api

import (
	db "example.com/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// this server will be used to handle all the requests from the client
type Server struct {
	store  db.Store
	router *gin.Engine
}

// Creates new server instance
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.GET("hello", server.helloWorld)

	// account handlers
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccounts)
	router.POST("/transfers", server.CreateTransfer)
	router.POST("/register", server.CreateUser)
	router.POST("/login", server.LoginUser)
	router.GET("/users/:username", server.GetUser)

	server.router = router

	return server

}

func (server *Server) helloWorld(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Hello World",
	})
}

// Starts the server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
