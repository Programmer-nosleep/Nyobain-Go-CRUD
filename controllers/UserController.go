package controllers

import (
	"authentication/config"
	"authentication/helpers"
	"authentication/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

// REGISTER
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		// Bind & validate input
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if email already exists
		var existing models.User
		if err := config.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}

		// Hash password
		hashedPassword, err := helpers.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Assign fields
		user.Password = hashedPassword
		user.UserID = uuid.New().String()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		// Generate JWT tokens
		accessToken, refreshToken, err := helpers.GenerateTokens(user.UserID, user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		user.Token = accessToken
		user.RefreshToken = refreshToken

		// Save user
		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":       "Signup successful",
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

// LOGIN
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.User
		var found models.User

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := config.DB.Where("email = ?", input.Email).First(&found).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		accessToken, refreshToken, err := helpers.GenerateTokens(found.UserID, found.Email, found.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		// Optional: simpan token ke DB (boleh dihilangkan kalau tidak mau simpan)
		found.Token = accessToken
		found.RefreshToken = refreshToken
		_ = config.DB.Save(&found)

		c.JSON(http.StatusOK, gin.H{
			"message":       "Login successful",
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

// GET ALL USERS (admin only)
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsAny, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := claimsAny.(*helpers.Claims)
		if !ok || claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var users []models.User
		if err := config.DB.WithContext(ctx).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// GET SINGLE USER BY ID
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		var user models.User

		if err := config.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
