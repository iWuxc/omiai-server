package omiai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMatchTestDB(t *testing.T) *data.DB {
	// Use SQLite in-memory DB with unique name per test to ensure isolation
	// t.Name() might contain characters like '/', so we sanitize it.
	dbName := strings.ReplaceAll(t.Name(), "/", "_")
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_busy_timeout=5000", dbName)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Clean up previous tables if any (due to shared cache with same name)
	db.Migrator().DropTable(&biz_omiai.Client{}, &biz_omiai.MatchRecord{}, &biz_omiai.MatchStatusHistory{}, &biz_omiai.FollowUpRecord{})

	// Migrate schemas
	err = db.AutoMigrate(&biz_omiai.Client{}, &biz_omiai.MatchRecord{}, &biz_omiai.MatchStatusHistory{}, &biz_omiai.FollowUpRecord{})
	assert.NoError(t, err)

	return &data.DB{DB: db}
}

func TestMatchRepo_Compare(t *testing.T) {
	db := setupMatchTestDB(t)
	repo := NewMatchRepo(db)

	// Seed data
	c1 := &biz_omiai.Client{
		Name:      "Client 1",
		Gender:    1,
		Age:       30,
		Height:    180,
		Education: 3,
		Income:    20000,
	}
	c2 := &biz_omiai.Client{
		Name:      "Client 2",
		Gender:    2,
		Age:       28,
		Height:    165,
		Education: 3,
		Income:    15000,
	}
	db.Create(c1)
	db.Create(c2)

	// Test Compare
	ctx := context.Background()
	comp, err := repo.Compare(ctx, c1.ID, c2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, comp)

	// Verify JSON serialization
	bytes, err := json.Marshal(comp)
	assert.NoError(t, err)
	t.Logf("Compare Result JSON: %s", string(bytes))

	// Verify fields
	assert.NotNil(t, comp.BasicInfo)
	assert.NotNil(t, comp.PersonalityRadar)
	assert.NotNil(t, comp.Interests)
	assert.NotNil(t, comp.Values)
	assert.NotNil(t, comp.RelationshipExpectations)

	// Verify specific values
	assert.Equal(t, 30, comp.BasicInfo["age"]["client"])
	assert.Equal(t, 28, comp.BasicInfo["age"]["candidate"])

	// Verify PersonalityRadar keys
	assert.Contains(t, comp.PersonalityRadar, "openness")
	assert.Contains(t, comp.PersonalityRadar, "conscientiousness")
}

func TestMatchRepo_GetCandidates(t *testing.T) {
	db := setupMatchTestDB(t)
	repo := NewMatchRepo(db)

	// Seed data
	c1 := &biz_omiai.Client{
		Name:      "Client 1",
		Gender:    1, // Male
		Age:       30,
		Education: 3,
		Status:    biz_omiai.ClientStatusSingle,
	}
	// Candidate (Female)
	c2 := &biz_omiai.Client{
		Name:      "Candidate 1",
		Gender:    2, // Female
		Age:       28,
		Education: 3,
		Status:    biz_omiai.ClientStatusSingle,
	}

	db.Create(c1)
	db.Create(c2)

	ctx := context.Background()

	// 1. Test Fallback (No Cache)
	candidates, err := repo.GetCandidates(ctx, c1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, candidates)
	assert.Equal(t, c2.ID, candidates[0].CandidateID)
	assert.NotEmpty(t, candidates[0].Tags)

	// 2. Test Cache
	// Create dummy cache
	cacheData := []*biz_omiai.Candidate{
		{
			CandidateID: c2.ID,
			Name:        c2.Name,
			MatchScore:  99,
			Tags:        []string{"Cached Tag"},
		},
	}
	bytes, _ := json.Marshal(cacheData)
	c1.CandidateCacheJSON = string(bytes)
	db.Save(c1)

	candidatesCached, err := repo.GetCandidates(ctx, c1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, candidatesCached)
	assert.Equal(t, 99, candidatesCached[0].MatchScore)
	assert.Equal(t, "Cached Tag", candidatesCached[0].Tags[0])
}

func TestMatchRepo_CreateFollowUp(t *testing.T) {
	db := setupMatchTestDB(t)
	repo := NewMatchRepo(db)
	ctx := context.Background()

	record := &biz_omiai.FollowUpRecord{
		MatchRecordID: 1,
		Method:        "电话",
		Content:       "Test Content",
		Feedback:      "Test Feedback",
		Satisfaction:  5,
		Attachments:   `["file1.jpg"]`,
	}

	err := repo.CreateFollowUp(ctx, record)
	assert.NoError(t, err)
	assert.NotZero(t, record.ID)

	// Verify persistence
	var saved biz_omiai.FollowUpRecord
	err = db.DB.First(&saved, record.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Test Feedback", saved.Feedback)
	assert.Equal(t, int8(5), saved.Satisfaction)
	assert.Equal(t, `["file1.jpg"]`, saved.Attachments)
}
