package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/okanay/file-upload-go/db"
	"github.com/okanay/file-upload-go/internal/asset"
	"github.com/okanay/file-upload-go/internal/upload"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Load Environments
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Error loading .env file")
		return
	}
	// Set Database Connection
	sqlDB, err := db.Init(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close(sqlDB)

	// ->> Middlewares
	router := gin.Default()
	router.Use(db.SecureMiddleware)
	router.Use(db.CorsConfig())
	router.Use(db.TimeoutMiddleware(150 * time.Second))

	store := cookie.NewStore([]byte("your-secret-key"))
	router.Use(sessions.Sessions("my-session", store))

	// ->> Auth Middleware
	auth := router.Group("auth")
	auth.Use(db.AuthMiddleware())

	// Repositories
	uploadRepo := upload.NewRepository(sqlDB)
	assetRepo := asset.NewRepository(sqlDB)
	// Services
	uploadService := upload.NewService(uploadRepo)
	assetService := asset.NewService(assetRepo)
	// Handlers
	uploadHandler := upload.NewHandler(uploadService)
	assetHandler := asset.NewAssetHandler(assetService, "./public", "./public/blur", "./public/optimized", true, 60*time.Minute)

	// Main Route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to File Upload API", "Language": "Go Lang", "Framework": "Gin Gonic"})
	})

	// Assets Route
	router.GET("/assets/:filename", assetHandler.GetAsset)
	router.GET("/assets/all", assetHandler.GetAllAssets)

	// Auth Routes
	auth.POST("/upload", uploadHandler.UploadFile)
	auth.POST("/assets/delete", assetHandler.DeleteAsset)

	// Login Route
	router.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("session_token", "your-session-token")
		session.Set("auth-status", "login")
		session.Save()

		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	})

	// Logout Route
	auth.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear() // Tüm session verilerini temizle
		session.Save()

		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	})

	// 404 Handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "The requested " + c.Request.URL.Path + " was not found."})
	})

	err = router.Run(":8080")
	if err != nil {
		return
	}
}
