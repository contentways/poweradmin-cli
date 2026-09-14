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

func TestUsersCreate(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		CreateFn: func(ctx context.Context, opts poweradmin.UserCreateOpts) (int, *poweradmin.Response, error) {
			return 42, nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewCreateCmd(), []string{
		"--username", "max",
		"--password", "secret123",
		"--email", "max@example.com",
	})
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

func TestUsersCreateJSON(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		CreateFn: func(ctx context.Context, opts poweradmin.UserCreateOpts) (int, *poweradmin.Response, error) {
			return 42, nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewCreateCmd(), []string{
		"--username", "max",
		"--password", "secret123",
		"--email", "max@example.com",
		"--output", "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"username": "max"`) {
		t.Errorf("expected JSON to contain max, got:\n%s", out)
	}
	if !strings.Contains(out, `"id": 42`) {
		t.Errorf("expected JSON to contain id 42, got:\n%s", out)
	}
}

func TestUsersCreateMissingUsername(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, &testutil.MockUserClient{}, nil, nil)

	err := fx.Run(users.NewCreateCmd(), []string{
		"--password", "secret123",
		"--email", "max@example.com",
	})
	if err == nil {
		t.Fatal("expected error when username is missing")
	}
}

func TestUsersCreateError(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		CreateFn: func(ctx context.Context, opts poweradmin.UserCreateOpts) (int, *poweradmin.Response, error) {
			return 0, nil, fmt.Errorf("api error")
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewCreateCmd(), []string{
		"--username", "max",
		"--password", "secret123",
		"--email", "max@muster.com",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
