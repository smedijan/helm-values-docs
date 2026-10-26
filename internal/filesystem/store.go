// Package filesystem reads templates and persists generated documents.
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

// Read reads a generated document from disk.
func Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- path is explicitly provided by the caller.
	if err != nil {
		return nil, fmt.Errorf("read output %q: %w", path, err)
	}
	return data, nil
}

// WriteAtomic replaces a generated document after writing it successfully.
func WriteAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	file, err := os.CreateTemp(dir, ".hvd-*")
	if err != nil {
		return fmt.Errorf("create temporary output for %q: %w", path, err)
	}
	temporary := file.Name()
	defer func() { _ = os.Remove(temporary) }()
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return fmt.Errorf("write output %q: %w", path, err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return fmt.Errorf("sync output %q: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close output %q: %w", path, err)
	}
	if err := os.Rename(temporary, path); err != nil {
		return fmt.Errorf("replace output %q: %w", path, err)
	}
	return nil
}
