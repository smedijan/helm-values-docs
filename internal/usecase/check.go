package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/smedijan/helm-values-docs/internal/filesystem"
)

// ErrOutputMissing marks a check request whose generated document does not exist.
var ErrOutputMissing = errors.New("output document is missing")

// Check compares generated documentation without mutating the output document.
func Check(ctx context.Context, request Request) (bool, error) {
	want, err := render(ctx, request)
	if err != nil {
		return false, err
	}
	have, err := filesystem.Read(request.OutputPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, fmt.Errorf("%w: %w", ErrOutputMissing, err)
		}
		return false, err
	}
	return bytes.Equal(want, have), nil
}
