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

func TestUsersUpdate(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.UserUpdateOpts) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: "max", Email: "new@example.com"}, nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewUpdateCmd(nil), []string{
		"--name", "max",
		"--email", "new@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "max") {
		t.Errorf("expected output to contain max, got:\n%s", out)
	}
}

func TestUsersUpdateJSON(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.UserUpdateOpts) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: "max", Email: "new@example.com"}, nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewUpdateCmd(nil), []string{
		"--name", "max",
		"--email", "new@example.com",
		"--output", "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"username": "max"`) {
		t.Errorf("expected JSON to contain max, got:\n%s", out)
	}
}

func TestUsersUpdateMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, &testutil.MockUserClient{}, nil, nil)

	err := fx.Run(users.NewUpdateCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestUsersUpdateError(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.UserUpdateOpts) (*poweradmin.User, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewUpdateCmd(nil), []string{
		"--name", "max",
		"--email", "new@example.com",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
