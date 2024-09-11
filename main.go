package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/okanay/file-upload-go/db"
	"github.com/okanay/file-upload-go/internal/asset"
	"github.com/okanay/file-upload-go/internal/upload"
	"log"
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

	// ->> Routers Groups
	router := gin.Default()
	// ->> Middlewares
	router.Use(db.SecureMiddleware)
	router.Use(db.CorsConfig())
	router.Use(db.TimeoutMiddleware(150 * time.Second))

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

	// Upload Route
	router.POST("/upload", uploadHandler.UploadFile)

	// Assets Route
	router.GET("/assets/:filename", assetHandler.GetAsset)
	router.GET("/assets/all", assetHandler.GetAllAssets)

	// 404 Handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "The requested " + c.Request.URL.Path + " was not found."})
	})

	err = router.Run(":8080")
	if err != nil {
		return
	}
}
