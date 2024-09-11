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
	"strings"
	"time"
)

var SecureMiddleware = secure.New(secure.Config{
	FrameDeny:             true,
	ContentTypeNosniff:    true,
	BrowserXssFilter:      true,
	ContentSecurityPolicy: "default-src 'self'",
})

func CorsConfig() gin.HandlerFunc {
	allowedOriginEnv := os.Getenv("ALLOWED_ORIGIN")
	if allowedOriginEnv == "" {
		allowedOriginEnv = "http://localhost:3000" // Bu kısma backend'in çalıştığı URL'yi ekle
	}
	origins := strings.Split(allowedOriginEnv, " ")

	return cors.New(cors.Config{
		AllowOrigins:     origins, // Dikkat: '*' yerine belirli bir origin ayarla
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Cookie'ler için gerekli
		MaxAge:           12 * time.Hour,
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
		cookieOptions := &http.Cookie{
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 30, // 30 gün
			HttpOnly: true,
			Secure:   true, // HTTPS için zorunlu
			SameSite: http.SameSiteNoneMode,
		}

		c.Set("cookie_options", cookieOptions)
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
