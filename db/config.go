package db

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var SecureMiddleware = secure.New(secure.Config{
	FrameDeny:             true,
	ContentTypeNosniff:    true,
	BrowserXssFilter:      true,
	ContentSecurityPolicy: "default-src 'self'",
})

func CorsConfig() gin.HandlerFunc {
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowOrigins:     []string{"http://localhost:3000"}, // You can specify more origins if needed
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,                       // If you need credentials (cookies, HTTP auth), set this to true
		ExposeHeaders:    []string{"Content-Length"}, // This is optional, depending on your needs
		MaxAge:           12 * time.Hour,             // Cache the CORS response for 12 hours
	}

	return cors.New(config)
}
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		done := make(chan bool, 1)
		go func() {
			c.Next()
			done <- true
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			if !c.IsAborted() {
				fmt.Println("Timeout occurred")
				c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{"error": "Request Timeout"})
			}
		}
	}
}
