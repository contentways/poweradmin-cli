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

func TestRecordsDelete(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewDeleteCmd(nil), []string{
		"--zone-name", "example.com",
		"--id", "rec-42",
		"--yes",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "rec-42") {
		t.Errorf("expected output to contain rec-42, got:\n%s", out)
	}
}

func TestRecordsDeleteJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewDeleteCmd(nil), []string{
		"--zone-name", "example.com",
		"--id", "rec-42",
		"--output", "json",
		"--yes",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"id": "rec-42"`) {
		t.Errorf("expected JSON to contain rec-42, got:\n%s", out)
	}
}

func TestRecordsDeleteMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, &testutil.MockRecordClient{})
	err := fx.Run(records.NewDeleteCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestRecordsDeleteError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)
	err := fx.Run(records.NewDeleteCmd(nil), []string{
		"--zone-name", "example.com",
		"--id", "rec-42",
		"--yes",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRecordsDeleteByZoneID(t *testing.T) {
	mockRecord := &testutil.MockRecordClient{
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			if zoneID != 7 {
				t.Errorf("expected zoneID 7, got %d", zoneID)
			}
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, mockRecord)

	err := fx.Run(records.NewDeleteCmd(nil), []string{
		"--zone-id", "7",
		"--id", "rec-42",
		"--yes",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecordsDeleteQuiet(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewDeleteCmd(nil), []string{
		"--zone-name", "example.com",
		"--id", "rec-42",
		"--yes",
		"--quiet",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fx.Stdout.String() != "" {
		t.Errorf("expected no output in quiet mode, got:\n%s", fx.Stdout.String())
	}
}

// TestRecordsDeleteDryRun verifies that --dry-run prints what would be
// deleted without calling Delete or prompting for confirmation.
func TestRecordsDeleteDryRun(t *testing.T) {
	deleteCalled := false
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			deleteCalled = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)
	// No stdin provided — dry-run must never prompt.

	err := fx.Run(records.NewDeleteCmd(nil), []string{
		"--zone-name", "example.com",
		"--id", "rec-42",
		"--dry-run",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleteCalled {
		t.Error("expected Delete to never be called with --dry-run")
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "Would delete record (id rec-42)") {
		t.Errorf("expected dry-run message, got:\n%s", out)
	}
}
