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

type UpdateTimezoneResponse struct {
	Location string `json:"location"`
}

func UpdateTimezone(c *gin.Context) {
	var body UpdateTimezoneResponse
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(body.Location) == 0 {
		c.JSON(400, gin.H{"error": fmt.Errorf("timezone is empty")})
		return
	}

	location, err := process.SetTimezoneFromLocation(body.Location)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("The timezone has been updated to %s", location),
	})
}

func GetTimezone(c *gin.Context) {
	timezone := os.Getenv("LOCATION")

	c.JSON(200, gin.H{
		"timezone": timezone,
	})
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
		"LOCATION",
		"COURSE_PAGE_ID",
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
			"message": "All of the environment variables have been loaded",
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

	if strings.Contains(body.Message, process.TimezonePrefix) {
		location, err := process.ExtractLocationFromSMS(body.Message)
		if err != nil {
			c.JSON(400, gin.H{"error": fmt.Errorf("Error: Failed to extract location from string.\n" + err.Error())})
			return
		}

		iana, err := process.SetTimezoneFromLocation(location)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message": fmt.Sprintf("Updated timezone to %s", iana+" ✅"),
		})
		return
	}

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

	if strings.Contains(c.PostForm("Body"), process.TimezonePrefix) {
		location, err := process.ExtractLocationFromSMS(c.PostForm("Body"))
		if err != nil {
			process.Message(c.PostForm("From"), "Error: Failed to extract location from string.\n"+err.Error())
			return
		}

		iana, err := process.SetTimezoneFromLocation(location)
		if err != nil {
			process.Message(c.PostForm("From"), "Error: "+err.Error())
			return
		}

		process.Message(c.PostForm("From"), fmt.Sprintf("Updated timezone to %s", iana+" ✅"))
		return
	}

	// Add line to process
	_, err := process.Process(c.PostForm("Body"))

	if err != nil {
		process.Message(c.PostForm("From"), "Error: "+process.TruncateString(err.Error()))
	} else {
		process.Message(c.PostForm("From"), "Logged ✅")
	}
}
