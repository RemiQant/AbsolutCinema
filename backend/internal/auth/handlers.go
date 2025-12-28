package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"absolutcinema-backend/internal/models"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the user data returned in responses
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Register handles user registration
// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already registered",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process registration",
		})
		return
	}

	// Create user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "customer", // Default role
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Generate access token
	accessToken, err := GenerateAccessToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	// Generate refresh token
	refreshToken, expiresAt, err := GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		return
	}

	// Store refresh token in database
	dbRefreshToken := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}

	if err := h.db.Create(&dbRefreshToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save refresh token",
		})
		return
	}

	// Set cookies
	setAuthCookies(c, accessToken, refreshToken)

	// Return user data without password
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Login handles user login and sets auth cookies
// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Find user by email
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Generate access token
	accessToken, err := GenerateAccessToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	// Generate refresh token
	refreshToken, expiresAt, err := GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		return
	}

	// Store refresh token in database
	dbRefreshToken := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}

	if err := h.db.Create(&dbRefreshToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save refresh token",
		})
		return
	}

	// Set cookies
	setAuthCookies(c, accessToken, refreshToken)

	// Return user data
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Refresh handles access token refresh using refresh token
// POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token not found",
		})
		return
	}

	// Validate refresh token
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired refresh token",
		})
		return
	}

	// Check if refresh token exists in database and is not revoked
	var dbRefreshToken models.RefreshToken
	if err := h.db.Where("token = ? AND user_id = ?", refreshToken, claims.UserID).
		First(&dbRefreshToken).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token not found or revoked",
		})
		return
	}

	// Check if token is expired
	if dbRefreshToken.IsExpired() {
		h.db.Delete(&dbRefreshToken)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token has expired",
		})
		return
	}

	// Get user data
	var user models.User
	if err := h.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// Generate new access token
	newAccessToken, err := GenerateAccessToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	// REFRESH TOKEN ROTATION: Generate new refresh token
	newRefreshToken, expiresAt, err := GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		return
	}

	// Delete old refresh token and create new one (rotation)
	h.db.Delete(&dbRefreshToken)
	newDBRefreshToken := models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: expiresAt,
	}

	if err := h.db.Create(&newDBRefreshToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save refresh token",
		})
		return
	}

	// Set new cookies
	setAuthCookies(c, newAccessToken, newRefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Tokens refreshed successfully",
	})
}

// Logout handles user logout and invalidates refresh token
// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		// Delete refresh token from database
		h.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{})
	}

	// Clear cookies
	clearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// GetCurrentUser returns the currently authenticated user
// GET /auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Get user from database
	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	})
}

// setAuthCookies sets access and refresh token cookies
func setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	// Determine if we're in production (for Secure flag)
	isProduction := strings.ToLower(os.Getenv("APP_ENV")) == "production"

	// Use SameSite=Lax for both local dev and production
	// This works because frontend and backend are now on the same domain:
	// - Local: both on localhost
	// - Production: both on absolut-cinema-umwih.ondigitalocean.app
	c.SetSameSite(http.SameSiteLaxMode)

	// Set access token cookie
	c.SetCookie(
		"access_token",                          // name
		accessToken,                             // value
		int(GetAccessTokenDuration().Seconds()), // maxAge in seconds
		"/",                                     // path
		"",                                      // domain (empty for current domain)
		isProduction,                            // secure (HTTPS only in production)
		true,                                    // httpOnly
	)

	// Set refresh token cookie
	c.SetCookie(
		"refresh_token",                          // name
		refreshToken,                             // value
		int(GetRefreshTokenDuration().Seconds()), // maxAge in seconds
		"/",                                      // path
		"",                                       // domain
		isProduction,                             // secure
		true,                                     // httpOnly
	)
}

// clearAuthCookies clears authentication cookies
func clearAuthCookies(c *gin.Context) {
	isProduction := strings.ToLower(os.Getenv("APP_ENV")) == "production"

	// Use same SameSite policy as when setting cookies
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		"access_token",
		"",
		-1, // maxAge -1 deletes the cookie
		"/",
		"",
		isProduction,
		true,
	)

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		isProduction,
		true,
	)
}
