package omiai

import (
	"context"
	"testing"

	biz_omiai "omiai-server/internal/biz/omiai"

	"github.com/stretchr/testify/assert"
)

func TestMatchRepo_UniqueConstraint(t *testing.T) {
	db := setupMatchTestDB(t)
	repo := &MatchRepo{db: db}

	// Seed Clients
	// Client A (Male)
	c1 := &biz_omiai.Client{Name: "Client A", Gender: 1, Status: biz_omiai.ClientStatusSingle}
	// Client B (Female) - The target
	c2 := &biz_omiai.Client{Name: "Client B", Gender: 2, Status: biz_omiai.ClientStatusSingle}
	// Client C (Male) - The intruder
	c3 := &biz_omiai.Client{Name: "Client C", Gender: 1, Status: biz_omiai.ClientStatusSingle}

	db.Create(c1)
	db.Create(c2)
	db.Create(c3)

	ctx := context.Background()
	adminID := "admin_tester"

	// 1. Match A and B (Success)
	matchRecord, err := repo.confirmMatchDB(ctx, c1.ID, c2.ID, adminID, "")
	assert.NoError(t, err)
	assert.NotNil(t, matchRecord)

	// Verify Status and PartnerID
	var clientA, clientB biz_omiai.Client
	db.First(&clientA, c1.ID)
	db.First(&clientB, c2.ID)
	assert.Equal(t, int8(biz_omiai.ClientStatusMatched), clientA.Status)
	assert.NotNil(t, clientA.PartnerID)
	assert.Equal(t, c2.ID, *clientA.PartnerID)
	assert.Equal(t, int8(biz_omiai.ClientStatusMatched), clientB.Status)
	assert.NotNil(t, clientB.PartnerID)
	assert.Equal(t, c1.ID, *clientB.PartnerID)

	// 2. Try to Match C and B (Should Fail due to Status Check)
	// Because confirmMatchDB checks status first.
	_, err = repo.confirmMatchDB(ctx, c3.ID, c2.ID, adminID, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already matched")

	// 3. Try to Bypass Status Check and force DB Constraint Violation
	// We manually update C to point to B to test Unique Index
	// A.PartnerID is B.ID (2)
	// We try to set C.PartnerID = 2
	err = db.Model(&biz_omiai.Client{}).Where("id = ?", c3.ID).
		Update("partner_id", c2.ID).Error

	// Should fail due to UNIQUE constraint
	assert.Error(t, err)
	// SQLite error message for unique constraint
	// "UNIQUE constraint failed: client.partner_id"
	assert.Contains(t, err.Error(), "UNIQUE constraint failed")
}

func TestMatchRepo_ConfirmMatchConcurrency(t *testing.T) {
	db := setupMatchTestDB(t)
	repo := &MatchRepo{db: db}

	// Setup: Client A and B
	c1 := &biz_omiai.Client{Name: "Client A", Gender: 1, Status: biz_omiai.ClientStatusSingle}
	c2 := &biz_omiai.Client{Name: "Client B", Gender: 2, Status: biz_omiai.ClientStatusSingle}
	db.Create(c1)
	db.Create(c2)

	// Concurrency
	concurrency := 3 // Reduced from 10 to avoid SQLite excessive locking
	errChan := make(chan error, concurrency)
	successCount := 0

	ctx := context.Background()

	for i := 0; i < concurrency; i++ {
		go func(idx int) {
			// Call confirmMatchDB directly (bypassing Redis lock to test DB lock)
			_, err := repo.confirmMatchDB(ctx, c1.ID, c2.ID, "admin", "")
			errChan <- err
		}(i)
	}

	for i := 0; i < concurrency; i++ {
		err := <-errChan
		if err == nil {
			successCount++
		} else {
			// t.Logf("Goroutine error: %v", err)
		}
	}

	// In SQLite with high concurrency, it's possible all fail with "database table is locked" or exactly one succeeds.
	// For the purpose of "Anti-duplicate", we just need to ensure successCount <= 1.
	assert.LessOrEqual(t, successCount, 1, "At most one match confirmation should succeed")

	if successCount == 1 {
		// Verify DB state only if one succeeded
		var clientA biz_omiai.Client
		db.First(&clientA, c1.ID)
		assert.Equal(t, int8(biz_omiai.ClientStatusMatched), clientA.Status)
		assert.NotNil(t, clientA.PartnerID)
		if clientA.PartnerID != nil {
			assert.Equal(t, c2.ID, *clientA.PartnerID)
		}

		// Verify only 1 match record exists
		var matchCount int64
		db.Model(&biz_omiai.MatchRecord{}).Count(&matchCount)
		assert.Equal(t, int64(1), matchCount)
	} else {
		t.Log("All attempts failed due to DB locking (expected in SQLite concurrency)")
	}
}
