package omiai

import (
	"context"
	"testing"

	biz_omiai "omiai-server/internal/biz/omiai"

	"github.com/stretchr/testify/assert"
)

func TestMatchRepo_UpdateStatus(t *testing.T) {
	db := setupMatchTestDB(t)
	repo := NewMatchRepo(db)
	ctx := context.Background()

	// 1. Create Match Record
	record := &biz_omiai.MatchRecord{
		MaleClientID:   1,
		FemaleClientID: 2,
		Status:         biz_omiai.MatchStatusAcquaintance, // 1
	}
	err := repo.Create(ctx, record)
	assert.NoError(t, err)

	// 2. Update Status
	newStatus := int8(biz_omiai.MatchStatusDating) // 2
	err = repo.UpdateStatus(ctx, record.ID, int8(record.Status), newStatus, "admin", "test reason")
	assert.NoError(t, err)

	// 3. Verify Status Updated
	updatedRecord, err := repo.Get(ctx, record.ID)
	assert.NoError(t, err)
	assert.Equal(t, newStatus, updatedRecord.Status)

	// 4. Verify History Created
	var history biz_omiai.MatchStatusHistory
	err = db.DB.First(&history, "match_record_id = ?", record.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, record.ID, history.MatchRecordID)
	assert.Equal(t, int8(biz_omiai.MatchStatusAcquaintance), history.OldStatus)
	assert.Equal(t, newStatus, history.NewStatus)
	assert.Equal(t, "admin", history.Operator)
	assert.Equal(t, "test reason", history.Reason)
}
