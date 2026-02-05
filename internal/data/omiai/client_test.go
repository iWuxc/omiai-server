package omiai

import (
	"context"
	"fmt"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupClientTestDB(t *testing.T) *data.DB {
	dbName := strings.ReplaceAll(t.Name(), "/", "_")
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_busy_timeout=5000", dbName)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	db.Migrator().DropTable(&biz_omiai.Client{})
	err = db.AutoMigrate(&biz_omiai.Client{})
	assert.NoError(t, err)

	return &data.DB{DB: db}
}

func TestClientRepo_Get_WithPartner(t *testing.T) {
	db := setupClientTestDB(t)
	repo := NewClientRepo(db)
	ctx := context.Background()

	// 1. Create Partner
	partner := &biz_omiai.Client{
		Name:   "Partner Client",
		Gender: 2,
		Status: biz_omiai.ClientStatusMatched,
	}
	err := repo.Create(ctx, partner)
	assert.NoError(t, err)

	// 2. Create Client linked to Partner
	partnerID := partner.ID
	client := &biz_omiai.Client{
		Name:      "Main Client",
		Gender:    1,
		Status:    biz_omiai.ClientStatusMatched,
		PartnerID: &partnerID,
	}
	err = repo.Create(ctx, client)
	assert.NoError(t, err)

	// 3. Get Client and verify Partner is loaded
	gotClient, err := repo.Get(ctx, client.ID)
	assert.NoError(t, err)
	assert.NotNil(t, gotClient)
	assert.NotNil(t, gotClient.Partner)
	assert.Equal(t, partner.ID, gotClient.Partner.ID)
	assert.Equal(t, partner.Name, gotClient.Partner.Name)

	// 4. Verify PartnerID is correct
	assert.NotNil(t, gotClient.PartnerID)
	assert.Equal(t, partnerID, *gotClient.PartnerID)
}

func TestClientRepo_Get_NoPartner(t *testing.T) {
	db := setupClientTestDB(t)
	repo := NewClientRepo(db)
	ctx := context.Background()

	// Create Client without Partner
	client := &biz_omiai.Client{
		Name:   "Single Client",
		Gender: 1,
		Status: biz_omiai.ClientStatusSingle,
	}
	err := repo.Create(ctx, client)
	assert.NoError(t, err)

	// Get Client
	gotClient, err := repo.Get(ctx, client.ID)
	assert.NoError(t, err)
	assert.NotNil(t, gotClient)
	assert.Nil(t, gotClient.Partner)
	assert.Nil(t, gotClient.PartnerID)
}
