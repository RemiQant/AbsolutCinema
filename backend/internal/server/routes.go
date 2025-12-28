package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"absolutcinema-backend/internal/auth"
	"absolutcinema-backend/internal/controllers"
	"absolutcinema-backend/internal/middleware"
	"absolutcinema-backend/internal/services"
)

func (s *Server) RegisterRoutes() http.Handler {
	// Initialize Sentry
	sentryDsn := os.Getenv("SENTRY_DSN")
	if sentryDsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: sentryDsn,
		}); err != nil {
			fmt.Printf("Sentry initialization failed: %v\n", err)
		}
	} else {
		log.Println("Sentry DSN not configured, skipping Sentry initialization")
	}

	r := gin.Default()

	r.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",                           // Local development
			"https://absolut-cinema.vercel.app",               // Vercel (legacy)
			"https://absolut-cinema-umwih.ondigitalocean.app", // DigitalOcean frontend
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Public routes
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	// Initialize auth handler
	authHandler := auth.NewAuthHandler(s.db.DB())

	// Initialize services
	studioService := services.NewStudioService(s.db.DB())
	movieService := services.NewMovieService(s.db.DB())
	showtimeService := services.NewShowtimeService(s.db.DB())

	// Initialize payment service (optional - may fail if XENDIT_SECRET_KEY not set)
	var paymentService *services.PaymentService
	var err error
	paymentService, err = services.NewPaymentService()
	if err != nil {
		log.Printf("Warning: Payment service initialization failed: %v", err)
		log.Println("Bookings will be created without payment links")
		paymentService = nil
	}

	// Initialize booking service
	bookingService := services.NewBookingService(s.db.DB(), paymentService)

	// Initialize controllers
	studioController := controllers.NewStudioController(studioService)
	movieController := controllers.NewMovieController(movieService)
	showtimeController := controllers.NewShowtimeController(showtimeService)
	bookingController := controllers.NewBookingController(bookingService)
	publicController := controllers.NewPublicController(movieService, showtimeService, studioService, bookingService)
	webhookController := controllers.NewWebhookController(bookingService)

	// Note: DigitalOcean routes /api/* to this backend, so we don't need /api prefix here
	// Routes are defined from root since DO strips the /api prefix

	// Auth routes (public)
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)
		authRoutes.POST("/logout", authHandler.Logout)

		// Protected: Get current user
		authRoutes.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
	}

	// Public showtime routes (read-only)
	showtimeRoutes := r.Group("/showtimes")
	{
		showtimeRoutes.GET("", showtimeController.GetAllShowtimes)          // List with filters
		showtimeRoutes.GET("/:id", showtimeController.GetShowtimeByID)      // Get single showtime
		showtimeRoutes.GET("/:id/seats", publicController.GetOccupiedSeats) // Get occupied seats
	}

	// Public movie routes (read-only)
	movieRoutes := r.Group("/movies")
	{
		movieRoutes.GET("", publicController.ListMovies)          // List all movies
		movieRoutes.GET("/:id", publicController.GetMovieDetails) // Movie details + showtimes
	}

	// Public studio routes (read-only for seat layout)
	studioRoutes := r.Group("/studios")
	{
		studioRoutes.GET("/:id", publicController.GetStudioLayout) // Get studio seat layout
	}

	// Webhook routes (public but secured by callback token)
	// IMPORTANT: These routes must NOT have JWT middleware
	// Security is handled by validating the x-callback-token header
	webhookRoutes := r.Group("/webhooks")
	{
		webhookRoutes.POST("/xendit", webhookController.HandleXenditCallback)
	}

	// Protected routes - require authentication
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Example: User routes (authenticated users only)
		protected.GET("/profile", s.getProfileHandler)

		// Booking routes (Customer/Admin)
		bookingRoutes := protected.Group("/bookings")
		bookingRoutes.Use(middleware.RequireAdminOrCustomer())
		{
			bookingRoutes.GET("", bookingController.GetBookings)                     // List own bookings
			bookingRoutes.POST("", bookingController.CreateBooking)                  // Create booking
			bookingRoutes.GET("/:id", bookingController.GetBookingByID)              // Get booking by ID
			bookingRoutes.DELETE("/:id", bookingController.CancelBooking)            // Cancel booking
			bookingRoutes.POST("/:id/retry-payment", bookingController.RetryPayment) // Retry payment
		}

		// Admin-only routes for Master Data Management
		adminRoutes := protected.Group("/admin")
		adminRoutes.Use(middleware.RequireAdmin())
		{
			// Studio CRUD endpoints
			adminRoutes.POST("/studios", studioController.CreateStudio)
			adminRoutes.GET("/studios", studioController.GetAllStudios)
			adminRoutes.GET("/studios/:id", studioController.GetStudioByID)
			adminRoutes.PUT("/studios/:id", studioController.UpdateStudio)
			adminRoutes.DELETE("/studios/:id", studioController.DeleteStudio)

			// Movie CRUD endpoints
			adminRoutes.POST("/movies", movieController.CreateMovie)
			adminRoutes.GET("/movies", movieController.GetAllMovies)
			adminRoutes.GET("/movies/:id", movieController.GetMovieByID)
			adminRoutes.PUT("/movies/:id", movieController.UpdateMovie)
			adminRoutes.DELETE("/movies/:id", movieController.DeleteMovie)

			// Showtime CRUD endpoints (Admin only for CUD operations)
			adminRoutes.POST("/showtimes", showtimeController.CreateShowtime)
			adminRoutes.GET("/showtimes", showtimeController.GetAllShowtimes)
			adminRoutes.PUT("/showtimes/:id", showtimeController.UpdateShowtime)
			adminRoutes.DELETE("/showtimes/:id", showtimeController.DeleteShowtime)

			// Example: User management (keep existing)
			adminRoutes.GET("/users", s.getAllUsersHandler)
			adminRoutes.DELETE("/users/:id", s.deleteUserHandler)
		}
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

// Example handler implementations (you can implement these based on your needs)
func (s *Server) getProfileHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User profile",
		"user_id": userID,
	})
}

func (s *Server) getAllUsersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin: Get all users",
	})
}

func (s *Server) deleteUserHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin: Delete user",
	})
}
