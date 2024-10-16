package database

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"aibo/internal/types"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	GetDB() *gorm.DB
	Migrate() error
}

// GetDB returns the underlying Gorm DB instance.
func (s *service) GetDB() *gorm.DB {
	return s.db
}

type service struct {
	db *gorm.DB
}

var (
	dbname     = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

// New returns a new database Service instance.
//
// The service instance is configured using the following environment variables:
//
// * DB_HOST: The hostname of the database server.
// * DB_PORT: The port number of the database server.
// * DB_DATABASE: The name of the database to connect to.
// * DB_USERNAME: The username to use when connecting to the database.
// * DB_PASSWORD: The password to use when connecting to the database.
//
// The service instance is reused if it has already been initialized. This
// means that calling New multiple times will return the same service instance
// each time.
//
// The service instance is configured with the following connection parameters:
//
// * MaxLifetime: 0 (connections can be kept open indefinitely)
// * MaxIdleConns: 50 (up to 50 idle connections are kept open)
// * MaxOpenConns: 50 (up to 50 open connections are allowed)
//
// If there is an error creating the service instance, it returns a non-nil error.
func New() (Service, error) {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance, nil
	}

	// Opening a driver typically will not attempt to connect to the database.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		return nil, err
	}

	slog.Info("Setting connection parameters", "maxLifetime", 0, "maxIdleConns", 50, "maxOpenConns", 50)
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	gormDb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	slog.Info("GORM", "GORM", gormDb)

	if err != nil {
		return nil, err
	}

	dbInstance = &service{
		db: gormDb,
	}
	return dbInstance, nil
}

// Health returns a map of health status information.
//
// The keys and values in the map are as follows:
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
func (s *service) Health() map[string]string {
	stats := make(map[string]string)

	sqlDB, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err))
		return stats
	}

	err = sqlDB.Ping()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err))
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}
	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}
	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close terminates the database connection.
// It returns an error if the connection cannot be closed.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dbname)

	sqldb, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqldb.Close()
}

// Migrate runs the database migrations. It is called automatically during the startup of the server.
// If there is an error migrating the database, it returns a non-nil error.
func (s *service) Migrate() error {
	// Auto-migrate the Aibo and CatBud models
	err := s.db.AutoMigrate(&types.Aibo{}, &types.CatBud{})
	if err != nil {
		return err
	}

	// Create index on CatBud's AiboID
	err = s.db.Exec("CREATE INDEX idx_catbuds_aibo_id ON cat_buds(aibo_id)").Error
	if err != nil {
		// If the error is because the index already exists, ignore it
		if !strings.Contains(err.Error(), "Duplicate key name") {
			return err
		}
	}

	// Create index on CatBud's Category
	err = s.db.Exec("CREATE INDEX idx_catbuds_category ON cat_buds(category)").Error
	if err != nil {
		// If the error is because the index already exists, ignore it
		if !strings.Contains(err.Error(), "Duplicate key name") {
			return err
		}
	}

	return nil
}
