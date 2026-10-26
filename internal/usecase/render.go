package usecase

import (
	"context"

	"github.com/smedijan/helm-values-docs/internal/filesystem"
	"github.com/smedijan/helm-values-docs/internal/markdown"
	"github.com/smedijan/helm-values-docs/internal/schema"
	"github.com/smedijan/helm-values-docs/internal/values"
)

// Render generates and atomically persists chart values documentation.
func Render(ctx context.Context, request Request) error {
	data, err := render(ctx, request)
	if err != nil {
		return err
	}
	return filesystem.WriteAtomic(request.OutputPath, data)
}

func render(ctx context.Context, request Request) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	model, err := schema.Load(request.SchemaPath)
	if err != nil {
		return nil, err
	}
	effectiveValues, err := values.Load(request.ValuesPaths)
	if err != nil {
		return nil, err
	}
	data, err := markdown.Render(model, effectiveValues)
	if err != nil {
		return nil, err
	}
	if request.TemplatePath == "" {
		return data, nil
	}
	return filesystem.Compose(data, request.TemplatePath)
}
