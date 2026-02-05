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

func TestClientRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClientRepo(db)
	ctx := context.Background()

	// 1. Create a client
	client := &biz_omiai.Client{
		Name: "Original Name",
		Age:  25,
	}
	err := repo.Create(ctx, client)
	assert.NoError(t, err)
	assert.NotZero(t, client.ID)

	// 2. Update the client
	client.Name = "Updated Name"
	client.Age = 26
	err = repo.Update(ctx, client)
	assert.NoError(t, err)

	// 3. Verify update
	updatedClient, err := repo.Get(ctx, client.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedClient.Name)
	assert.Equal(t, 26, updatedClient.Age)
}
