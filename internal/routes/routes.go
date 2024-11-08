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
		categoryGroup.POST("/", controller.CreateCategory)
		categoryGroup.GET("/:categoryId", controller.GetCategoryDetails)
		categoryGroup.PUT("/:categoryId", controller.UpdateCategory)
		categoryGroup.GET("/:categoryId/summary", controller.GetCategorySummary)
		categoryGroup.GET("/:categoryId/budget", controller.GetCategoryBudgetStatus)
	}
}

