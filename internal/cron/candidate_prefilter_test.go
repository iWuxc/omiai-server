package cron

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"

	logger "github.com/iWuxc/go-wit/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *data.DB {
	// Setup Logger
	log = logger.NewLogger("test", func(l *logrus.Logger) {
		l.SetOutput(os.Stdout)
	})

	// Use SQLite in-memory DB for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate schemas
	err = db.AutoMigrate(&biz_omiai.Client{})
	assert.NoError(t, err)

	return &data.DB{DB: db}
}

func TestCandidatePreFilterService_Run(t *testing.T) {
	db := setupTestDB(t)
	service := NewCandidatePreFilterService(db)

	// Seed data
	// Client A: Male, 30, Bachelor
	clientA := &biz_omiai.Client{
		Name:      "Client A",
		Gender:    1,
		Age:       30,
		Education: 3,
		Status:    biz_omiai.ClientStatusSingle,
	}
	// Match B: Female, 28, Bachelor (Should match)
	matchB := &biz_omiai.Client{
		Name:      "Match B",
		Gender:    2,
		Age:       28,
		Education: 3,
		Status:    biz_omiai.ClientStatusSingle,
	}
	// Match C: Female, 40, Master (Should match but lower score/tags)
	matchC := &biz_omiai.Client{
		Name:      "Match C",
		Gender:    2,
		Age:       40,
		Education: 4,
		Status:    biz_omiai.ClientStatusSingle,
	}
	// Match D: Male (Should not match)
	matchD := &biz_omiai.Client{
		Name:   "Match D",
		Gender: 1,
		Age:    25,
		Status: biz_omiai.ClientStatusSingle,
	}
	// Match E: Female, Age 0 but Birthday set (Should calculate age and match)
	// Current date in env is 2026-02-05. Birthday 1996-02 => 30 years old.
	matchE := &biz_omiai.Client{
		Name:      "Match E",
		Gender:    2,
		Age:       0,
		Birthday:  "1996-02",
		Education: 3,
		Status:    biz_omiai.ClientStatusSingle,
	}

	db.Create(clientA)
	db.Create(matchB)
	db.Create(matchC)
	db.Create(matchD)
	db.Create(matchE)

	// Run logic
	ctx := context.Background()
	service.Execute(ctx)

	// Verify results
	var updatedClientA biz_omiai.Client
	err := db.First(&updatedClientA, clientA.ID).Error
	assert.NoError(t, err)

	assert.NotEmpty(t, updatedClientA.CandidateCacheJSON)

	var candidates []*biz_omiai.Candidate
	err = json.Unmarshal([]byte(updatedClientA.CandidateCacheJSON), &candidates)
	assert.NoError(t, err)

	// Expect Match B, C, E
	assert.GreaterOrEqual(t, len(candidates), 3)

	foundB := false
	foundC := false
	foundD := false
	foundE := false

	for _, c := range candidates {
		if c.CandidateID == matchB.ID {
			foundB = true
			// Check Tags
			assert.Contains(t, c.Tags, "年龄相仿")
			assert.Contains(t, c.Tags, "学历相当")
		}
		if c.CandidateID == matchE.ID {
			foundE = true
			// Age should be calculated. 1996-02 vs 2026-02 => 30.
			// clientA is 30. Diff is 0. Should have "年龄相仿".
			assert.Equal(t, 30, c.Age)
			assert.Contains(t, c.Tags, "年龄相仿")
		}
		if c.CandidateID == matchC.ID {
			foundC = true
			// Check Tags (Age diff 10 > 3)
			assert.NotContains(t, c.Tags, "年龄相仿")
			// Edu 4 >= 3, so it matches
			assert.Contains(t, c.Tags, "学历相当")
			// Has tags, so no "缘分推荐"
			assert.NotContains(t, c.Tags, "缘分推荐")
		}
		if c.CandidateID == matchD.ID {
			foundD = true
		}
	}

	assert.True(t, foundB, "Should find Match B")
	assert.True(t, foundC, "Should find Match C")
	assert.False(t, foundD, "Should NOT find Match D (same gender)")
	assert.True(t, foundE, "Should find Match E (age calculated from birthday)")
}
