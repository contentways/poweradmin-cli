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

func TestRecordsGet(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4", TTL: 3600},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, mockRecord, nil, nil, nil)

	err := fx.Run(records.NewGetCmd(nil), []string{
		"--zone-name", "example.com",
		"--id", "rec-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "www.example.com") {
		t.Errorf("expected output to contain www.example.com, got:\n%s", out)
	}
}

func TestRecordsGetJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4", TTL: 3600},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, mockRecord, nil, nil, nil)

	err := fx.Run(records.NewGetCmd(nil), []string{
		"--zone-name", "example.com",
		"--id", "rec-1",
		"-o", "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"name": "www.example.com"`) {
		t.Errorf("expected JSON to contain www.example.com, got:\n%s", out)
	}
}

func TestRecordsGetNotFound(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{}, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, mockRecord, nil, nil, nil)

	err := fx.Run(records.NewGetCmd(nil), []string{
		"--zone-name", "example.com",
		"--id", "nonexistent",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent record")
	}
}

func TestRecordsGetMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, &testutil.MockZoneClient{}, &testutil.MockRecordClient{}, nil, nil, nil)
	err := fx.Run(records.NewGetCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestRecordsGetError(t *testing.T) {
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

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, mockRecord, nil, nil, nil)

	err := fx.Run(records.NewGetCmd(nil), []string{
		"--zone-name", "example.com",
		"--id", "rec-1",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
