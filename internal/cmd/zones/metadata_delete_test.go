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

func TestZonesMetadataDelete(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 42, Name: name, Type: "NATIVE"}, nil, nil
		},
		DeleteMetadataFn: func(ctx context.Context, zoneID int, kind string) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewMetadataDeleteCmd(nil), []string{
		"--name", "example.com",
		"--kind", "ALLOW-AXFR-FROM",
		"--yes",
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

func TestZonesMetadataDeleteQuiet(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 42, Name: name, Type: "NATIVE"}, nil, nil
		},
		DeleteMetadataFn: func(ctx context.Context, zoneID int, kind string) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewMetadataDeleteCmd(nil), []string{
		"--name", "example.com",
		"--kind", "ALLOW-AXFR-FROM",
		"--yes",
		"--quiet",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out := fx.Stdout.String(); out != "" {
		t.Errorf("expected no output in quiet mode, got:\n%s", out)
	}
}

func TestZonesMetadataDeleteMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, nil)
	err := fx.Run(zones.NewMetadataDeleteCmd(nil), []string{"--name", "example.com", "--yes"})
	if err == nil {
		t.Fatal("expected error when --kind is missing")
	}
}

func TestZonesMetadataDeleteError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
		DeleteMetadataFn: func(ctx context.Context, zoneID int, kind string) (*poweradmin.Response, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	err := fx.Run(zones.NewMetadataDeleteCmd(nil), []string{
		"--name", "example.com",
		"--kind", "ALLOW-AXFR-FROM",
		"--yes",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
