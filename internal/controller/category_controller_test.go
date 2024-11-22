package controller

/* Endpoints Covered:

GET /api/v1/categories/defaults
GET /api/v1/categories/
POST /api/v1/categories/
GET /api/v1/categories/:categoryId
PUT /api/v1/categories/:categoryId
DELETE /api/v1/categories/:categoryId

*/

import (
	"expense-mgmt/internal/middleware"
	"expense-mgmt/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Database
type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetDefaultCategories() ([]models.Category, error) {
	args := m.Called()
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockDB) GetUserCategories(userID string) ([]models.Category, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockDB) CreateCategory(category models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockDB) GetCategoryDetails(categoryID string) (models.Category, error) {
	args := m.Called(categoryID)
	return args.Get(0).(models.Category), args.Error(1)
}

func (m *MockDB) UpdateCategory(category models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockDB) DeleteCategory(categoryID string) error {
	args := m.Called(categoryID)
	return args.Error(0)
}

// Test for Category Routes
func TestCategoryRoutes(t *testing.T) {
	mockDB := new(MockDB)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())

	// Define the routes
	// routes.CategoryRoutes(router)

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzIzMTQ2NTgsInVzZXJfaWQiOiIwNmQ1MDU0NS1mMzNhLTRmZjktYTExZC01YjQ4ZDI2MGExNmQifQ.poDbCTsJJrHz0gohE_JQcCNNgDp2oNb_oDtXV8QA7bg"
	// Create a request with the Authorization header
	req, _ := http.NewRequest("GET", "/api/v1/categories/defaults", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	// Mock Data for Test Cases
	defaultCategories := []models.Category{
		{ID: uuid.Must(uuid.Parse("8ff95915-e811-4ff8-afd4-fe71daaa5f85")), Name: "Food"},
		{ID: uuid.Must(uuid.Parse("9aa6d915-13a6-4f3d-895c-378d3f84c5e9")), Name: "Transport"},
	}

	// Test Case for "Get Default Categories"
	t.Run("Get Default Categories - Success", func(t *testing.T) {
		mockDB.On("GetDefaultCategories").Return(defaultCategories, nil)

		req, _ := http.NewRequest("GET", "/api/v1/categories/defaults", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Food")
		mockDB.AssertExpectations(t)
	})

	// // Test Case for "List User Categories"
	// t.Run("List User Categories - Success", func(t *testing.T) {
	// 	mockDB.On("GetUserCategories", "userID").Return(defaultCategories, nil)

	// 	req, _ := http.NewRequest("GET", "/api/v1/categories/", nil)
	// 	w := httptest.NewRecorder()
	// 	router.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	mockDB.AssertExpectations(t)
	// })

	// // Test Case for "Create Custom Category"
	// t.Run("Create Custom Category - Success", func(t *testing.T) {
	// 	newCategory := models.Category{Name: "Entertainment"}
	// 	mockDB.On("CreateCategory", newCategory).Return(nil)

	// 	req, _ := http.NewRequest("POST", "/api/v1/categories/", nil) // Include category data in request body
	// 	w := httptest.NewRecorder()
	// 	router.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusCreated, w.Code)
	// 	mockDB.AssertExpectations(t)
	// })

	// // Test Case for "Get Category Details"
	// t.Run("Get Category Details - Success", func(t *testing.T) {
	// 	mockDB.On("GetCategoryDetails", "123").Return(defaultCategories[0], nil)

	// 	req, _ := http.NewRequest("GET", "/api/v1/categories/123", nil)
	// 	w := httptest.NewRecorder()
	// 	router.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	assert.Contains(t, w.Body.String(), "Food")
	// 	mockDB.AssertExpectations(t)
	// })

	// // Test Case for "Update Category"
	// t.Run("Update Category - Success", func(t *testing.T) {
		
	// 	updatedCategory := models.Category{ID: uuid.Must(uuid.Parse("123e4567-e89b-12d3-a456-426614174000")), Name: "Food & Dining"}
	// 	mockDB.On("UpdateCategory", updatedCategory).Return(nil)

	// 	req, _ := http.NewRequest("PUT", "/api/v1/categories/123", nil) // Include updated category data
	// 	w := httptest.NewRecorder()
	// 	router.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	mockDB.AssertExpectations(t)
	// })

	// // Test Case for "Delete Category"
	// t.Run("Delete Category - Success", func(t *testing.T) {
	// 	mockDB.On("DeleteCategory", "123").Return(nil)

	// 	req, _ := http.NewRequest("DELETE", "/api/v1/categories/123", nil)
	// 	w := httptest.NewRecorder()
	// 	router.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusNoContent, w.Code)
	// 	mockDB.AssertExpectations(t)
	// })
}
