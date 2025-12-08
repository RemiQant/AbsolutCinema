package main

import (
	"fmt"
	"net/http"
	"os"   // New import for os package
	"time" // New import for time package

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize Sentry with DSN from Environment
	dsn := os.Getenv("SENTRY_DSN")
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
		// Recommended for performance monitoring
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	// IMPORTANT: Ensure Sentry events are sent before the program terminates
	defer sentry.Flush(2 * time.Second)

	// Then create your app
	app := gin.Default()

	// Once it's done, you can attach the handler as one of your middleware
	app.Use(sentrygin.New(sentrygin.Options{}))

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
