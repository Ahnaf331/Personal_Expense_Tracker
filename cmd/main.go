package main

import (
	"log"

	"expense-tracker/internal/database"
	"expense-tracker/internal/handlers"
	"expense-tracker/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()

	authHandler := handlers.NewAuthHandler(db)
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		expenseHandler := handlers.NewExpenseHandler(db)
		api.POST("/expenses", expenseHandler.CreateExpense)
		api.GET("/expenses", expenseHandler.GetExpenses)
		api.GET("/expenses/:id", expenseHandler.GetExpense)
		api.PUT("/expenses/:id", expenseHandler.UpdateExpense)
		api.DELETE("/expenses/:id", expenseHandler.DeleteExpense)

		categoryHandler := handlers.NewCategoryHandler(db)
		api.POST("/categories", categoryHandler.CreateCategory)
		api.GET("/categories", categoryHandler.GetCategories)
		api.PUT("/categories/:id", categoryHandler.UpdateCategory)
		api.DELETE("/categories/:id", categoryHandler.DeleteCategory)

		analyticsHandler := handlers.NewAnalyticsHandler(db)
		api.GET("/analytics/monthly", analyticsHandler.GetMonthlySummary)
		api.GET("/analytics/categories", analyticsHandler.GetCategoryBreakdown)
	}

	log.Println("Server running on :8080")
	r.Run(":8080")
}
