package controllers

import (
	"bytes"
	"encoding/json"
	"gin-fleamarket/dto"
	"gin-fleamarket/infra"
	"gin-fleamarket/models"
	"gin-fleamarket/repositories"
	"gin-fleamarket/services"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gotest.tools/v3/assert"
	"gorm.io/gorm"
)

func setupAuthTest() (*gin.Engine, *gorm.DB) {
	// Load test environment
	if err := godotenv.Load("../.env.test"); err != nil {
		os.Setenv("ENV", "test") // Fallback if .env.test can't be loaded
	}

	// Setup database
	db := infra.SetupDB()
	db.AutoMigrate(&models.User{})

	// Clear any existing users
	db.Exec("DELETE FROM users")

	// Setup router
	r := gin.Default()
	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := NewAuthController(authService)

	// Setup routes
	authGroup := r.Group("/auth")
	authGroup.POST("/signup", authController.Signup)
	authGroup.POST("/login", authController.Login)

	return r, db
}

func TestSignup(t *testing.T) {
	r, _ := setupAuthTest()

	// Create signup request
	signupInput := dto.SignupInput{
		Email:    "test@example.com",
		Password: "password123",
	}
	reqBody, _ := json.Marshal(signupInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)
}
