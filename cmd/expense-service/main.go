package main

import (
	"expense-mgmt/db"
	"expense-mgmt/internal/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	if _,err := db.ConnectDatabase(); err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Seed default categories
	if err := db.SeedDefaultCategories(db.DB); err != nil {
		log.Fatalf("Failed to seed default categories: %v", err)
	}

	// Initialize Gin engine
	server := gin.Default()

	// Register routes
	routes.CategoryRoutes(server) // Public routes
  routes.ExpenseRoutes(server)
  routes.BudgetRoutes(server)
	routes.AddHealthCheckRoute(server)
	// Check for environment variable port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Start the server
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}