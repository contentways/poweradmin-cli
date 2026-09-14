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

func TestZonesExport(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{
					ID:      "rec-soa",
					Name:    "example.com",
					Type:    "SOA",
					Content: "ns1.example.com. admin.example.com. 2026060601 28800 7200 1209600 86400",
					TTL:     86400,
				},
				{
					ID:      "rec-ns1",
					Name:    "example.com",
					Type:    "NS",
					Content: "ns1.example.com",
					TTL:     3600,
				},
				{
					ID:      "rec-a",
					Name:    "www.example.com",
					Type:    "A",
					Content: "1.2.3.4",
					TTL:     3600,
				},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, mockRecord, nil, nil, nil)

	err := fx.Run(zones.NewExportCmd(nil), []string{"--name", "example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "$ORIGIN example.com.") {
		t.Errorf("expected $ORIGIN, got:\n%s", out)
	}
	if !strings.Contains(out, "SOA") {
		t.Errorf("expected SOA record, got:\n%s", out)
	}
	if !strings.Contains(out, "ns1.example.com.") {
		t.Errorf("expected nameserver, got:\n%s", out)
	}
	if !strings.Contains(out, "www") {
		t.Errorf("expected www A record, got:\n%s", out)
	}
}

func TestZonesExportMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, &testutil.MockZoneClient{}, &testutil.MockRecordClient{}, nil, nil, nil)
	err := fx.Run(zones.NewExportCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestZonesExportError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, mockZone, &testutil.MockRecordClient{}, nil, nil, nil)
	err := fx.Run(zones.NewExportCmd(nil), []string{"--name", "example.com"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
