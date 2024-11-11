package common

import (
	"expense-mgmt/db"
	"expense-mgmt/internal/models"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Checks if token is valid and in database
func IsTokenActive(token string) bool {
	// Get the DB instance
	db := db.GetDBInstance()

	err := db.Where("token = ?", token).First(&models.AuthToken{}).Error
	return err == nil
}

// GenerateRandomColor generates a random 6-character hexadecimal color code (e.g., "#A3C113")
func GenerateRandomColor() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Generate a random number and format it as a 6-digit hex code
	color := fmt.Sprintf("#%06X", r.Intn(16777215))

	// Return the color in uppercase
	return strings.ToUpper(color)
}
