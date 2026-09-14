// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/internal/testutil"
)

func TestZonesList(t *testing.T) {
	// Arrange — set up a mock Zone client that returns two zones.
	mockZone := &testutil.MockZoneClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Zone, error) {
			return []*poweradmin.Zone{
				{ID: 1, Name: "example.com", Type: "NATIVE"},
				{ID: 2, Name: "example.org", Type: "NATIVE"},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	// Act — run the list command.
	err := fx.Run(zones.NewListCmd(nil), []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Assert — both zone names appear in the output.
	out := fx.Stdout.String()
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected output to contain example.com, got:\n%s", out)
	}
	if !strings.Contains(out, "example.org") {
		t.Errorf("expected output to contain example.org, got:\n%s", out)
	}
}

func TestZonesListJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Zone, error) {
			return []*poweradmin.Zone{
				{ID: 1, Name: "example.com", Type: "NATIVE"},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewListCmd(nil), []string{"--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"name": "example.com"`) {
		t.Errorf("expected JSON output to contain example.com, got:\n%s", out)
	}
}

func TestZonesListError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Zone, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	err := fx.Run(zones.NewListCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestZonesListFilterByType(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Zone, error) {
			return []*poweradmin.Zone{
				{ID: 1, Name: "example.com", Type: "NATIVE"},
				{ID: 2, Name: "example.org", Type: "MASTER"},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, nil, nil, nil, nil)

	err := fx.Run(zones.NewListCmd(nil), []string{"--type", "NATIVE"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected output to contain example.com, got:\n%s", out)
	}
	if strings.Contains(out, "example.org") {
		t.Errorf("expected output to NOT contain example.org, got:\n%s", out)
	}
}

func TestZonesListFilterByName(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Zone, error) {
			return []*poweradmin.Zone{
				{ID: 1, Name: "example.com", Type: "NATIVE"},
				{ID: 2, Name: "contentways.org", Type: "NATIVE"},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, nil, nil, nil, nil)

	err := fx.Run(zones.NewListCmd(nil), []string{"--name-filter", "contentways"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "contentways.org") {
		t.Errorf("expected output to contain contentways.org, got:\n%s", out)
	}
	if strings.Contains(out, "example.com") {
		t.Errorf("expected output to NOT contain example.com, got:\n%s", out)
	}
}
