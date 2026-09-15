// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v2/internal/cmd/users"
	"github.com/contentways/poweradmin-cli/v2/internal/testutil"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestUsersGetByName(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 1, Username: username, Email: "max@example.com", Active: true}, nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewGetCmd(nil), []string{"--name", "max"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "max") {
		t.Errorf("expected output to contain max, got:\n%s", out)
	}
}

func TestUsersGetJSON(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 1, Username: username, Email: "max@example.com", Active: true}, nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewGetCmd(nil), []string{"--name", "max", "--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"username": "max"`) {
		t.Errorf("expected JSON to contain max, got:\n%s", out)
	}
}

func TestUsersGetMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, &testutil.MockUserClient{}, nil, nil)

	err := fx.Run(users.NewGetCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestUsersGetError(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewGetCmd(nil), []string{"--name", "max"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
