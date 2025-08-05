package main

import (
	"log"

	"authentication/config"
	"authentication/helpers"
	"authentication/models"
	"authentication/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect ke database dan auto migrate
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{})

	// Generate JWT Key
	key, err := config.GenerateRandomKey()
	if err != nil {
		log.Fatalf("Failed to generate random key: %v", err)
	}
	helpers.SetJWTKey([]byte(key))

	// Setup router
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*.html")

	// Default route
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"Title":   "Home | Nyobain Golang",
			"Message": "Halo, ini halaman utama",
		})
	})

	// Setup route grup lainnya
	routes.SetupRoutes(router)

	// Start server
	port := "8080"
	log.Println("Server is running on port", port)
	router.Run(":" + port)
}
