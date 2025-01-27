package main

import (
	"encoding/json"
	"gin-fleamarket/infra"
	"gin-fleamarket/models"
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
