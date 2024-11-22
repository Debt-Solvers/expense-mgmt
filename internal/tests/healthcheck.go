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



// func TestHealthCheck(t *testing.T) {
//     // Create a Gin router with the health check route
//     router := gin.Default()
//     AddHealthCheckRoute(router)

//     // Perform a GET request to the /health endpoint
//     req := httptest.NewRequest(http.MethodGet, "/health", nil)
//     rec := httptest.NewRecorder()
//     router.ServeHTTP(rec, req)

//     // Assert the response
//     assert.Equal(t, http.StatusOK, rec.Code)
//     assert.JSONEq(t, `{"status":"success","message":"Backend is reachable"}`, rec.Body.String())
// }
