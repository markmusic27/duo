package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markmusic27/workspace/handlers"
)

func main() {
	// Setup ENV
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading ENV: " + err.Error())
		return
	}

	// Setup API
	api := gin.Default()

	const port = "8080"

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api.GET("/sms", handlers.InboundSMSRequest)
	api.POST("/api", handlers.InboundHTTPRequest)

	log.Println("Starting server on port " + port)
	api.Run(":" + port)
}
