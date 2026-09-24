// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
	"github.com/spf13/cobra"
)

func TestDeleteSelectedRecordsAllSucceed(t *testing.T) {
	var deletedIDs []string
	mockRecord := &testutil.MockRecordClient{
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			deletedIDs = append(deletedIDs, recordID)
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, mockRecord)
	testutil.WithStdin(t, "y\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Record{
		{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4"},
		{ID: "rec-2", Name: "mail.example.com", Type: "A", Content: "1.2.3.5"},
	}

	err := deleteSelectedRecords(cmd, fx.State.MockClient, 42, selected, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deletedIDs) != 2 || deletedIDs[0] != "rec-1" || deletedIDs[1] != "rec-2" {
		t.Errorf("expected both records deleted in order, got %v", deletedIDs)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "deleted record www.example.com") || !strings.Contains(out, "deleted record mail.example.com") {
		t.Errorf("expected success lines for both records, got:\n%s", out)
	}
}

func TestDeleteSelectedRecordsPartialFailure(t *testing.T) {
	mockRecord := &testutil.MockRecordClient{
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			if recordID == "rec-2" {
				return nil, fmt.Errorf("record not found in this zone")
			}
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, mockRecord)
	testutil.WithStdin(t, "y\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Record{
		{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4"},
		{ID: "rec-2", Name: "mail.example.com", Type: "A", Content: "1.2.3.5"},
		{ID: "rec-3", Name: "ftp.example.com", Type: "A", Content: "1.2.3.6"},
	}

	err := deleteSelectedRecords(cmd, fx.State.MockClient, 42, selected, false)
	if err == nil {
		t.Fatal("expected error due to partial failure, got nil")
	}
	if !strings.Contains(err.Error(), "1 of 3") {
		t.Errorf("expected error to mention 1 of 3 failures, got: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "deleted record www.example.com") {
		t.Errorf("expected www.example.com to be deleted despite mail.example.com failing, got:\n%s", out)
	}
	if !strings.Contains(out, "deleted record ftp.example.com") {
		t.Errorf("expected ftp.example.com to be deleted despite mail.example.com failing, got:\n%s", out)
	}
	if !strings.Contains(out, "failed to delete record mail.example.com") {
		t.Errorf("expected failure line for mail.example.com, got:\n%s", out)
	}
}

func TestDeleteSelectedRecordsEmptySelectionNoOp(t *testing.T) {
	called := false
	mockRecord := &testutil.MockRecordClient{
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, mockRecord)

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	err := deleteSelectedRecords(cmd, fx.State.MockClient, 42, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected Delete to never be called for an empty selection")
	}
	if !strings.Contains(fx.Stdout.String(), "nothing to do") {
		t.Errorf("expected 'nothing to do' message, got:\n%s", fx.Stdout.String())
	}
}

func TestDeleteSelectedRecordsDeclineAbortsWithoutDeleting(t *testing.T) {
	called := false
	mockRecord := &testutil.MockRecordClient{
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, mockRecord)
	testutil.WithStdin(t, "n\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Record{
		{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4"},
	}

	err := deleteSelectedRecords(cmd, fx.State.MockClient, 42, selected, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected Delete to never be called after declining confirmation")
	}
}

// TestDeleteSelectedRecordsDryRunMakesNoAPICalls verifies that dry-run mode
// prints what would be deleted without calling Delete or prompting for
// confirmation at all.
func TestDeleteSelectedRecordsDryRunMakesNoAPICalls(t *testing.T) {
	called := false
	mockRecord := &testutil.MockRecordClient{
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, mockRecord)
	// No stdin provided — dry-run must never reach base.Confirm.

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Record{
		{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4"},
	}

	err := deleteSelectedRecords(cmd, fx.State.MockClient, 42, selected, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected Delete to never be called in dry-run mode")
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "Would delete 1 record(s)") {
		t.Errorf("expected dry-run summary, got:\n%s", out)
	}
	if !strings.Contains(out, "www.example.com") {
		t.Errorf("expected record name listed, got:\n%s", out)
	}
	if !strings.Contains(out, "No changes made") {
		t.Errorf("expected 'No changes made' confirmation, got:\n%s", out)
	}
}
