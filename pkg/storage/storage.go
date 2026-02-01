package storage

import (
	"context"
	"io"
)

// Driver defines the interface for storage drivers
type Driver interface {
	Put(ctx context.Context, key string, r io.Reader) (string, error)
	Delete(ctx context.Context, key string) error
}
