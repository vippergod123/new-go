package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/vipeergod123/simple_bank/db/sqlc"
	"github.com/vipeergod123/simple_bank/token"
	"github.com/vipeergod123/simple_bank/util"
)

type Server struct {
	store       db.Store
	tokenMarker token.Marker
	config      util.Config
	router      *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMarker, err := token.NewJWTMarker(config.TokenKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token: %w", err)
	}
	server := &Server{store: store, config: config, tokenMarker: tokenMarker}

	//Binding validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMarker))
	// accounts
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	// transfer
	authRoutes.POST("/transfer", server.createTransfer)
	// user
	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUser)
	// login
	authRoutes.POST("/login", server.loginUser)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
