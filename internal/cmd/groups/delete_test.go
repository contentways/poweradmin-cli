// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/groups"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestGroupsDelete(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 42, Name: name}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewDeleteCmd(nil), []string{"--name", "TestGroup", "--yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "42") {
		t.Errorf("expected output to contain id 42, got:\n%s", out)
	}
}

func TestGroupsDeleteJSON(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 42, Name: name}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewDeleteCmd(nil), []string{"--name", "TestGroup", "--yes", "-o", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, `"id": 42`) {
		t.Errorf("expected JSON to contain id 42, got:\n%s", out)
	}
}

func TestGroupsDeleteMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, &testutil.MockGroupClient{}, nil)
	err := fx.Run(groups.NewDeleteCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestGroupsDeleteError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 42, Name: name}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewDeleteCmd(nil), []string{"--name", "TestGroup", "--yes"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGroupsDeleteByID(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: id, Name: "TestGroup"}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewDeleteCmd(nil), []string{"--id", "42", "--yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGroupsDeleteQuiet(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 42, Name: name}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewDeleteCmd(nil), []string{"--name", "TestGroup", "--yes", "-q"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fx.Stdout.String() != "" {
		t.Errorf("expected no output in quiet mode, got:\n%s", fx.Stdout.String())
	}
}

// TestGroupsDeleteDryRun verifies that --dry-run prints what would be
// deleted without calling Delete or prompting for confirmation.
func TestGroupsDeleteDryRun(t *testing.T) {
	deleteCalled := false
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 42, Name: name}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			deleteCalled = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	// No stdin provided — dry-run must never prompt.

	err := fx.Run(groups.NewDeleteCmd(nil), []string{"--name", "TestGroup", "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleteCalled {
		t.Error("expected Delete to never be called with --dry-run")
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "Would delete group TestGroup (id 42)") {
		t.Errorf("expected dry-run message, got:\n%s", out)
	}
}
