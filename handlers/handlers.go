package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	process "github.com/markmusic27/workspace/utils"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != fmt.Sprintf("Bearer %s", os.Getenv("HTTP_KEY")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
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
		"HTTP_KEY",
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

type APIResponse struct {
	Message string `json:"message"`
}

func InboundHTTPRequest(c *gin.Context) {
	var body APIResponse
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(body.Message) == 0 {
		c.JSON(400, gin.H{"error": fmt.Errorf("message is empty")})
		return
	}

	// Add line to process
	id, err := process.Process(body.Message)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(id) == 0 {
		c.JSON(400, gin.H{"error": fmt.Errorf("returned ID was empty")})
		return
	}

	c.JSON(200, gin.H{
		"id": id,
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
	_, err := process.Process(c.PostForm("Body"))

	if err != nil {
		process.Message(c.PostForm("From"), "Error: "+process.TruncateString(err.Error()))
	} else {
		process.Message(c.PostForm("From"), "Logged ✅")
	}
}
