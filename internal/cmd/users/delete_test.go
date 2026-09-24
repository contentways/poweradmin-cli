// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/users"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestUsersDelete(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username, Email: "max@example.com"}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int, opts poweradmin.UserDeleteOpts) (int, *poweradmin.Response, error) {
			return 0, nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewDeleteCmd(nil), []string{"--name", "max", "--yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "max") {
		t.Errorf("expected output to contain max, got:\n%s", out)
	}
	if !strings.Contains(out, "42") {
		t.Errorf("expected output to contain id 42, got:\n%s", out)
	}
}

func TestUsersDeleteJSON(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int, opts poweradmin.UserDeleteOpts) (int, *poweradmin.Response, error) {
			return 0, nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewDeleteCmd(nil), []string{"--name", "max", "--yes", "--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"id": 42`) {
		t.Errorf("expected JSON to contain id 42, got:\n%s", out)
	}
}

func TestUsersDeleteMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, &testutil.MockUserClient{}, nil, nil)

	err := fx.Run(users.NewDeleteCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestUsersDeleteError(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int, opts poweradmin.UserDeleteOpts) (int, *poweradmin.Response, error) {
			return 0, nil, fmt.Errorf("api error")
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewDeleteCmd(nil), []string{"--name", "max", "--yes"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUsersDeleteByID(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: id, Username: "max"}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int, opts poweradmin.UserDeleteOpts) (int, *poweradmin.Response, error) {
			return 0, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewDeleteCmd(nil), []string{"--id", "42", "--yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUsersDeleteQuiet(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int, opts poweradmin.UserDeleteOpts) (int, *poweradmin.Response, error) {
			return 0, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewDeleteCmd(nil), []string{"--name", "max", "--yes", "--quiet"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fx.Stdout.String() != "" {
		t.Errorf("expected no output in quiet mode, got:\n%s", fx.Stdout.String())
	}
}

// TestUsersDeleteDryRun verifies that --dry-run prints what would be
// deleted without calling Delete or prompting for confirmation.
func TestUsersDeleteDryRun(t *testing.T) {
	deleteCalled := false
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int, opts poweradmin.UserDeleteOpts) (int, *poweradmin.Response, error) {
			deleteCalled = true
			return 0, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)
	// No stdin provided — dry-run must never prompt.

	err := fx.Run(users.NewDeleteCmd(nil), []string{"--name", "max", "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleteCalled {
		t.Error("expected Delete to never be called with --dry-run")
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "Would delete user max (id 42)") {
		t.Errorf("expected dry-run message, got:\n%s", out)
	}
}

func TestUsersDeleteTransferToByID(t *testing.T) {
	var gotOpts poweradmin.UserDeleteOpts
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int, opts poweradmin.UserDeleteOpts) (int, *poweradmin.Response, error) {
			gotOpts = opts
			return 3, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	if err := fx.Run(users.NewDeleteCmd(nil), []string{"--name", "max", "--transfer-to", "7", "--yes"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotOpts.TransferToUserID == nil || *gotOpts.TransferToUserID != 7 {
		t.Errorf("TransferToUserID = %v, want 7", gotOpts.TransferToUserID)
	}
	if out := fx.Stdout.String(); !strings.Contains(out, "3 zone(s) transferred") {
		t.Errorf("expected transfer count in output, got:\n%s", out)
	}
}

func TestUsersDeleteTransferToByName(t *testing.T) {
	var gotOpts poweradmin.UserDeleteOpts
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			ids := map[string]int{"max": 42, "anna": 9}
			return &poweradmin.User{ID: ids[username], Username: username}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int, opts poweradmin.UserDeleteOpts) (int, *poweradmin.Response, error) {
			if id != 42 {
				t.Errorf("deleted id = %d, want 42", id)
			}
			gotOpts = opts
			return 1, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	if err := fx.Run(users.NewDeleteCmd(nil), []string{"--name", "max", "--transfer-to", "anna", "--yes", "-o", "json"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotOpts.TransferToUserID == nil || *gotOpts.TransferToUserID != 9 {
		t.Errorf("TransferToUserID = %v, want 9", gotOpts.TransferToUserID)
	}
	if out := fx.Stdout.String(); !strings.Contains(out, `"zones_transferred": 1`) {
		t.Errorf("expected zones_transferred in JSON, got:\n%s", out)
	}
}

func TestUsersDeleteTransferToUnknownUser(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			if username == "ghost" {
				return nil, nil, fmt.Errorf("user not found: ghost")
			}
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int, opts poweradmin.UserDeleteOpts) (int, *poweradmin.Response, error) {
			t.Error("Delete must not be called when --transfer-to cannot be resolved")
			return 0, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewDeleteCmd(nil), []string{"--name", "max", "--transfer-to", "ghost", "--yes"})
	if err == nil || !strings.Contains(err.Error(), "--transfer-to") {
		t.Errorf("err = %v, want --transfer-to resolution error", err)
	}
}
