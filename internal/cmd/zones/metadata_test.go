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
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestZonesMetadataList(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
		ListMetadataFn: func(ctx context.Context, zoneID int) ([]*poweradmin.ZoneMetadata, *poweradmin.Response, error) {
			return []*poweradmin.ZoneMetadata{
				{Kind: "ALLOW-AXFR-FROM", Values: []string{"192.0.2.10", "AUTO-NS"}},
			}, nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewMetadataCmd(nil), []string{"--name", "example.com"})
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

func TestZonesMetadataListJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
		ListMetadataFn: func(ctx context.Context, zoneID int) ([]*poweradmin.ZoneMetadata, *poweradmin.Response, error) {
			return []*poweradmin.ZoneMetadata{
				{Kind: "ALLOW-AXFR-FROM", Values: []string{"192.0.2.10"}},
			}, nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewMetadataCmd(nil), []string{"--name", "example.com", "--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"kind": "ALLOW-AXFR-FROM"`) {
		t.Errorf("expected JSON to contain kind ALLOW-AXFR-FROM, got:\n%s", out)
	}
	if !strings.Contains(out, `"count": 1`) {
		t.Errorf("expected JSON to contain count 1, got:\n%s", out)
	}
}

func TestZonesMetadataListMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, nil)
	err := fx.Run(zones.NewMetadataCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestZonesMetadataListError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
		ListMetadataFn: func(ctx context.Context, zoneID int) ([]*poweradmin.ZoneMetadata, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	err := fx.Run(zones.NewMetadataCmd(nil), []string{"--name", "example.com"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
