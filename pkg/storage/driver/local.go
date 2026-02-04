package driver

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	Root string
	Url  string
}

func NewLocal(root, url string) *Local {
	return &Local{Root: root, Url: url}
}

func (l *Local) Put(ctx context.Context, key string, r io.Reader, contentType string) (string, error) {
	path := filepath.Join(l.Root, key)
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", l.Url, key), nil
}

func (l *Local) Delete(ctx context.Context, key string) error {
	return os.Remove(filepath.Join(l.Root, key))
}
