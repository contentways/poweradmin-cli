// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/records"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestRecordsList(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4", TTL: 3600},
				{ID: "rec-2", Name: "mail.example.com", Type: "MX", Content: "mail.example.com", TTL: 3600},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewListCmd(nil), []string{"--zone-name", "example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "www.example.com") {
		t.Errorf("expected output to contain www.example.com, got:\n%s", out)
	}
	if !strings.Contains(out, "1.2.3.4") {
		t.Errorf("expected output to contain 1.2.3.4, got:\n%s", out)
	}
}

func TestRecordsListJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4", TTL: 3600},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewListCmd(nil), []string{"--zone-name", "example.com", "--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"name": "www.example.com"`) {
		t.Errorf("expected JSON to contain www.example.com, got:\n%s", out)
	}
}

func TestRecordsListMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, &testutil.MockRecordClient{})
	err := fx.Run(records.NewListCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestRecordsListError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)
	err := fx.Run(records.NewListCmd(nil), []string{"--zone-name", "example.com"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRecordsListByZoneID(t *testing.T) {
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			if zoneID != 9 {
				t.Errorf("expected zoneID 9, got %d", zoneID)
			}
			return []*poweradmin.Record{
				{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4"},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, mockRecord)

	err := fx.Run(records.NewListCmd(nil), []string{"--zone-id", "9"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestRecordsListFilterByType assumes "records list" filters client-side by
// --type, mirroring "zones list --type". If --type isn't implemented this
// way, this test will fail at runtime (unknown flag) rather than at compile
// time — check the failure message if so.
func TestRecordsListFilterByType(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4"},
				{ID: "rec-2", Name: "example.com", Type: "MX", Content: "mail.example.com"},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewListCmd(nil), []string{"--zone-name", "example.com", "--type", "MX"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if strings.Contains(out, "www.example.com") {
		t.Errorf("expected A record to be filtered out, got:\n%s", out)
	}
	if !strings.Contains(out, "mail.example.com") {
		t.Errorf("expected MX record content present, got:\n%s", out)
	}
}
