// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Command poweradmin is the entry point for the Poweradmin CLI.
// It resolves credentials from three sources in increasing order of precedence:
//
//  1. Config file    (~/.config/poweradmin/config.yaml)
//  2. Environment    (POWERADMIN_URL, POWERADMIN_API_KEY)
//  3. CLI flags      (--url, --api-key)
package main

import (
	"fmt"
	"os"

	"github.com/contentways/poweradmin-cli/v2/internal/cli"
	"github.com/contentways/poweradmin-cli/v2/internal/config"
	"github.com/contentways/poweradmin-cli/v2/internal/state"
)

func main() {
	// Load configuration from the default config file.
	// A missing file is not treated as an error — credentials can be
	// supplied via environment variables or CLI flags instead.
	cfg, err := config.Load(config.DefaultPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load config: %s\n", err)
		cfg = &config.Config{}
	}

	// Start with values from the config file.
	url := cfg.URL
	apiKey := cfg.APIKey

	// Environment variables override the config file.
	if v := os.Getenv("POWERADMIN_URL"); v != "" {
		url = v
	}
	if v := os.Getenv("POWERADMIN_API_KEY"); v != "" {
		apiKey = v
	}

	// CLI flags (--url, --api-key) override everything else.
	// This is handled in PersistentPreRunE inside cli.NewRootCommand,
	// after Cobra has parsed the flags.
	s := state.New(url, apiKey)

	// Build the complete command tree and hand control to Cobra.
	root := cli.NewRootCommand(s)
	cli.Execute(root)
}
