// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
	"github.com/spf13/cobra"
)

func TestDeleteSelectedUsersAllSucceed(t *testing.T) {
	var deletedIDs []int
	mockUser := &testutil.MockUserClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			deletedIDs = append(deletedIDs, id)
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)
	testutil.WithStdin(t, "y\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.User{
		{ID: 1, Username: "alice"},
		{ID: 2, Username: "bob"},
	}

	err := deleteSelectedUsers(cmd, fx.State.MockClient, selected, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deletedIDs) != 2 || deletedIDs[0] != 1 || deletedIDs[1] != 2 {
		t.Errorf("expected both users deleted in order, got %v", deletedIDs)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "deleted user alice") || !strings.Contains(out, "deleted user bob") {
		t.Errorf("expected success lines for both users, got:\n%s", out)
	}
}

func TestDeleteSelectedUsersPartialFailure(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			if id == 2 {
				return nil, fmt.Errorf("permission denied")
			}
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)
	testutil.WithStdin(t, "y\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.User{
		{ID: 1, Username: "alice"},
		{ID: 2, Username: "bob"},
		{ID: 3, Username: "carol"},
	}

	err := deleteSelectedUsers(cmd, fx.State.MockClient, selected, false)
	if err == nil {
		t.Fatal("expected error due to partial failure, got nil")
	}
	if !strings.Contains(err.Error(), "1 of 3") {
		t.Errorf("expected error to mention 1 of 3 failures, got: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "deleted user alice") {
		t.Errorf("expected alice to be deleted despite bob failing, got:\n%s", out)
	}
	if !strings.Contains(out, "deleted user carol") {
		t.Errorf("expected carol to be deleted despite bob failing, got:\n%s", out)
	}
	if !strings.Contains(out, "failed to delete user bob") {
		t.Errorf("expected failure line for bob, got:\n%s", out)
	}
}

func TestDeleteSelectedUsersEmptySelectionNoOp(t *testing.T) {
	called := false
	mockUser := &testutil.MockUserClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	err := deleteSelectedUsers(cmd, fx.State.MockClient, nil, false)
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

func TestDeleteSelectedUsersDeclineAbortsWithoutDeleting(t *testing.T) {
	called := false
	mockUser := &testutil.MockUserClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)
	testutil.WithStdin(t, "n\n")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.User{{ID: 1, Username: "alice"}}

	err := deleteSelectedUsers(cmd, fx.State.MockClient, selected, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected Delete to never be called after declining confirmation")
	}
}

// TestDeleteSelectedUsersDryRunMakesNoAPICalls verifies that dry-run mode
// prints what would be deleted without calling Delete or prompting for
// confirmation at all.
func TestDeleteSelectedUsersDryRunMakesNoAPICalls(t *testing.T) {
	called := false
	mockUser := &testutil.MockUserClient{
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)
	// No stdin provided — dry-run must never reach base.Confirm.

	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(fx.Stdout)

	selected := []*poweradmin.User{
		{ID: 1, Username: "alice"},
		{ID: 2, Username: "bob"},
	}

	err := deleteSelectedUsers(cmd, fx.State.MockClient, selected, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected Delete to never be called in dry-run mode")
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "Would delete 2 user(s)") {
		t.Errorf("expected dry-run summary, got:\n%s", out)
	}
	if !strings.Contains(out, "alice") || !strings.Contains(out, "bob") {
		t.Errorf("expected both usernames listed, got:\n%s", out)
	}
	if !strings.Contains(out, "No changes made") {
		t.Errorf("expected 'No changes made' confirmation, got:\n%s", out)
	}
}
