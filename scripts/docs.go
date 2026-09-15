// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

//go:build ignore

// Command docs generates Markdown reference documentation for the poweradmin CLI.
// Run with: go run ./scripts/docs.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/contentways/poweradmin-cli/v3/internal/cli"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra/doc"
)

func main() {
	dir := "./docs/reference"

	// Ensure the output directory is empty and exists.
	if err := os.RemoveAll(dir); err != nil {
		log.Fatalf("could not remove directory: %v", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("could not create directory: %v", err)
	}

	// Build the command tree with empty credentials — docs only need the structure.
	s := state.New("", "")
	root := cli.NewRootCommand(s)

	// Disable the default completion command from appearing in the docs.
	root.CompletionOptions.DisableDefaultCmd = true

	if err := doc.GenMarkdownTree(root, dir); err != nil {
		log.Fatalf("could not generate docs: %v", err)
	}

	fmt.Printf("Documentation generated in %s\n", dir)
}
