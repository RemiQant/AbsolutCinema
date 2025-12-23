package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"absolutcinema-backend/internal/auth"
	"absolutcinema-backend/internal/controllers"
	"absolutcinema-backend/internal/middleware"
	"absolutcinema-backend/internal/services"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
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
	
	// Initialize controllers
	studioController := controllers.NewStudioController(studioService)
	movieController := controllers.NewMovieController(movieService)
	showtimeController := controllers.NewShowtimeController(showtimeService)

	// API group
	api := r.Group("/api")
	{
		// Auth routes (public)
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.Refresh)
			authRoutes.POST("/logout", authHandler.Logout)
			
			// Protected: Get current user
			authRoutes.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
		}
		
		// Public showtime routes (read-only)
		showtimeRoutes := api.Group("/showtimes")
		{
			showtimeRoutes.GET("", showtimeController.GetAllShowtimes)       // List with filters
			showtimeRoutes.GET("/:id", showtimeController.GetShowtimeByID)   // Get single showtime
		}
	}

	// Protected routes - require authentication
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Example: User routes (authenticated users only)
		protected.GET("/profile", s.getProfileHandler)
		
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
			adminRoutes.PUT("/showtimes/:id", showtimeController.UpdateShowtime)
			adminRoutes.DELETE("/showtimes/:id", showtimeController.DeleteShowtime)
			
			// Example: User management (keep existing)
			adminRoutes.GET("/users", s.getAllUsersHandler)
			adminRoutes.DELETE("/users/:id", s.deleteUserHandler)
		}
		
		// Example: Customer routes
		customerRoutes := protected.Group("/bookings")
		customerRoutes.Use(middleware.RequireAdminOrCustomer())
		{
			customerRoutes.GET("/", s.getBookingsHandler)
			customerRoutes.POST("/", s.createBookingHandler)
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

func (s *Server) getBookingsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get user bookings",
	})
}

func (s *Server) createBookingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Create booking",
	})
}
