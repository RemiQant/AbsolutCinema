package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"absolutcinema-backend/internal/models"
	"absolutcinema-backend/internal/services"
)

type ShowtimeController struct {
	service *services.ShowtimeService
}

func NewShowtimeController(service *services.ShowtimeService) *ShowtimeController {
	return &ShowtimeController{service: service}
}

// CreateShowtimeRequest represents the request body for creating a showtime
type CreateShowtimeRequest struct {
	MovieID   uint      `json:"movie_id" binding:"required"`
	StudioID  uint      `json:"studio_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	Price     float64   `json:"price" binding:"required,gt=0"`
}

// UpdateShowtimeRequest represents the request body for updating a showtime
type UpdateShowtimeRequest struct {
	MovieID   uint      `json:"movie_id" binding:"required"`
	StudioID  uint      `json:"studio_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	Price     float64   `json:"price" binding:"required,gt=0"`
}

// CreateShowtime handles
// POST /api/admin/showtimes
func (sc *ShowtimeController) CreateShowtime(c *gin.Context) {
	var req CreateShowtimeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	showtime := models.Showtime{
		MovieID:   req.MovieID,
		StudioID:  req.StudioID,
		StartTime: req.StartTime,
		Price:     req.Price,
	}

	if err := sc.service.CreateShowtime(&showtime); err != nil {
		// Check for specific error types
		if err.Error() == "movie not found" || err.Error() == "studio not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err.Error() == "schedule conflict: studio is already occupied during this time slot" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Showtime created successfully",
		"data":    showtime,
	})
}

// GetAllShowtimes handles
// GET /api/showtimes (Public) and GET /api/admin/showtimes
func (sc *ShowtimeController) GetAllShowtimes(c *gin.Context) {
	// Optional query parameters for filtering
	var movieID *uint
	var date *time.Time

	// Parse movie_id query parameter
	if movieIDStr := c.Query("movie_id"); movieIDStr != "" {
		id, err := strconv.ParseUint(movieIDStr, 10, 32)
		if err == nil {
			movieIDUint := uint(id)
			movieID = &movieIDUint
		}
	}

	// Parse date query parameter (format: YYYY-MM-DD)
	if dateStr := c.Query("date"); dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = &parsedDate
		}
	}

	showtimes, err := sc.service.GetAllShowtimes(movieID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve showtimes",
			"details": err.Error(),
		})
		return
	}

	// Build response with movie and studio explicitly included
	response := make([]gin.H, len(showtimes))
	for i, st := range showtimes {
		response[i] = gin.H{
			"id":         st.ID,
			"movie_id":   st.MovieID,
			"studio_id":  st.StudioID,
			"start_time": st.StartTime,
			"end_time":   st.EndTime,
			"price":      st.Price,
			"movie":      st.Movie,
			"studio":     st.Studio,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Showtimes retrieved successfully",
		"data":    response,
	})
}

// GetShowtimeByID handles
// GET /api/showtimes/:id (Public) and GET /api/admin/showtimes/:id
func (sc *ShowtimeController) GetShowtimeByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid showtime ID",
		})
		return
	}

	showtime, err := sc.service.GetShowtimeByID(uint(id))
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Showtime retrieved successfully",
		"data": gin.H{
			"id":         showtime.ID,
			"movie_id":   showtime.MovieID,
			"studio_id":  showtime.StudioID,
			"start_time": showtime.StartTime,
			"end_time":   showtime.EndTime,
			"price":      showtime.Price,
			"movie":      showtime.Movie,
			"studio":     showtime.Studio,
		},
	})
}

// UpdateShowtime handles 
// PUT /api/admin/showtimes/:id
func (sc *ShowtimeController) UpdateShowtime(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid showtime ID",
		})
		return
	}

	var req UpdateShowtimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	showtime := models.Showtime{
		MovieID:   req.MovieID,
		StudioID:  req.StudioID,
		StartTime: req.StartTime,
		Price:     req.Price,
	}

	if err := sc.service.UpdateShowtime(uint(id), &showtime); err != nil {
		if err.Error() == "showtime not found" || err.Error() == "movie not found" || err.Error() == "studio not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err.Error() == "schedule conflict: studio is already occupied during this time slot" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Showtime updated successfully",
	})
}

// DeleteShowtime handles 
// DELETE /api/admin/showtimes/:id
func (sc *ShowtimeController) DeleteShowtime(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid showtime ID",
		})
		return
	}

	if err := sc.service.DeleteShowtime(uint(id)); err != nil {
		if err.Error() == "showtime not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete showtime",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Showtime deleted successfully",
	})
}
