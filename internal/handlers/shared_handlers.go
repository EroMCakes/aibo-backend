package handlers

import (
	"aibo/internal/database"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DBHealthHandler is a gin.HandlerFunc that returns the health status of the
// database.
//
// It takes a database.Service instance as an argument and returns a gin.HandlerFunc.
// The returned handler queries the database for its health status and returns it as
// a JSON response with a 200 status code.
//
// The health status is a map of strings with the following keys and values:
//
// * "status": The overall status of the database connection, either "up" or "down".
// * "message": A human-readable message describing the health status.
// * "open_connections": The number of open connections to the database.
// * "in_use": The number of connections currently in use.
// * "idle": The number of idle connections.
// * "wait_count": The number of times a connection was requested and the pool was empty.
// * "wait_duration": The total time waited for a connection.
// * "max_idle_closed": The number of connections closed due to idle timeout.
// * "max_lifetime_closed": The number of connections closed due to max lifetime.
//
// If the database connection is down, the "error" key will be present with
// the error message as its value.
//
// If there is an error retrieving the health status, it returns a non-nil error.
// @Summary Get database health
// @Description Get the health status of the database
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func DBHealthHandler(dbService database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		healthStatus := dbService.Health()
		c.JSON(http.StatusOK, healthStatus)
	}
}

// HandleMigrate is a gin.HandlerFunc that runs the database migrations and
// returns a JSON response with a 200 status code if the migrations are successful.
//
// It takes a gin.Context as an argument and queries the database for its health
// status. If the database connection is down, it returns a 503 error with a JSON
// response containing the error message.
//
// If there is an error retrieving the health status, it returns a non-nil error.
// @Summary Run database migrations
// @Description Run the database migrations
// @Tags database
// @Produce plain
// @Success 200 {string} string "Database migrated"
// @Failure 500 {string} string "Error message"
// @Router /migrate [post]
func HandleMigrate(dbService database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := dbService.Migrate()
		if err != nil {
			slog.Error("Failed to migrate database", "error", err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.String(http.StatusOK, "Database migrated")
	}
}
