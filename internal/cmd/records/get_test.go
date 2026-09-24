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
	"github.com/contentways/poweradmin-go/v4/poweradmin"
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

func TestRecordsGetByZoneID(t *testing.T) {
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
	fx := testutil.NewFixtureWithAllMocks(t, &testutil.MockZoneClient{}, mockRecord, nil, nil, nil)

	err := fx.Run(records.NewGetCmd(nil), []string{
		"--zone-id", "9",
		"--id", "rec-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// A priority of 0 is meaningful for MX/SRV and must be shown; types without
// a priority must not show the line at all.
func TestRecordsGetPriority(t *testing.T) {
	for _, tc := range []struct {
		rec  poweradmin.Record
		want bool
	}{
		{poweradmin.Record{ID: "r1", Name: "_imaps._tcp.example.com", Type: "SRV", Content: "1 993 mail.example.com", TTL: 60}, true},
		{poweradmin.Record{ID: "r1", Name: "www.example.com", Type: "A", Content: "192.0.2.1", TTL: 60}, false},
	} {
		t.Run(tc.rec.Type, func(t *testing.T) {
			rec := tc.rec
			mockZone := &testutil.MockZoneClient{
				GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
					return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
				},
			}
			mockRecord := &testutil.MockRecordClient{
				AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
					return []*poweradmin.Record{&rec}, nil
				},
			}
			fx := testutil.NewFixtureWithAllMocks(t, mockZone, mockRecord, nil, nil, nil)
			if err := fx.Run(records.NewGetCmd(nil), []string{"--zone-name", "example.com", "--id", "r1"}); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := strings.Contains(fx.Stdout.String(), "Priority: 0"); got != tc.want {
				t.Errorf("Priority line shown = %v, want %v\n%s", got, tc.want, fx.Stdout.String())
			}
		})
	}
}
