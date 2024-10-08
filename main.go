package main

import (
	"log"

	"github.com/joho/godotenv"
	process "github.com/markmusic27/workspace/utils"
)

func main() {
	// Setup ENV
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading ENV: " + err.Error())
		return
	}

	// Setup API
	// api := gin.Default()

	// const port = "8080"

	// api.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	// api.GET("/env", handlers.VerifyEnv)
	// api.GET("/timezone", handlers.GetTimezone)

	// api.POST("/sms", handlers.InboundSMSRequest)
	// api.POST("/api", handlers.Authenticate(), handlers.InboundHTTPRequest)
	// api.POST("/set-timezone", handlers.Authenticate(), handlers.UpdateTimezone)

	// log.Println("Starting server on port " + port)
	// api.Run(":" + port)

	process.IngestNote("Look into Interface for next AI project (cheap compute) https://www.inference.net/")
}
