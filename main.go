package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"authentication/middleware"
	"authentication/routes"

	"github.com/gin-gonic/gin"
)

func GenerateRandomKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func main() {
	key, err := GenerateRandomKey()

	if err != nil {
		log.Fatalf("Failed to generate random key: %v", err)
	}
	middleware.SetJWTKey([]byte(key))

	port := "8080"

	r := gin.Default()

	routes.SetupRoutes(r)

	log.Println("Server is running on port", port)
	r.Run(":" + port)
}