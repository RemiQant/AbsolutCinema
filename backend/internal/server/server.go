package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"

	"absolutcinema-backend/internal/database"
	"absolutcinema-backend/internal/utils"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	// Load environment variables based on APP_ENV
	envFile := utils.GetEnvFile()
	if err := utils.LoadEnvFile(envFile); err != nil {
		panic(fmt.Sprintf("Failed to load env file %s: %v", envFile, err))
	}
	// Load PORT from environment variables (dev or prod)
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080" // Default port if not set
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Sprintf("Invalid PORT value: %v", err))
	}

	// Initialize database connection
	db := database.New()

	// Run database migrations
	if err := db.Migrate(); err != nil {
		panic(fmt.Sprintf("Failed to run migrations: %v", err))
	}

	NewServer := &Server{
		port: port,
		db:   db,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

// GetDB returns the GORM database instance
func (s *Server) GetDB() *gorm.DB {
	return s.db.DB()
}
