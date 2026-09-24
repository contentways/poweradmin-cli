// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestZonesDelete(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 42, Name: name, Type: "NATIVE"}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewDeleteCmd(nil), []string{"--name", "example.com", "--yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected output to contain example.com, got:\n%s", out)
	}
	if !strings.Contains(out, "42") {
		t.Errorf("expected output to contain id 42, got:\n%s", out)
	}
}

func TestZonesDeleteJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 42, Name: name, Type: "NATIVE"}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewDeleteCmd(nil), []string{"--name", "example.com", "--yes", "--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"name": "example.com"`) {
		t.Errorf("expected JSON to contain example.com, got:\n%s", out)
	}
	if !strings.Contains(out, `"id": 42`) {
		t.Errorf("expected JSON to contain id 42, got:\n%s", out)
	}
}

func TestZonesDeleteMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, nil)
	err := fx.Run(zones.NewDeleteCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestZonesDeleteError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	err := fx.Run(zones.NewDeleteCmd(nil), []string{"--name", "example.com", "--yes"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestZonesDeleteNilWithoutError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return nil, nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	err := fx.Run(zones.NewDeleteCmd(nil), []string{"--name", "ghost.com", "--yes"})
	if err == nil {
		t.Fatal("expected error for nil zone with nil error, got nil")
	}
}

// TestZonesDeleteInteractiveFullFlow drives the multi-select prompt for
// "zones delete --interactive" end to end: listing zones, toggling one,
// confirming the selection, then confirming deletion.
func TestZonesDeleteInteractiveFullFlow(t *testing.T) {
	testutil.WithAccessiblePrompts(t)

	var deletedID int
	mockZone := &testutil.MockZoneClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Zone, error) {
			return []*poweradmin.Zone{
				{ID: 1, Name: "example.com"},
				{ID: 2, Name: "other.com"},
			}, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			deletedID = id
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	// toggle zone 1 (example.com), 0 confirms the selection, then the
	// final "Proceed?" confirmation.
	testutil.WithDelayedStdin(t, "1\n", "0\n", "y\n")

	err := fx.Run(zones.NewDeleteCmd(nil), []string{"--interactive"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deletedID != 1 {
		t.Errorf("expected zone id 1 to be deleted, got %d", deletedID)
	}
	if !strings.Contains(fx.Stdout.String(), "deleted zone example.com") {
		t.Errorf("expected success output, got:\n%s", fx.Stdout.String())
	}
}

// TestZonesDeleteDryRun verifies that --dry-run prints what would be
// deleted without calling Delete or prompting for confirmation.
func TestZonesDeleteDryRun(t *testing.T) {
	deleteCalled := false
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 42, Name: name, Type: "NATIVE"}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			deleteCalled = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	// No stdin provided — dry-run must never prompt.

	err := fx.Run(zones.NewDeleteCmd(nil), []string{"--name", "example.com", "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleteCalled {
		t.Error("expected Delete to never be called with --dry-run")
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "Would delete zone example.com (id 42)") {
		t.Errorf("expected dry-run message, got:\n%s", out)
	}
}
