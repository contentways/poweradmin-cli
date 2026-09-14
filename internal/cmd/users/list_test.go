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

func TestUsersList(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.User, error) {
			return []*poweradmin.User{
				{ID: 1, Username: "max", Email: "max@example.com", Active: true},
				{ID: 2, Username: "markus", Email: "markus@example.com", Active: false},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewListCmd(), []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "max") {
		t.Errorf("expected output to contain max, got:\n%s", out)
	}
	if !strings.Contains(out, "markus") {
		t.Errorf("expected output to contain markus, got:\n%s", out)
	}
}

func TestUsersListJSON(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.User, error) {
			return []*poweradmin.User{
				{ID: 1, Username: "max", Email: "max@example.com", Active: true},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewListCmd(), []string{"--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"username": "max"`) {
		t.Errorf("expected JSON to contain max, got:\n%s", out)
	}
}

func TestUsersListError(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.User, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewListCmd(), []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
