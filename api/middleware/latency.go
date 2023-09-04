package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Latency() gin.HandlerFunc {
	return func(c *fiber.Ctx) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")
		//c.AbortWithStatusJSON(http.StatusConflict, gin.H{"status": false, "message": "error"})
		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}
