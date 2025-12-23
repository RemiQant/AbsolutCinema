package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"absolutcinema-backend/internal/models"
	"absolutcinema-backend/internal/services"
)

type StudioController struct {
	service *services.StudioService
}

func NewStudioController(service *services.StudioService) *StudioController {
	return &StudioController{service: service}
}

// CreateStudio handles POST /api/admin/studios
func (sc *StudioController) CreateStudio(c *gin.Context) {
	var studio models.Studio
	
	if err := c.ShouldBindJSON(&studio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	if err := sc.service.CreateStudio(&studio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Studio created successfully",
		"data": studio,
	})
}

// GetAllStudios handles GET /api/admin/studios
func (sc *StudioController) GetAllStudios(c *gin.Context) {
	studios, err := sc.service.GetAllStudios()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve studios",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Studios retrieved successfully",
		"data": studios,
	})
}

// GetStudioByID handles GET /api/admin/studios/:id
func (sc *StudioController) GetStudioByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid studio ID",
		})
		return
	}
	
	studio, err := sc.service.GetStudioByID(uint(id))
	if err != nil {
		if err.Error() == "studio not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve studio",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Studio retrieved successfully",
		"data": studio,
	})
}

// UpdateStudio handles PUT /api/admin/studios/:id
func (sc *StudioController) UpdateStudio(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid studio ID",
		})
		return
	}
	
	var studio models.Studio
	if err := c.ShouldBindJSON(&studio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	if err := sc.service.UpdateStudio(uint(id), &studio); err != nil {
		if err.Error() == "studio not found" {
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
		"message": "Studio updated successfully",
	})
}

// DeleteStudio handles DELETE /api/admin/studios/:id
func (sc *StudioController) DeleteStudio(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid studio ID",
		})
		return
	}
	
	if err := sc.service.DeleteStudio(uint(id)); err != nil {
		if err.Error() == "studio not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete studio",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Studio deleted successfully",
	})
}
