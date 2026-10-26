package main

import (
	"os"

	"github.com/smedijan/helm-values-docs/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
