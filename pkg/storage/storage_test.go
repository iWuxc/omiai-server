package storage

import (
	"context"
	"os"
	"strings"
	"testing"

	"omiai-server/pkg/storage/driver"
	"github.com/stretchr/testify/assert"
)

func TestLocalDriver(t *testing.T) {
	root := "./test_uploads"
	defer os.RemoveAll(root)

	d := driver.NewLocal(root, "http://localhost")
	ctx := context.Background()
	content := "test content"
	r := strings.NewReader(content)
	key := "test.txt"

	url, err := d.Put(ctx, key, r)
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost/test.txt", url)

	_, err = os.Stat(root + "/" + key)
	assert.NoError(t, err)

	err = d.Delete(ctx, key)
	assert.NoError(t, err)
}
