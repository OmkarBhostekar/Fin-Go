package api

import (
	"fmt"

	db "example.com/simplebank/db/sqlc"
	"example.com/simplebank/token"
	"example.com/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// this server will be used to handle all the requests from the client
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
	router     *gin.Engine
}

// Creates new server instance
func NewServer(config util.Config, store db.Store) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}
	server := &Server{store: store, tokenMaker: maker, config: config}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil

}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.GET("hello", server.helloWorld)
	router.POST("/auth/register", server.CreateUser)
	router.POST("/auth/login", server.LoginUser)
	router.POST("/auth/token/renew", server.RenewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// account handlers
	authRoutes.POST("/accounts", server.CreateAccount)
	authRoutes.GET("/accounts/:id", server.GetAccount)
	authRoutes.GET("/accounts", server.ListAccounts)
	authRoutes.POST("/transfers", server.CreateTransfer)
	authRoutes.GET("/users/:username", server.GetUser)
	server.router = router
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
