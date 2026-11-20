// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones_test

import (
	"context"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func metadataMock() *testutil.MockZoneClient {
	return &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
		ListMetadataFn: func(ctx context.Context, zoneID int) ([]*poweradmin.ZoneMetadata, *poweradmin.Response, error) {
			return []*poweradmin.ZoneMetadata{{Kind: "ALLOW-AXFR-FROM", Values: []string{"192.0.2.10"}}}, nil, nil
		},
		GetMetadataFn: func(ctx context.Context, zoneID int, kind string) (*poweradmin.ZoneMetadata, *poweradmin.Response, error) {
			return &poweradmin.ZoneMetadata{Kind: kind, Values: []string{"192.0.2.10"}}, nil, nil
		},
	}
}

func TestZonesMetadataGroupListsWithoutSubcommand(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, metadataMock(), nil)

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"metadata", "--name", "example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "ALLOW-AXFR-FROM") {
		t.Errorf("expected metadata listing, got:\n%s", fx.Stdout.String())
	}
}

func TestZonesMetadataListSubcommand(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, metadataMock(), nil)

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"metadata", "list", "--name", "example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "192.0.2.10") {
		t.Errorf("expected metadata listing, got:\n%s", fx.Stdout.String())
	}
}

func TestZonesMetadataGetNested(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, metadataMock(), nil)

	err := fx.Run(zones.NewZonesCommand(nil), []string{"metadata", "get", "--name", "example.com", "--kind", "ALLOW-AXFR-FROM"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "Kind: ALLOW-AXFR-FROM") {
		t.Errorf("expected metadata value, got:\n%s", out)
	}
	if strings.Contains(out, "deprecated") {
		t.Errorf("nested command must not print a deprecation note, got:\n%s", out)
	}
}

func TestZonesMetadataGetDeprecatedAlias(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, metadataMock(), nil)

	err := fx.Run(zones.NewZonesCommand(nil), []string{"metadata-get", "--name", "example.com", "--kind", "ALLOW-AXFR-FROM"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "Kind: ALLOW-AXFR-FROM") {
		t.Errorf("alias must still work, got:\n%s", out)
	}
	if !strings.Contains(out+fx.Stderr.String(), `use "zones metadata get" instead`) {
		t.Errorf("expected deprecation note, got stdout:\n%s\nstderr:\n%s", out, fx.Stderr.String())
	}
}

func TestZonesHelpHidesDeprecatedAliases(t *testing.T) {
	fx := testutil.NewFixture(t)

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"--help"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if strings.Contains(out, "metadata-get") {
		t.Errorf("help must not list deprecated aliases, got:\n%s", out)
	}
	if !strings.Contains(out, "metadata") {
		t.Errorf("help must list the metadata group, got:\n%s", out)
	}
}
