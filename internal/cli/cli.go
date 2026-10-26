// Package cli parses command-line requests and invokes their use cases.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/smedijan/helm-values-docs/internal/usecase"
)

type stringList []string

func (s *stringList) String() string { return strings.Join(*s, ",") }
func (s *stringList) Set(value string) error {
	if value == "" {
		return fmt.Errorf("value path must not be empty")
	}
	*s = append(*s, value)
	return nil
}

// Run executes the CLI and returns its process exit status.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return 2
	}
	switch args[0] {
	case "render":
		return render(args[1:], stderr)
	case "check":
		return check(args[1:], stderr)
	default:
		_, _ = fmt.Fprintf(stderr, "hvd: unknown command %q\n", args[0])
		usage(stderr)
		return 2
	}
}

func render(args []string, stderr io.Writer) int {
	request, ok := parseFlags("render", args, stderr)
	if !ok {
		return 2
	}
	if err := usecase.Render(context.Background(), request); err != nil {
		_, _ = fmt.Fprintf(stderr, "render: %v\n", err)
		return 2
	}
	return 0
}

func check(args []string, stderr io.Writer) int {
	request, ok := parseFlags("check", args, stderr)
	if !ok {
		return 2
	}
	current, err := usecase.Check(context.Background(), request)
	if err != nil {
		if errors.Is(err, usecase.ErrOutputMissing) {
			_, _ = fmt.Fprintf(stderr, "check: output %q is missing\n", request.OutputPath)
			return 1
		}
		_, _ = fmt.Fprintf(stderr, "check: %v\n", err)
		return 2
	}
	if !current {
		_, _ = fmt.Fprintf(stderr, "check: output %q is stale\n", request.OutputPath)
		return 1
	}
	return 0
}

func parseFlags(name string, args []string, stderr io.Writer) (usecase.Request, bool) {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(stderr)
	schemaPath := flags.String("schema", "", "path to values.schema.json")
	var values stringList
	flags.Var(&values, "values", "path to a values YAML file (repeatable)")
	outputPath := flags.String("output", "", "path for generated Markdown")
	templatePath := flags.String("template", "", "optional Markdown template path")
	if err := flags.Parse(args); err != nil {
		return usecase.Request{}, false
	}
	if *schemaPath == "" || len(values) == 0 || *outputPath == "" {
		_, _ = fmt.Fprintf(stderr, "%s: --schema, --values (at least once), and --output are required\n", name)
		flags.PrintDefaults()
		return usecase.Request{}, false
	}
	if flags.NArg() != 0 {
		_, _ = fmt.Fprintf(stderr, "%s: unexpected arguments: %s\n", name, strings.Join(flags.Args(), " "))
		return usecase.Request{}, false
	}
	return usecase.Request{SchemaPath: *schemaPath, ValuesPaths: values, OutputPath: *outputPath, TemplatePath: *templatePath}, true
}

func usage(stderr io.Writer) {
	_, _ = fmt.Fprintln(stderr, "usage: hvd {render|check} --schema SCHEMA --values VALUES [--values VALUES ...] --output OUTPUT")
}
