package db

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

var SecureMiddleware = secure.New(secure.Config{
	FrameDeny:             true,
	ContentTypeNosniff:    true,
	BrowserXssFilter:      true,
	ContentSecurityPolicy: "default-src 'self'",
})

func CorsConfig() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowOrigins:     []string{"file.pdfrouters.com", "http://localhost:3000", "78.172.216.87:3000", "78.172.216.87"},
		AllowCredentials: true,
		MaxAge:           60 * 24 * 30,
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := os.Getenv("SECRET_SESSION_KEY")
		fmt.Println("[SECRET-KEY]", secretKey)

		if secretKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		token, err := c.Cookie("session_token")
		if err != nil {
			fmt.Println("[ERROR]", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		fmt.Println("[TOKEN]", token)
		if token != secretKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		log.Println("Authentication successful")
		c.Next()
	}
}

func CookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
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
