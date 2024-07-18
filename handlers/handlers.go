package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	process "github.com/markmusic27/workspace/utils"
)

func InboundHTTPRequest(c *gin.Context) {
	//TODO: Replace the logic with the HTTP logic
}

func VerifyEnv(c *gin.Context) {
	envs := []string{
		"PHONES",
		"OPENAI",
		"NOTION",
		"COURSEID",
		"TASKID",
		"PROJECTID",
		"AREAINTERESTID",
		"NOTESID",
		"TWILIO_ACCOUNT_SID",
		"TWILIO_AUTH_TOKEN",
		"TWILIO_PHONE",
		"GCP_KEY",
	}

	missing := []string{}

	for _, key := range envs {
		value := os.Getenv(key)

		if len(value) == 0 {
			missing = append(missing, key)
		}
	}

	if len(missing) == 0 {
		c.JSON(200, gin.H{
			"message": "All ENVs are loaded",
		})

		return
	}

	c.JSON(400, gin.H{
		"error": fmt.Sprintf("Error: The following ENVs have not been found: %s", missing),
	})
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

	// Add line to process
	err := process.Process(c.PostForm("Body"), c.PostForm("From"))

	if err != nil {
		process.Message(c.PostForm("From"), process.TruncateString(err.Error()))
	}
}
