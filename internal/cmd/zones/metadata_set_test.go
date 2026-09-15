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

func TestZonesMetadataSet(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
		SetMetadataFn: func(ctx context.Context, zoneID int, kind string, values []string) (*poweradmin.Response, error) {
			if kind != "ALLOW-AXFR-FROM" {
				t.Errorf("kind = %s, want ALLOW-AXFR-FROM", kind)
			}
			if len(values) != 2 {
				t.Errorf("values = %v, want 2 entries", values)
			}
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewMetadataSetCmd(nil), []string{
		"--name", "example.com",
		"--kind", "ALLOW-AXFR-FROM",
		"--values", "192.0.2.10,AUTO-NS",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "ALLOW-AXFR-FROM") {
		t.Errorf("expected output to contain ALLOW-AXFR-FROM, got:\n%s", out)
	}
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected output to contain example.com, got:\n%s", out)
	}
}

func TestZonesMetadataSetJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
		SetMetadataFn: func(ctx context.Context, zoneID int, kind string, values []string) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewMetadataSetCmd(nil), []string{
		"--name", "example.com",
		"--kind", "ALLOW-AXFR-FROM",
		"--values", "192.0.2.10",
		"--output", "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"kind": "ALLOW-AXFR-FROM"`) {
		t.Errorf("expected JSON to contain kind ALLOW-AXFR-FROM, got:\n%s", out)
	}
}

func TestZonesMetadataSetQuiet(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
		SetMetadataFn: func(ctx context.Context, zoneID int, kind string, values []string) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewMetadataSetCmd(nil), []string{
		"--name", "example.com",
		"--kind", "ALLOW-AXFR-FROM",
		"--values", "192.0.2.10",
		"--quiet",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out := fx.Stdout.String(); out != "" {
		t.Errorf("expected no output in quiet mode, got:\n%s", out)
	}
}

func TestZonesMetadataSetMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, nil)
	err := fx.Run(zones.NewMetadataSetCmd(nil), []string{"--name", "example.com", "--kind", "ALLOW-AXFR-FROM"})
	if err == nil {
		t.Fatal("expected error when --values is missing")
	}
}

func TestZonesMetadataSetError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
		SetMetadataFn: func(ctx context.Context, zoneID int, kind string, values []string) (*poweradmin.Response, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	err := fx.Run(zones.NewMetadataSetCmd(nil), []string{
		"--name", "example.com",
		"--kind", "ALLOW-AXFR-FROM",
		"--values", "192.0.2.10",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
