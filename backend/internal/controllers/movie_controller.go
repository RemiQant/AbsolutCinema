package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"absolutcinema-backend/internal/models"
	"absolutcinema-backend/internal/services"
)

type MovieController struct {
	service *services.MovieService
}

func NewMovieController(service *services.MovieService) *MovieController {
	return &MovieController{service: service}
}

// CreateMovie handles POST /api/admin/movies
func (mc *MovieController) CreateMovie(c *gin.Context) {
	var movie models.Movie
	
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	if err := mc.service.CreateMovie(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Movie created successfully",
		"data": movie,
	})
}

// GetAllMovies handles GET /api/admin/movies
func (mc *MovieController) GetAllMovies(c *gin.Context) {
	movies, err := mc.service.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve movies",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Movies retrieved successfully",
		"data": movies,
	})
}

// GetMovieByID handles GET /api/admin/movies/:id
func (mc *MovieController) GetMovieByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid movie ID",
		})
		return
	}
	
	movie, err := mc.service.GetMovieByID(uint(id))
	if err != nil {
		if err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve movie",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Movie retrieved successfully",
		"data": movie,
	})
}

// UpdateMovie handles PUT /api/admin/movies/:id
func (mc *MovieController) UpdateMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid movie ID",
		})
		return
	}
	
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	if err := mc.service.UpdateMovie(uint(id), &movie); err != nil {
		if err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, gin.H{
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
		"message": "Movie updated successfully",
	})
}

// DeleteMovie handles DELETE /api/admin/movies/:id
func (mc *MovieController) DeleteMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid movie ID",
		})
		return
	}
	
	if err := mc.service.DeleteMovie(uint(id)); err != nil {
		if err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete movie",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Movie deleted successfully",
	})
}
