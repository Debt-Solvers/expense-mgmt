package tests

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Expense mgmt Backend is reachable",
	})
}
