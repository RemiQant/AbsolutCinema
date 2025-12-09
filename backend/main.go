package main

import (
	"fmt"
	"net/http"
	"os"   // New import for os package
	"time" // New import for time package

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️  No .env file found, using system environment variables")
	}

	// 1. Initialize Sentry with DSN from Environment
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		fmt.Println("⚠️  SENTRY_DSN not set!")
	} else {
		fmt.Printf("✅ Sentry DSN loaded: %s...\n", dsn[:40])
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
		// Recommended for performance monitoring
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	} else {
		fmt.Println("✅ Sentry initialized successfully!")
	}
	// IMPORTANT: Ensure Sentry events are sent before the program terminates
	defer sentry.Flush(2 * time.Second)

	// Then create your app
	app := gin.Default()

	// Once it's done, you can attach the handler as one of your middleware
	app.Use(sentrygin.New(sentrygin.Options{}))

	// Database connection string
	dbURL := os.Getenv("DATABASE_URL")

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		panic(err)
	}
	defer db.Close()

	// Test the connection
	if err = db.Ping(); err != nil {
		fmt.Printf("Failed to ping database: %v\n", err)
		panic(err)
	} else {
		fmt.Println("✅ Successfully connected to database!")
	}

	// Set up routes
	app.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world!")
	})

	// A test route to check Sentry: this will panic and Sentry will capture it
	app.GET("/panic", func(c *gin.Context) {
		panic("This is a test panic caught by Sentry!")
	})

	// And run it
	app.Run(":8080") // Standard port for DO App Platform Go apps is 8080
}
