package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "pwa-backend/docs"
	"pwa-backend/internal/config"
	"pwa-backend/internal/database"
	"pwa-backend/internal/handlers"
	"pwa-backend/internal/middleware"
	"pwa-backend/internal/repositories"
)

// @title PWA Offline-First Backend API
// @version 1.0
// @description Backend MVP untuk PWA dengan PowerSync
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @tokenUrl /api/v1/auth/login
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
	productHandler := handlers.NewProductHandler(productRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionRepo, productRepo)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/products", productHandler.GetProducts)
			protected.GET("/products/:id", productHandler.GetProductByID)

			protected.GET("/transactions/:id", transactionHandler.GetTransaction)
			protected.POST("/transactions/checkout", transactionHandler.Checkout)
			protected.PUT("/transactions/:id/status", transactionHandler.UpdateStatus)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}