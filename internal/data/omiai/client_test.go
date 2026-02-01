package omiai

import (
	"context"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *data.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	err = db.AutoMigrate(&biz_omiai.Client{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return &data.DB{DB: db}
}

func TestClientRepo_Stats(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClientRepo(db)
	ctx := context.Background()

	// 1. Empty Stats
	stats, err := repo.Stats(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), stats["total"])
	assert.Equal(t, int64(0), stats["today"])

	// 2. Add some data
	now := time.Now()
	clients := []*biz_omiai.Client{
		{Name: "Test 1", Gender: 1, CreatedAt: now},
		{Name: "Test 2", Gender: 2, CreatedAt: now},
		{Name: "Old 1", Gender: 1, CreatedAt: now.AddDate(0, 0, -2)},
	}

	for _, c := range clients {
		err := repo.Create(ctx, c)
		assert.NoError(t, err)
	}

	// 3. Verify Stats
	stats, err = repo.Stats(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), stats["total"])
	assert.Equal(t, int64(2), stats["today"])
}
