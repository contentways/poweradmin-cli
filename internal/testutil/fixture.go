// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package testutil provides shared test helpers for the poweradmin-cli commands.
// It follows the fixture pattern used by hetznercloud/cli — each test creates
// a Fixture that wires up a mock State and exposes a Run method to execute
// commands and assert their output.
package testutil

import (
	"bytes"
	"context"
	"testing"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// Fixture holds the test state and output buffers for a single command test.
type Fixture struct {
	State  *state.State
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

// NewFixture creates a new Fixture with real (but unused) credentials.
// Use NewFixtureWithMocks to inject mock sub-clients.
func NewFixture(t *testing.T) *Fixture {
	t.Helper()
	return &Fixture{
		State:  state.New("https://test.example.com", "test-key"),
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
	}
}

// NewFixtureWithMocks creates a Fixture with mock Zone and Record clients
// injected into a real *poweradmin.Client shell. This allows commands to
// call client.Zone.All(...) etc. without making real HTTP requests.
func NewFixtureWithMocks(t *testing.T, zone poweradmin.IZoneClient, record poweradmin.IRecordClient) *Fixture {
	t.Helper()

	// Build a real client shell — credentials are unused since all sub-clients
	// are replaced with mocks immediately after construction.
	client, _ := poweradmin.NewClient(
		poweradmin.WithBaseURL("https://test.example.com"),
		poweradmin.WithAPIKey("test-key"),
	)
	client.Zone = zone
	client.Record = record

	s := state.New("https://test.example.com", "test-key")
	s.MockClient = client

	return &Fixture{
		State:  s,
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
	}
}

// Run executes the given Cobra command with the provided arguments.
// The State is injected into the context so all subcommands can access it
// via state.FromContext(cmd.Context()).
func (f *Fixture) Run(cmd *cobra.Command, args []string) error {
	ctx := f.State.WithContext(context.Background())
	cmd.SetContext(ctx)
	cmd.SetOut(f.Stdout)
	cmd.SetErr(f.Stderr)
	cmd.SetArgs(args)

	// Disable usage printing on error to keep test output clean.
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd.Execute()
}

// NewFixtureWithAllMocks creates a Fixture with mock Zone, Record, User, Group
// and PermissionTemplate clients.
func NewFixtureWithAllMocks(t *testing.T, zone poweradmin.IZoneClient, record poweradmin.IRecordClient, user poweradmin.IUserClient, group poweradmin.IGroupClient, permTemplate poweradmin.IPermissionTemplateClient) *Fixture {
	t.Helper()
	client, _ := poweradmin.NewClient(
		poweradmin.WithBaseURL("https://test.example.com"),
		poweradmin.WithAPIKey("test-key"),
	)
	client.Zone = zone
	client.Record = record
	client.User = user
	client.Group = group
	client.PermissionTemplate = permTemplate
	s := state.New("https://test.example.com", "test-key")
	s.MockClient = client
	return &Fixture{
		State:  s,
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
	}
}
