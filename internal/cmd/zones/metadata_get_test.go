// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v2/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/v2/internal/testutil"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestZonesMetadataGet(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
		GetMetadataFn: func(ctx context.Context, zoneID int, kind string) (*poweradmin.ZoneMetadata, *poweradmin.Response, error) {
			return &poweradmin.ZoneMetadata{Kind: kind, Values: []string{"192.0.2.10", "AUTO-NS"}}, nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewMetadataGetCmd(nil), []string{"--name", "example.com", "--kind", "ALLOW-AXFR-FROM"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "ALLOW-AXFR-FROM") {
		t.Errorf("expected output to contain ALLOW-AXFR-FROM, got:\n%s", out)
	}
	if !strings.Contains(out, "192.0.2.10") {
		t.Errorf("expected output to contain 192.0.2.10, got:\n%s", out)
	}
}

func TestZonesMetadataGetJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
		GetMetadataFn: func(ctx context.Context, zoneID int, kind string) (*poweradmin.ZoneMetadata, *poweradmin.Response, error) {
			return &poweradmin.ZoneMetadata{Kind: kind, Values: []string{"192.0.2.10"}}, nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewMetadataGetCmd(nil), []string{"--name", "example.com", "--kind", "ALLOW-AXFR-FROM", "--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"kind": "ALLOW-AXFR-FROM"`) {
		t.Errorf("expected JSON to contain kind ALLOW-AXFR-FROM, got:\n%s", out)
	}
}

func TestZonesMetadataGetMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, nil)
	err := fx.Run(zones.NewMetadataGetCmd(nil), []string{"--name", "example.com"})
	if err == nil {
		t.Fatal("expected error when --kind is missing")
	}
}

func TestZonesMetadataGetError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
		GetMetadataFn: func(ctx context.Context, zoneID int, kind string) (*poweradmin.ZoneMetadata, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	err := fx.Run(zones.NewMetadataGetCmd(nil), []string{"--name", "example.com", "--kind", "ALLOW-AXFR-FROM"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
