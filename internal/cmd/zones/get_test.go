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

func TestZonesGetByName(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{Type: "NS", Content: "ns1.example.com"},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(zones.NewGetCmd(nil), []string{"--name", "example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected output to contain example.com, got:\n%s", out)
	}
	if !strings.Contains(out, "ns1.example.com") {
		t.Errorf("expected output to contain ns1.example.com, got:\n%s", out)
	}
}

func TestZonesGetByNameJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{Type: "NS", Content: "ns1.example.com"},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(zones.NewGetCmd(nil), []string{"--name", "example.com", "--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"name": "example.com"`) {
		t.Errorf("expected JSON to contain example.com, got:\n%s", out)
	}
	if !strings.Contains(out, "ns1.example.com") {
		t.Errorf("expected JSON to contain ns1.example.com, got:\n%s", out)
	}
}

func TestZonesGetMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, nil)
	err := fx.Run(zones.NewGetCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestZonesGetByNameError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	err := fx.Run(zones.NewGetCmd(nil), []string{"--name", "example.com"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
