// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
	"github.com/spf13/cobra"
)

func TestDeleteSelectedGroupsAllSucceed(t *testing.T) {
	var deletedIDs []int
	mockGroup := &testutil.MockGroupClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			deletedIDs = append(deletedIDs, id)
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	testutil.WithStdin(t, "y\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Group{
		{ID: 1, Name: "Editors"},
		{ID: 2, Name: "Viewers"},
	}

	err := deleteSelectedGroups(cmd, fx.State.MockClient, selected, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deletedIDs) != 2 || deletedIDs[0] != 1 || deletedIDs[1] != 2 {
		t.Errorf("expected both groups deleted in order, got %v", deletedIDs)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "deleted group Editors") || !strings.Contains(out, "deleted group Viewers") {
		t.Errorf("expected success lines for both groups, got:\n%s", out)
	}
}

func TestDeleteSelectedGroupsPartialFailure(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			if id == 2 {
				return nil, fmt.Errorf("permission denied")
			}
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	testutil.WithStdin(t, "y\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Group{
		{ID: 1, Name: "Editors"},
		{ID: 2, Name: "Viewers"},
		{ID: 3, Name: "Admins"},
	}

	err := deleteSelectedGroups(cmd, fx.State.MockClient, selected, false)
	if err == nil {
		t.Fatal("expected error due to partial failure, got nil")
	}
	if !strings.Contains(err.Error(), "1 of 3") {
		t.Errorf("expected error to mention 1 of 3 failures, got: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "deleted group Editors") {
		t.Errorf("expected Editors to be deleted despite Viewers failing, got:\n%s", out)
	}
	if !strings.Contains(out, "deleted group Admins") {
		t.Errorf("expected Admins to be deleted despite Viewers failing, got:\n%s", out)
	}
	if !strings.Contains(out, "failed to delete group Viewers") {
		t.Errorf("expected failure line for Viewers, got:\n%s", out)
	}
}

func TestDeleteSelectedGroupsEmptySelectionNoOp(t *testing.T) {
	called := false
	mockGroup := &testutil.MockGroupClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	err := deleteSelectedGroups(cmd, fx.State.MockClient, nil, false)
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

func TestDeleteSelectedGroupsDeclineAbortsWithoutDeleting(t *testing.T) {
	called := false
	mockGroup := &testutil.MockGroupClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	testutil.WithStdin(t, "n\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Group{{ID: 1, Name: "Editors"}}

	err := deleteSelectedGroups(cmd, fx.State.MockClient, selected, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected Delete to never be called after declining confirmation")
	}
}

// TestDeleteSelectedGroupsDryRunMakesNoAPICalls verifies that dry-run mode
// prints what would be deleted without calling Delete or prompting for
// confirmation at all.
func TestDeleteSelectedGroupsDryRunMakesNoAPICalls(t *testing.T) {
	called := false
	mockGroup := &testutil.MockGroupClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	// No stdin provided — dry-run must never reach base.Confirm.

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.Group{
		{ID: 1, Name: "Editors"},
		{ID: 2, Name: "Viewers"},
	}

	err := deleteSelectedGroups(cmd, fx.State.MockClient, selected, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected Delete to never be called in dry-run mode")
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "Would delete 2 group(s)") {
		t.Errorf("expected dry-run summary, got:\n%s", out)
	}
	if !strings.Contains(out, "Editors") || !strings.Contains(out, "Viewers") {
		t.Errorf("expected both group names listed, got:\n%s", out)
	}
	if !strings.Contains(out, "No changes made") {
		t.Errorf("expected 'No changes made' confirmation, got:\n%s", out)
	}
}
