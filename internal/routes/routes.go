package routes

import (
	"expense-mgmt/internal/controller"
	"expense-mgmt/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.Engine) {
	categoryGroup := router.Group("/api/v1/categories")
	categoryGroup.Use(middleware.AuthMiddleware())
	{
		categoryGroup.GET("/defaults", controller.GetDefaultCategories)
		categoryGroup.GET("/", controller.ListUserCategories)
		categoryGroup.POST("/", controller.CreateCustomCategory)
		categoryGroup.GET("/:categoryId", controller.GetCategoryDetails)
		categoryGroup.PUT("/:categoryId", controller.UpdateCategory)
		categoryGroup.DELETE("/:categoryId", controller.DeleteCategory)
	}
}

func ExpenseRoutes(router *gin.Engine) {
	expenseGroup := router.Group("/api/v1/expenses")
	expenseGroup.Use(middleware.AuthMiddleware())
	{
		expenseGroup.POST("/", controller.CreateExpense)
		expenseGroup.GET("/", controller.ListUserExpenses)
		expenseGroup.GET("/:expenseId", controller.GetExpense)
		expenseGroup.DELETE("/:expenseId", controller.DeleteExpense)
		expenseGroup.PUT("/:expenseId", controller.UpdateExpense)
		expenseGroup.GET("/analysis", controller.ExpenseAnalysis)
	}
}

func BudgetRoutes(router *gin.Engine) {
	budgetGroup := router.Group("/api/v1/budgets")
	budgetGroup.Use(middleware.AuthMiddleware())
	{
		budgetGroup.POST("/", controller.CreateBudget)   
		budgetGroup.GET("/", controller.ListBudgets)      
		budgetGroup.GET("/:budgetId", controller.GetSingleBudget) 
		budgetGroup.PUT("/:budgetId", controller.UpdateBudget)   
		budgetGroup.DELETE("/:budgetId", controller.DeleteBudget)
		budgetGroup.GET("/analysis", controller.BudgetAnalysis)
	}
}
