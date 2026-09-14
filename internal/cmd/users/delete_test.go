// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/users"
	"github.com/contentways/poweradmin-cli/internal/testutil"
)

func TestUsersDelete(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username, Email: "max@example.com"}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
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
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
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
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewDeleteCmd(nil), []string{"--name", "max", "--yes"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
