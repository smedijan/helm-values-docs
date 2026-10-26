// Package usecase implements the chart documentation use cases.
package usecase

// Request identifies the files used to produce a document.
type Request struct {
	SchemaPath   string
	ValuesPaths  []string
	OutputPath   string
	TemplatePath string
}
