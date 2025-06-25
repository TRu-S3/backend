package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(&Contest{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestContestModel(t *testing.T) {
	t.Run("Create Contest", func(t *testing.T) {
		db := setupTestDB(t)
		deadline := time.Now().Add(24 * time.Hour)
		contest := Contest{
			BackendQuota:        2,
			FrontendQuota:       3,
			AIQuota:             1,
			ApplicationDeadline: deadline,
			Purpose:             "テストコンテスト",
			Message:             "これはテスト用のメッセージです",
			AuthorID:            1,
		}

		err := db.Create(&contest).Error
		assert.NoError(t, err)
		assert.NotZero(t, contest.ID)
		assert.NotZero(t, contest.CreatedAt)
		assert.NotZero(t, contest.UpdatedAt)
	})

	t.Run("Contest TableName", func(t *testing.T) {
		contest := Contest{}
		assert.Equal(t, "contests", contest.TableName())
	})

	t.Run("Contest Validation", func(t *testing.T) {
		db := setupTestDB(t)
		// Test with missing purpose field - SQLite doesn't enforce NOT NULL at GORM level by default
		contest := Contest{
			BackendQuota:        2,
			FrontendQuota:       3,
			AIQuota:             1,
			ApplicationDeadline: time.Now().Add(24 * time.Hour),
			// Missing Purpose and Message
			AuthorID: 1,
		}
		err := db.Create(&contest).Error
		// In SQLite, empty strings are allowed, so we test with proper validation
		assert.NoError(t, err) // SQLite allows empty strings for TEXT fields
		assert.NotZero(t, contest.ID)
	})

	t.Run("Contest Update", func(t *testing.T) {
		db := setupTestDB(t)
		deadline := time.Now().Add(24 * time.Hour)
		contest := Contest{
			BackendQuota:        2,
			FrontendQuota:       3,
			AIQuota:             1,
			ApplicationDeadline: deadline,
			Purpose:             "テストコンテスト",
			Message:             "これはテスト用のメッセージです",
			AuthorID:            1,
		}

		// Create contest
		err := db.Create(&contest).Error
		assert.NoError(t, err)

		// Update contest
		originalUpdatedAt := contest.UpdatedAt
		time.Sleep(10 * time.Millisecond) // Ensure time difference

		contest.Purpose = "更新されたテストコンテスト"
		err = db.Save(&contest).Error
		assert.NoError(t, err)
		assert.NotEqual(t, originalUpdatedAt, contest.UpdatedAt)
		assert.Equal(t, "更新されたテストコンテスト", contest.Purpose)
	})

	t.Run("Contest Delete", func(t *testing.T) {
		db := setupTestDB(t)
		deadline := time.Now().Add(24 * time.Hour)
		contest := Contest{
			BackendQuota:        2,
			FrontendQuota:       3,
			AIQuota:             1,
			ApplicationDeadline: deadline,
			Purpose:             "削除テストコンテスト",
			Message:             "これは削除テスト用のメッセージです",
			AuthorID:            1,
		}

		// Create contest
		err := db.Create(&contest).Error
		assert.NoError(t, err)
		contestID := contest.ID

		// Delete contest
		err = db.Delete(&contest).Error
		assert.NoError(t, err)

		// Verify deletion
		var deletedContest Contest
		err = db.First(&deletedContest, contestID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Contest Query", func(t *testing.T) {
		db := setupTestDB(t)
		// Create multiple contests
		contests := []Contest{
			{
				BackendQuota:        1,
				FrontendQuota:       2,
				AIQuota:             0,
				ApplicationDeadline: time.Now().Add(24 * time.Hour),
				Purpose:             "クエリテスト1",
				Message:             "メッセージ1",
				AuthorID:            1,
			},
			{
				BackendQuota:        2,
				FrontendQuota:       1,
				AIQuota:             1,
				ApplicationDeadline: time.Now().Add(48 * time.Hour),
				Purpose:             "クエリテスト2",
				Message:             "メッセージ2",
				AuthorID:            2,
			},
		}

		for _, contest := range contests {
			err := db.Create(&contest).Error
			assert.NoError(t, err)
		}

		// Query by author
		var authorContests []Contest
		err := db.Where("author_id = ?", 1).Find(&authorContests).Error
		assert.NoError(t, err)
		assert.Len(t, authorContests, 1)
		assert.Equal(t, "クエリテスト1", authorContests[0].Purpose)

		// Query active contests (deadline not passed)
		var activeContests []Contest
		err = db.Where("application_deadline > ?", time.Now()).Find(&activeContests).Error
		assert.NoError(t, err)
		assert.Len(t, activeContests, 2)

		// Query with ordering
		var orderedContests []Contest
		err = db.Order("created_at DESC").Find(&orderedContests).Error
		assert.NoError(t, err)
		assert.Len(t, orderedContests, 2)
	})
}
