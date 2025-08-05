package controllers

import (
	"authentication/config"
	"authentication/helpers"
	"authentication/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		var existingUser models.User
		config.DB.Where("email = ?", user.Email).First(&existingUser)
		if existingUser.ID != 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}

		hashedPassword, err := helpers.HashPassword(user.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.Password = string(hashedPassword)
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.UserID = uuid.New().String()

		accessToken, refreshToken, err := helpers.GenerateTokens(user.UserID, user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		user.Token = accessToken
		user.RefreshToken = refreshToken

		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "Signup successful",
			"access_token":   accessToken,
			"refresh_token":  refreshToken,
		})
	}
}


func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		config.DB.Where("email = ?", user.Email).First(&foundUser)
		if foundUser.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		accessToken, refreshToken, err := helpers.GenerateTokens(foundUser.UserID, foundUser.Email, foundUser.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		foundUser.Token = accessToken
		foundUser.RefreshToken = refreshToken
		config.DB.Save(&foundUser)

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "access_token": accessToken, "refresh_token": refreshToken})
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		config.DB.Find(&users)
		c.JSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		var user models.User
		config.DB.First(&user, userID)
		c.JSON(http.StatusOK, user)
	}
}