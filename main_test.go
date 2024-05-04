package main

import (
	"bytes"
	"encoding/json"
	"gin-freemarket-app/dto"
	"gin-freemarket-app/infra"
	"gin-freemarket-app/models"
	"gin-freemarket-app/services"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

/*
*testing.Mはmain関数用のテストランナー
他のテスト関数が呼び出される前に実行される特別なテスト関数
*/
func TestMain(m *testing.M) {
	if err := godotenv.Load(".env.test"); err != nil {
		log.Fatalln("Rrror loading .env.test file")
	}

	code := m.Run()

	os.Exit(code)
}

/*
sqliteに登録するテストデータをセットアップする関数
*/
func setupTestData(db *gorm.DB) {
	items := []models.Item{
		{Name: "testItem1", Price: 1000, Description: "", SoldOut: false, UserID: 1},
		{Name: "testItem2", Price: 2000, Description: "test2", SoldOut: true, UserID: 1},
		{Name: "testItem3", Price: 3000, Description: "test3", SoldOut: false, UserID: 2},
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

/*
test用のDBとrouterをセットアップする関数
*/
func setup() *gin.Engine {
	db := infra.SetupDB()
	db.AutoMigrate(&models.Item{}, &models.User{})

	setupTestData(db)
	router := setupRouter(db)

	return router
}

/*
Authorizeが不要な関数のテスト
*/
func TestFindAll(t *testing.T) {
	// setup for test
	router := setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items", nil) // 第3引数はrequestのbodyだがFindAll()では不要のためnil

	// APIrequestの実行
	router.ServeHTTP(w, req)

	// APIの実行結果を取得
	var res map[string][]models.Item
	json.Unmarshal([]byte(w.Body.String()), &res)

	// Assertion
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 3, len(res["data"]))
}

/*
Authorizationヘッダーがある場合のテスト
*/
func TestCreate(t *testing.T) {
	// setup for test
	router := setup()

	token, err := services.CreateToken(1, "test1@example.com")
	assert.Equal(t, nil, err)

	createItemInput := dto.CreateItemInput{
		Name:        "testItem4",
		Price:       4000,
		Description: "Create test",
	}

	reqBody, _ := json.Marshal(createItemInput) // convert struct to json

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+*token)

	// APIrequestの実行
	router.ServeHTTP(w, req)

	// APIの実行結果を取得
	var res map[string]models.Item
	json.Unmarshal([]byte(w.Body.String()), &res)

	// Assertion
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, uint(4), res["data"].ID)
}

/*
Authorizationヘッダーがない場合のテスト
goでは異常系のテストも正常系と同様に行う
*/
func TestCreateUnauthorized(t *testing.T) {
	// setup for test
	router := setup()

	createItemInput := dto.CreateItemInput{
		Name:        "testItem4",
		Price:       4000,
		Description: "Create test",
	}

	reqBody, _ := json.Marshal(createItemInput) // convert struct to json

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))

	// APIrequestの実行
	router.ServeHTTP(w, req)

	// APIの実行結果を取得
	var res map[string]models.Item
	json.Unmarshal([]byte(w.Body.String()), &res)

	// Assertion
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
