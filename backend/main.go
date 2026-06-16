package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"budgetbooking/internal/db"
	"budgetbooking/internal/handlers"
	"budgetbooking/internal/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables")
	}

	pool, err := db.New(context.Background())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	expiryHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil || expiryHours <= 0 {
		expiryHours = 168
	}

	h := handlers.New(pool, os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("JWT_SECRET"), time.Duration(expiryHours)*time.Hour)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	user := router.Group("/user")
	user.POST("/login", h.Login)
	user.GET("/me", middleware.Auth(), h.Me)
	user.POST("/setting", middleware.Auth(), h.UpdateSetting)

	setBudget := router.Group("/setBudget", middleware.Auth())
	setBudget.POST("/add", h.AddBudget)
	setBudget.POST("/change", h.ChangeBudget)
	setBudget.POST("/delete", h.DeleteBudget)

	booking := router.Group("/booking", middleware.Auth())
	booking.POST("/add", h.AddBooking)
	booking.POST("/delete", h.DeleteBooking)

	report := router.Group("/report", middleware.Auth())
	report.GET("/showSpendingReport", h.ShowSpendingReport)
	report.GET("/showBudgetReport", h.ShowBudgetReport)
	report.GET("/showIncomeReport", h.ShowIncomeReport)
	report.GET("/availablePeriods", h.AvailablePeriods)
	report.GET("/showAnnualSpendingReport", h.ShowAnnualSpendingReport)
	report.GET("/showAnnualBudgetReport", h.ShowAnnualBudgetReport)
	report.GET("/showAnnualIncomeReport", h.ShowAnnualIncomeReport)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
