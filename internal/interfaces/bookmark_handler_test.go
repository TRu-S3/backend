package interfaces

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TRu-S3/backend/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupBookmarkTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&database.User{}, &database.Bookmark{})
	if err != nil {
		panic("failed to migrate database")
	}

	return db
}

func TestCreateBookmark(t *testing.T) {
	db := setupBookmarkTestDB()
	handler := NewBookmarkHandler(db)

	// Create test users
	user1 := database.User{Gmail: "user1@example.com", Name: "User 1"}
	user2 := database.User{Gmail: "user2@example.com", Name: "User 2"}
	db.Create(&user1)
	db.Create(&user2)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/bookmarks", handler.CreateBookmark)

	t.Run("Valid bookmark creation", func(t *testing.T) {
		reqBody := CreateBookmarkRequest{
			UserID:           user1.ID,
			BookmarkedUserID: user2.ID,
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/bookmarks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var bookmark database.Bookmark
		json.Unmarshal(resp.Body.Bytes(), &bookmark)
		assert.Equal(t, user1.ID, bookmark.UserID)
		assert.Equal(t, user2.ID, bookmark.BookmarkedUserID)
	})

	t.Run("Cannot bookmark yourself", func(t *testing.T) {
		reqBody := CreateBookmarkRequest{
			UserID:           user1.ID,
			BookmarkedUserID: user1.ID,
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/bookmarks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Missing required fields", func(t *testing.T) {
		reqBody := CreateBookmarkRequest{
			UserID: user1.ID,
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/bookmarks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestListBookmarks(t *testing.T) {
	db := setupBookmarkTestDB()
	handler := NewBookmarkHandler(db)

	// Create test users
	user1 := database.User{Gmail: "user1@example.com", Name: "User 1"}
	user2 := database.User{Gmail: "user2@example.com", Name: "User 2"}
	user3 := database.User{Gmail: "user3@example.com", Name: "User 3"}
	db.Create(&user1)
	db.Create(&user2)
	db.Create(&user3)

	// Create test bookmarks
	bookmark1 := database.Bookmark{UserID: user1.ID, BookmarkedUserID: user2.ID}
	bookmark2 := database.Bookmark{UserID: user1.ID, BookmarkedUserID: user3.ID}
	db.Create(&bookmark1)
	db.Create(&bookmark2)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/bookmarks", handler.ListBookmarks)

	t.Run("List all bookmarks", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/bookmarks", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &response)
		bookmarks := response["bookmarks"].([]interface{})
		assert.Equal(t, 2, len(bookmarks))
	})

	t.Run("Filter by user_id", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/bookmarks?user_id=1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestDeleteBookmark(t *testing.T) {
	db := setupBookmarkTestDB()
	handler := NewBookmarkHandler(db)

	// Create test users
	user1 := database.User{Gmail: "user1@example.com", Name: "User 1"}
	user2 := database.User{Gmail: "user2@example.com", Name: "User 2"}
	db.Create(&user1)
	db.Create(&user2)

	// Create test bookmark
	bookmark := database.Bookmark{UserID: user1.ID, BookmarkedUserID: user2.ID}
	db.Create(&bookmark)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/bookmarks/:id", handler.DeleteBookmark)

	t.Run("Delete existing bookmark", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/bookmarks/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var count int64
		db.Model(&database.Bookmark{}).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Delete non-existent bookmark", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/bookmarks/999", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}