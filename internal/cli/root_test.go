// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package cli_test

import (
	"bytes"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cli"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
)

func TestNewRootCommand(t *testing.T) {
	s := state.New(
		"https://dns.example.com",
		"secret",
	)

	cmd := cli.NewRootCommand(s)

	if cmd.Use != "poweradmin" {
		t.Fatalf("Use = %q", cmd.Use)
	}

	if len(cmd.Commands()) != 6 {
		t.Fatalf("got %d commands, want 6", len(cmd.Commands()))
	}
}

// TestRootCommandPersistentFlagsExist verifies that --url, --api-key, and
// --verbose (with their short forms) are registered on the root command.
func TestRootCommandPersistentFlagsExist(t *testing.T) {
	s := state.New("https://dns.example.com", "secret")
	cmd := cli.NewRootCommand(s)

	for _, name := range []string{"url", "api-key", "verbose"} {
		if cmd.PersistentFlags().Lookup(name) == nil {
			t.Errorf("expected persistent flag --%s to be registered", name)
		}
	}

	shorthands := map[string]string{"u": "url", "k": "api-key", "v": "verbose"}
	for short, long := range shorthands {
		f := cmd.PersistentFlags().ShorthandLookup(short)
		if f == nil {
			t.Errorf("expected shorthand -%s to be registered", short)
			continue
		}
		if f.Name != long {
			t.Errorf("expected -%s to map to --%s, got --%s", short, long, f.Name)
		}
	}
}

// TestRootCommandURLFlagOverridesState verifies that --url takes precedence
// over the URL the State was constructed with (simulating config/env values).
func TestRootCommandURLFlagOverridesState(t *testing.T) {
	s := state.New("https://config.example.com", "secret")
	cmd := cli.NewRootCommand(s)
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--url", "https://flag.example.com", "version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.URL != "https://flag.example.com" {
		t.Errorf("expected URL to be overridden to https://flag.example.com, got %q", s.URL)
	}
}

// TestRootCommandAPIKeyFlagOverridesState verifies that --api-key takes
// precedence over the API key the State was constructed with.
func TestRootCommandAPIKeyFlagOverridesState(t *testing.T) {
	s := state.New("https://config.example.com", "config-key")
	cmd := cli.NewRootCommand(s)
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--api-key", "flag-key", "version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.APIKey != "flag-key" {
		t.Errorf("expected APIKey to be overridden to flag-key, got %q", s.APIKey)
	}
}

// TestRootCommandVerboseFlagEnablesState verifies that --verbose sets
// State.Verbose to true.
func TestRootCommandVerboseFlagEnablesState(t *testing.T) {
	s := state.New("https://dns.example.com", "secret")
	cmd := cli.NewRootCommand(s)
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--verbose", "version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Verbose {
		t.Error("expected Verbose to be true after --verbose flag")
	}
}

// TestRootCommandNoFlagsPreservesState verifies that omitting --url,
// --api-key, and --verbose leaves the State exactly as constructed —
// config/env-sourced values are not clobbered by absent flags.
func TestRootCommandNoFlagsPreservesState(t *testing.T) {
	s := state.New("https://config.example.com", "config-key")
	cmd := cli.NewRootCommand(s)
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.URL != "https://config.example.com" {
		t.Errorf("expected URL to remain https://config.example.com, got %q", s.URL)
	}
	if s.APIKey != "config-key" {
		t.Errorf("expected APIKey to remain config-key, got %q", s.APIKey)
	}
	if s.Verbose {
		t.Error("expected Verbose to remain false")
	}
}
