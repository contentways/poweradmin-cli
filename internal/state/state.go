// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package state holds the global runtime state passed to all CLI commands.
// It provides access to the Poweradmin API client and is threaded through
// the command tree via context.Context.
package state

import (
	"context"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
)

// State holds the global runtime configuration for the CLI.
// URL and APIKey are resolved from the config file, environment variables,
// or CLI flags — in increasing order of precedence.
type State struct {
	URL    string
	APIKey string
	// MockClient is used in tests to inject a mock client.
	// If nil, a real client is built from URL and APIKey.
	MockClient *poweradmin.Client
}

// ClientProvider is a function that returns a Poweradmin client.
// In production it builds a real client; in tests it returns a mock.
type ClientProvider func() (*poweradmin.Client, error)

// New creates a new State with the given Poweradmin instance URL and API key.
func New(url, apiKey string) *State {
	return &State{
		URL:    url,
		APIKey: apiKey,
	}
}

// Client builds and returns a new Poweradmin API client using the current
// URL and APIKey. A new client is created on each call — credentials may
// have been updated by a CLI flag after State was initialized.
func (s *State) Client() (*poweradmin.Client, error) {
	if s.MockClient != nil {
		return s.MockClient, nil
	}
	return poweradmin.NewClient(
		poweradmin.WithBaseURL(s.URL),
		poweradmin.WithAPIKey(s.APIKey),
	)
}

// WithContext returns a new context with the State embedded under stateKey.
// Call this in PersistentPreRunE to make the state available to all subcommands.
func (s *State) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, stateKey{}, s)
}

// stateKey is an unexported type used as the context key for State.
// Using a dedicated type prevents collisions with other context values.
type stateKey struct{}

// FromContext retrieves the State from the given context.
// Returns an empty State if no State was set — prevents nil pointer dereferences
// in commands when the context has not been initialized.
func FromContext(ctx context.Context) *State {
	s, ok := ctx.Value(stateKey{}).(*State)
	if !ok || s == nil {
		return &State{}
	}
	return s
}
