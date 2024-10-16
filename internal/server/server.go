package server

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"

	"aibo/internal/database"
)

// Server represents the server instance.
//
// It contains the gin.Engine instance for handling HTTP requests
// and the database.Service instance for interacting with the database.
type Server struct {
	Router *gin.Engine
	DB     database.Service
}

// NewServer creates a new Server instance.
//
// It creates a new database service instance and stores it in the Server instance.
// If there is an error creating the database service instance, it returns a non-nil error.
// It also sets up the routes for the server.
//
// If there is an error setting up the routes, it returns a non-nil error.
func NewServer() (*Server, error) {
	dbservice, err := database.New()
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		return nil, fmt.Errorf("database connection failed: %v", err)
	}

	slog.Info("Database connected successfully")

	server := &Server{
		Router: gin.Default(),
		DB:     dbservice,
	}

	server.setupRoutes()

	return server, nil
}

// setupRoutes sets up the routes for the server.
//
// It creates an instance of AuthService and assigns it to handle the "/register" and "/login" routes.
//
// It creates a route group "/premium" that requires authentication and a premium subscription.
// The premium routes are not implemented yet.
//
// The "/profile" route is accessible only if the user is authenticated.
func (s *Server) setupRoutes() {
	SetupRoutes(s.Router, s.DB)
}

// Run starts the server and listens on the given address.
//
// It blocks until the server is stopped.
//
// If there is an error starting the server, it returns a non-nil error.
func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}
