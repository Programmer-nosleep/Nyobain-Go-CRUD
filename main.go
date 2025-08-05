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
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{})

	key, err := config.GenerateRandomKey()

	if err != nil {
		log.Fatalf("Failed to generate random key: %v", err)
	}
	helpers.SetJWTKey([]byte(key))

	port := "8080"
	r := gin.Default()

	r.LoadHTMLFiles("views/*")
	routes.SetupRoutes(r)

	log.Println("Server is running on port", port)
	r.Run(":" + port)
}