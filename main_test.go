package main

import (
	"bytes"
	"encoding/json"
	"gin-fleamarket/dto"
	"gin-fleamarket/infra"
	"gin-fleamarket/models"
	"gin-fleamarket/services"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gotest.tools/v3/assert"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load(".env.test"); err != nil {
		log.Fatalln("Error loading .env.test file")
	}

	// Add this line to install the missing package
	if err := exec.Command("go", "get", "gotest.tools/v3/assert").Run(); err != nil {
		log.Fatalln("Error installing gotest.tools/v3/assert package")
	}

	code := m.Run()

	os.Exit(code)
}

func setupTestData(db *gorm.DB) {
	items := []models.Item{
		{Name: "テストアイテム1", Price: 1000, Description: "", SoldOut: false, UserID: 1},
		{Name: "テストアイテム2", Price: 2000, Description: "テスト2", SoldOut: true, UserID: 1},
		{Name: "テストアイテム3", Price: 3000, Description: "テスト3", SoldOut: false, UserID: 1},
	}

	users := []models.User{
		{Email: "test1@example.com", Password: "test1pass"},
		{Email: "test2@example.com", Password: "test2pass"},
	}

	for _, user := range users {
		db.Create(&user)
	}
	for _, item := range items {
		db.Create(&item)
	}
}

func setup() *gin.Engine {
	db := infra.SetupDB()
	db.AutoMigrate(&models.Item{}, &models.User{})

	setupTestData(db)
	router := setupRouter(db)

	return router
}

func TestFindAll(t *testing.T) {
	// テストのセットアップ
	router := setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items", nil)

	// APIリクエストの実行
	router.ServeHTTP(w, req)

	// APIの実行結果を取得
	var res map[string][]models.Item
	json.Unmarshal([]byte(w.Body.String()), &res)

	// アサーション
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 3, len(res["data"]))
}

func TestCreate(t *testing.T) {
	// テストのセットアップ
	router := setup()

	token, err := services.CreateToken(1, "test1@example.com")
	assert.Equal(t, nil, err)

	createItemInput := dto.CreateItemInput{
		Name:        "テストアイテム4",
		Price:       4000,
		Description: "Createテスト",
	}
	reqBody, _ := json.Marshal(createItemInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+*token)

	// APIリクエストの実行
	router.ServeHTTP(w, req)

	// APIの実行結果を取得
	var res map[string]models.Item
	json.Unmarshal([]byte(w.Body.String()), &res)

	// アサーション
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, uint(4), res["data"].ID)
}

func TestCreateUnauthorized(t *testing.T) {
	// テストのセットアップ
	router := setup()

	createItemInput := dto.CreateItemInput{
		Name:        "テストアイテム4",
		Price:       4000,
		Description: "Createテスト",
	}
	reqBody, _ := json.Marshal(createItemInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))

	// APIリクエストの実行
	router.ServeHTTP(w, req)

	// APIの実行結果を取得
	var res map[string]models.Item
	json.Unmarshal([]byte(w.Body.String()), &res)

	// アサーション
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdate(t *testing.T) {
	// テストのセットアップ
	router := setup()

	token, err := services.CreateToken(1, "test1@example.com")
	assert.Equal(t, nil, err)

	description := "Updateテスト"
	updateItemInput := dto.UpdateItemInput{
		Description: &description,
	}
	reqBody, _ := json.Marshal(updateItemInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/items/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+*token)

	router.ServeHTTP(w, req)
	var res map[string]models.Item
	json.Unmarshal([]byte(w.Body.String()), &res)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, description, res["data"].Description)

}

func TestDelete(t *testing.T) {
	// テストのセットアップ
	router := setup()

	token, err := services.CreateToken(1, "test1@example.com")
	assert.Equal(t, nil, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/items/1", nil)
	req.Header.Set("Authorization", "Bearer "+*token)

	router.ServeHTTP(w, req)
	var res map[string]models.Item
	json.Unmarshal([]byte(w.Body.String()), &res)

	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/items/1", nil)
	req.Header.Set("Authorization", "Bearer "+*token)

	router.ServeHTTP(w, req)
	json.Unmarshal([]byte(w.Body.String()), &res)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
