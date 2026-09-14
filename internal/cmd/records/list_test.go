// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/records"
	"github.com/contentways/poweradmin-cli/internal/testutil"
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
