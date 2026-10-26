package schema

import (
	"fmt"
	"os"

	"github.com/smedijan/helm-values-docs/internal/domain"
)

// Load reads a schema file and decodes it into the domain model.
func Load(path string) (*domain.Schema, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- path is provided explicitly by the caller.
	if err != nil {
		return nil, fmt.Errorf("read schema %q: %w", path, err)
	}
	model, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("schema %q: %w", path, err)
	}
	return model, nil
}
