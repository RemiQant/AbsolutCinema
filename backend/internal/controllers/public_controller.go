package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"absolutcinema-backend/internal/services"
)

// PublicController handles public-facing, read-only endpoints
type PublicController struct {
	movieService    *services.MovieService
	showtimeService *services.ShowtimeService
	studioService   *services.StudioService
	bookingService  *services.BookingService
}

// NewPublicController creates a new public controller
func NewPublicController(
	movieService *services.MovieService,
	showtimeService *services.ShowtimeService,
	studioService *services.StudioService,
	bookingService *services.BookingService,
) *PublicController {
	return &PublicController{
		movieService:    movieService,
		showtimeService: showtimeService,
		studioService:   studioService,
		bookingService:  bookingService,
	}
}

// ListMovies handles GET /api/movies
// Returns a list of all movies (public, no auth required)
func (pc *PublicController) ListMovies(c *gin.Context) {
	movies, err := pc.movieService.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve movies",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movies retrieved successfully",
		"data":    movies,
		"count":   len(movies),
	})
}

// GetMovieDetails handles GET /api/movies/:id
// Returns movie details with associated showtimes (public, no auth required)
func (pc *PublicController) GetMovieDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid movie ID",
		})
		return
	}

	// Get movie with showtimes
	movie, err := pc.movieService.GetMovieWithShowtimes(uint(id))
	if err != nil {
		if err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve movie",
			"details": err.Error(),
		})
		return
	}

	// Build showtimes response with studio explicitly included
	showtimesResponse := make([]gin.H, len(movie.Showtimes))
	for i, st := range movie.Showtimes {
		showtimesResponse[i] = gin.H{
			"id":         st.ID,
			"movie_id":   st.MovieID,
			"studio_id":  st.StudioID,
			"start_time": st.StartTime,
			"end_time":   st.EndTime,
			"price":      st.Price,
			"studio":     st.Studio,
		}
	}

	// Build response with showtimes explicitly included
	c.JSON(http.StatusOK, gin.H{
		"message": "Movie retrieved successfully",
		"data": gin.H{
			"id":               movie.ID,
			"title":            movie.Title,
			"description":      movie.Description,
			"duration_minutes": movie.DurationMinutes,
			"poster_url":       movie.PosterURL,
			"rating":           movie.Rating,
			"showtimes":        showtimesResponse,
		},
	})
}

// GetStudioLayout handles GET /api/studios/:id
// Returns studio seat layout/dimensions (public, no auth required)
func (pc *PublicController) GetStudioLayout(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid studio ID",
		})
		return
	}

	studio, err := pc.studioService.GetStudioByID(uint(id))
	if err != nil {
		if err.Error() == "studio not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve studio",
			"details": err.Error(),
		})
		return
	}

	// Generate seat layout for frontend
	seatLayout := generateSeatLayout(studio.TotalRows, studio.TotalCols)

	c.JSON(http.StatusOK, gin.H{
		"message": "Studio layout retrieved successfully",
		"data": gin.H{
			"id":         studio.ID,
			"name":       studio.Name,
			"total_rows": studio.TotalRows,
			"total_cols": studio.TotalCols,
			"seats":      seatLayout,
		},
	})
}

// GetOccupiedSeats handles GET /api/showtimes/:id/seats
// Returns occupied seats for a specific showtime (public, no auth required)
func (pc *PublicController) GetOccupiedSeats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid showtime ID",
		})
		return
	}

	// First verify showtime exists and get studio info
	showtime, err := pc.showtimeService.GetShowtimeByID(uint(id))
	if err != nil {
		if err.Error() == "showtime not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve showtime",
			"details": err.Error(),
		})
		return
	}

	// Get occupied seats for this showtime
	occupiedSeats, err := pc.bookingService.GetOccupiedSeats(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve occupied seats",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Occupied seats retrieved successfully",
		"data": gin.H{
			"showtime_id":    showtime.ID,
			"movie":          showtime.Movie.Title,
			"studio":         showtime.Studio.Name,
			"start_time":     showtime.StartTime,
			"price":          showtime.Price,
			"occupied_seats": occupiedSeats,
			"total_occupied": len(occupiedSeats),
		},
	})
}

// generateSeatLayout creates a 2D representation of seats
// Rows are labeled A-Z, columns are numbered 1-N
func generateSeatLayout(rows, cols int) [][]string {
	layout := make([][]string, rows)
	for i := 0; i < rows; i++ {
		rowLabel := string(rune('A' + i))
		layout[i] = make([]string, cols)
		for j := 0; j < cols; j++ {
			layout[i][j] = rowLabel + strconv.Itoa(j+1)
		}
	}
	return layout
}
