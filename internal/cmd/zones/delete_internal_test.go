// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
	"github.com/spf13/cobra"
)

func TestDeleteSelectedZonesAllSucceed(t *testing.T) {
	var deletedIDs []int
	mockZone := &testutil.MockZoneClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			deletedIDs = append(deletedIDs, id)
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	testutil.WithStdin(t, "y\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Zone{
		{ID: 1, Name: "a.com"},
		{ID: 2, Name: "b.com"},
	}

	err := deleteSelectedZones(cmd, fx.State.MockClient, selected, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deletedIDs) != 2 || deletedIDs[0] != 1 || deletedIDs[1] != 2 {
		t.Errorf("expected both zones deleted in order, got %v", deletedIDs)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "deleted zone a.com") || !strings.Contains(out, "deleted zone b.com") {
		t.Errorf("expected success lines for both zones, got:\n%s", out)
	}
}

func TestDeleteSelectedZonesPartialFailure(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			if id == 2 {
				return nil, fmt.Errorf("permission denied")
			}
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	testutil.WithStdin(t, "y\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Zone{
		{ID: 1, Name: "a.com"},
		{ID: 2, Name: "b.com"},
		{ID: 3, Name: "c.com"},
	}

	err := deleteSelectedZones(cmd, fx.State.MockClient, selected, false)
	if err == nil {
		t.Fatal("expected error due to partial failure, got nil")
	}
	if !strings.Contains(err.Error(), "1 of 3") {
		t.Errorf("expected error to mention 1 of 3 failures, got: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "deleted zone a.com") {
		t.Errorf("expected a.com to be deleted despite b.com failing, got:\n%s", out)
	}
	if !strings.Contains(out, "deleted zone c.com") {
		t.Errorf("expected c.com to be deleted despite b.com failing, got:\n%s", out)
	}
	if !strings.Contains(out, "failed to delete zone b.com") {
		t.Errorf("expected failure line for b.com, got:\n%s", out)
	}
}

func TestDeleteSelectedZonesEmptySelectionNoOp(t *testing.T) {
	called := false
	mockZone := &testutil.MockZoneClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	err := deleteSelectedZones(cmd, fx.State.MockClient, nil, false)
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

func TestDeleteSelectedZonesDeclineAbortsWithoutDeleting(t *testing.T) {
	called := false
	mockZone := &testutil.MockZoneClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	testutil.WithStdin(t, "n\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Zone{{ID: 1, Name: "a.com"}}

	err := deleteSelectedZones(cmd, fx.State.MockClient, selected, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected Delete to never be called after declining confirmation")
	}
}

// TestDeleteSelectedZonesDryRunMakesNoAPICalls verifies that dry-run mode
// prints what would be deleted without calling Delete or prompting for
// confirmation at all.
func TestDeleteSelectedZonesDryRunMakesNoAPICalls(t *testing.T) {
	called := false
	mockZone := &testutil.MockZoneClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	// No stdin provided — dry-run must never reach base.Confirm, or this
	// test would hang waiting for input.

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Zone{
		{ID: 1, Name: "a.com"},
		{ID: 2, Name: "b.com"},
	}

	err := deleteSelectedZones(cmd, fx.State.MockClient, selected, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected Delete to never be called in dry-run mode")
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "Would delete 2 zone(s)") {
		t.Errorf("expected dry-run summary, got:\n%s", out)
	}
	if !strings.Contains(out, "a.com") || !strings.Contains(out, "b.com") {
		t.Errorf("expected both zone names listed, got:\n%s", out)
	}
	if !strings.Contains(out, "No changes made") {
		t.Errorf("expected 'No changes made' confirmation, got:\n%s", out)
	}
}

// TestDeleteSelectedZonesDryRunEmptySelectionNoOp verifies dry-run with an
// empty selection behaves the same as the non-dry-run empty case.
func TestDeleteSelectedZonesDryRunEmptySelectionNoOp(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	err := deleteSelectedZones(cmd, fx.State.MockClient, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "nothing to do") {
		t.Errorf("expected 'nothing to do' message, got:\n%s", fx.Stdout.String())
	}
}
