package main

import (
	"aibo/internal/server"
	"fmt"
	"log/slog"
	"os"

	_ "aibo/docs"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// main is the main entry point for the application.
//
// It creates a new server instance, then calls its ListenAndServe method. If
// the server fails to start, it panics with an error message.

// @title           Aibo API
// @version         1.0
// @description     This is a server for managing Aibo's backend.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host           ${APP_URL}
// @BasePath       /api/v1
func main() {

	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			slog.Error("Error loading .env file", "err", err)
		}
	}

	slog.Info("Starting server...",
		"ENV", os.Getenv("ENV"),
		"PORT", os.Getenv("PORT"),
		"DB_HOST", os.Getenv("DB_HOST"),
		"DB_PORT", os.Getenv("DB_PORT"),
		"DB_DATABASE", os.Getenv("DB_DATABASE"),
		"DB_USERNAME", os.Getenv("DB_USERNAME"),
		"APP_URL", os.Getenv("APP_URL"))

	s, err := server.NewServer()
	if err != nil {
		slog.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")

	addr := ":" + port

	// Set up Swagger
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "localhost" + addr
	}
	url := ginSwagger.URL(appURL + "/swagger/doc.json") // The url pointing to API definition
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	slog.Info("Server created", "address", addr)
	slog.Info("Swagger UI available at", "url", "http://"+appURL+"/swagger/index.html")

	slog.Info("Server starting", "address", addr)

	err = s.Run(addr)
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
