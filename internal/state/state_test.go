// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package state_test

import (
	"context"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestNew(t *testing.T) {
	s := state.New("https://dns.example.com", "pwa_test")
	if s.URL != "https://dns.example.com" {
		t.Errorf("URL = %q, want https://dns.example.com", s.URL)
	}
	if s.APIKey != "pwa_test" {
		t.Errorf("APIKey = %q, want pwa_test", s.APIKey)
	}
}

func TestWithContextAndFromContext(t *testing.T) {
	s := state.New("https://dns.example.com", "pwa_test")
	ctx := s.WithContext(context.Background())

	got := state.FromContext(ctx)
	if got == nil {
		t.Fatal("expected State from context, got nil")
	}
	if got.URL != s.URL {
		t.Errorf("URL = %q, want %q", got.URL, s.URL)
	}
}

func TestFromContextEmpty(t *testing.T) {
	got := state.FromContext(context.Background())
	// FromContext returns an empty State, never nil.
	if got.URL != "" || got.APIKey != "" {
		t.Errorf("expected empty State, got: %+v", got)
	}
}

func TestClientWithMockClient(t *testing.T) {
	s := &state.State{
		MockClient: &poweradmin.Client{},
	}

	client, err := s.Client()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client != s.MockClient {
		t.Fatal("expected mock client")
	}
}

func TestClientBuildsRealClientWithoutMock(t *testing.T) {
	s := state.New("https://dns.example.com", "pwa_test")

	client, err := s.Client()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected a non-nil client")
	}
	if client == s.MockClient {
		t.Error("expected a freshly built client, not the (nil) mock")
	}
}

func TestClientVerboseEnablesDebugWriter(t *testing.T) {
	s := state.New("https://dns.example.com", "pwa_test")
	s.Verbose = true

	client, err := s.Client()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected a non-nil client")
	}
}

func TestClientMissingURLReturnsError(t *testing.T) {
	s := state.New("", "pwa_test")

	_, err := s.Client()
	if err == nil {
		t.Fatal("expected error for missing URL, got nil")
	}
}
