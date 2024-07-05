package handlers

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	process "github.com/markmusic27/workspace/utils"
)

func InboundHTTPRequest(c *gin.Context) {
	//TODO: Build out HTTP version of the API
}

func InboundSMSRequest(c *gin.Context) {
	// Checks if the inbound request should be processed
	authenticatedDevices := strings.Split(os.Getenv("PHONES"), ",")
	var authenticated = false

	for _, phone := range authenticatedDevices {
		if phone == c.PostForm("From") {
			authenticated = true
		}
	}

	if !authenticated {
		c.JSON(401, gin.H{
			"error": "Phone number is not authorized",
		})

		return
	}

	c.JSON(200, gin.H{
		"status": "Message is being processed!",
	})

	process.Process(c.PostForm("Body"))
}
