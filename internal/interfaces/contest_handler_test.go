package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TRu-S3/backend/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestHandler(t *testing.T) (*ContestHandler, *gorm.DB) {
	// Create in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(&database.Contest{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	handler := NewContestHandler(db)
	return handler, db
}

func setupTestRouter(handler *ContestHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	v1 := r.Group("/api/v1")
	contests := v1.Group("/contests")
	{
		contests.POST("", handler.CreateContest)
		contests.GET("", handler.ListContests)
		contests.GET("/:id", handler.GetContest)
		contests.PUT("/:id", handler.UpdateContest)
		contests.DELETE("/:id", handler.DeleteContest)
	}
	
	return r
}

func TestContestHandler_CreateContest(t *testing.T) {
	handler, _ := setupTestHandler(t)
	router := setupTestRouter(handler)

	t.Run("Valid Contest Creation", func(t *testing.T) {
		reqBody := CreateContestRequest{
			BackendQuota:        2,
			FrontendQuota:       3,
			AIQuota:             1,
			ApplicationDeadline: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			Purpose:             "テストコンテスト",
			Message:             "これはテスト用のメッセージです",
			AuthorID:            1,
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/contests", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response database.Contest
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, reqBody.BackendQuota, response.BackendQuota)
		assert.Equal(t, reqBody.Purpose, response.Purpose)
		assert.NotZero(t, response.ID)
	})

	t.Run("Invalid Request - Missing Required Fields", func(t *testing.T) {
		reqBody := CreateContestRequest{
			BackendQuota: 2,
			// Missing required fields
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/contests", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Deadline Format", func(t *testing.T) {
		reqBody := CreateContestRequest{
			BackendQuota:        2,
			FrontendQuota:       3,
			AIQuota:             1,
			ApplicationDeadline: "invalid-date",
			Purpose:             "テストコンテスト",
			Message:             "これはテスト用のメッセージです",
			AuthorID:            1,
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/contests", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestContestHandler_GetContest(t *testing.T) {
	handler, db := setupTestHandler(t)
	router := setupTestRouter(handler)

	// Create a test contest
	contest := database.Contest{
		BackendQuota:        2,
		FrontendQuota:       3,
		AIQuota:             1,
		ApplicationDeadline: time.Now().Add(24 * time.Hour),
		Purpose:             "テストコンテスト",
		Message:             "これはテスト用のメッセージです",
		AuthorID:            1,
	}
	db.Create(&contest)

	t.Run("Get Existing Contest", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/contests/%d", contest.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response database.Contest
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, contest.ID, response.ID)
		assert.Equal(t, contest.Purpose, response.Purpose)
	})

	t.Run("Get Non-existent Contest", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/contests/99999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestContestHandler_ListContests(t *testing.T) {
	handler, db := setupTestHandler(t)
	router := setupTestRouter(handler)

	// Create test contests
	contests := []database.Contest{
		{
			BackendQuota:        1,
			FrontendQuota:       2,
			AIQuota:             0,
			ApplicationDeadline: time.Now().Add(24 * time.Hour),
			Purpose:             "コンテスト1",
			Message:             "メッセージ1",
			AuthorID:            1,
		},
		{
			BackendQuota:        2,
			FrontendQuota:       1,
			AIQuota:             1,
			ApplicationDeadline: time.Now().Add(48 * time.Hour),
			Purpose:             "コンテスト2",
			Message:             "メッセージ2",
			AuthorID:            2,
		},
		{
			BackendQuota:        1,
			FrontendQuota:       1,
			AIQuota:             1,
			ApplicationDeadline: time.Now().Add(-24 * time.Hour), // Past deadline
			Purpose:             "期限切れコンテスト",
			Message:             "期限切れメッセージ",
			AuthorID:            1,
		},
	}

	for _, contest := range contests {
		db.Create(&contest)
	}

	t.Run("List All Contests", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/contests", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		contests := response["contests"].([]interface{})
		assert.Len(t, contests, 3)

		pagination := response["pagination"].(map[string]interface{})
		assert.Equal(t, float64(1), pagination["page"])
		assert.Equal(t, float64(10), pagination["limit"])
		assert.Equal(t, float64(3), pagination["total"])
	})

	t.Run("Filter by Author", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/contests?author_id=1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		contests := response["contests"].([]interface{})
		assert.Len(t, contests, 2) // Two contests by author 1
	})

	t.Run("Filter Active Contests", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/contests?active=true", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		contests := response["contests"].([]interface{})
		assert.Len(t, contests, 2) // Two active contests
	})

	t.Run("Pagination", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/contests?page=1&limit=2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		contests := response["contests"].([]interface{})
		assert.Len(t, contests, 2)

		pagination := response["pagination"].(map[string]interface{})
		assert.Equal(t, float64(1), pagination["page"])
		assert.Equal(t, float64(2), pagination["limit"])
		assert.Equal(t, float64(3), pagination["total"])
	})
}

func TestContestHandler_UpdateContest(t *testing.T) {
	handler, db := setupTestHandler(t)
	router := setupTestRouter(handler)

	// Create a test contest
	contest := database.Contest{
		BackendQuota:        2,
		FrontendQuota:       3,
		AIQuota:             1,
		ApplicationDeadline: time.Now().Add(24 * time.Hour),
		Purpose:             "テストコンテスト",
		Message:             "これはテスト用のメッセージです",
		AuthorID:            1,
	}
	db.Create(&contest)

	t.Run("Update Contest", func(t *testing.T) {
		newPurpose := "更新されたコンテスト"
		newBackendQuota := 5
		reqBody := UpdateContestRequest{
			Purpose:      &newPurpose,
			BackendQuota: &newBackendQuota,
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/contests/%d", contest.ID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response database.Contest
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, newPurpose, response.Purpose)
		assert.Equal(t, newBackendQuota, response.BackendQuota)
		assert.Equal(t, contest.FrontendQuota, response.FrontendQuota) // Unchanged
	})

	t.Run("Update Non-existent Contest", func(t *testing.T) {
		newPurpose := "更新されたコンテスト"
		reqBody := UpdateContestRequest{
			Purpose: &newPurpose,
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/api/v1/contests/99999", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestContestHandler_DeleteContest(t *testing.T) {
	handler, db := setupTestHandler(t)
	router := setupTestRouter(handler)

	// Create a test contest
	contest := database.Contest{
		BackendQuota:        2,
		FrontendQuota:       3,
		AIQuota:             1,
		ApplicationDeadline: time.Now().Add(24 * time.Hour),
		Purpose:             "削除テストコンテスト",
		Message:             "これは削除テスト用のメッセージです",
		AuthorID:            1,
	}
	db.Create(&contest)

	t.Run("Delete Existing Contest", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/contests/%d", contest.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify deletion
		var deletedContest database.Contest
		err := db.First(&deletedContest, contest.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Delete Non-existent Contest", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/contests/99999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}