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

func TestUsersListFilterActiveTrue(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.User, error) {
			return []*poweradmin.User{
				{ID: 1, Username: "max", Active: true},
				{ID: 2, Username: "markus", Active: false},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewListCmd(), []string{"--active"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "max") {
		t.Errorf("expected active user max in output, got:\n%s", out)
	}
	if strings.Contains(out, "markus") {
		t.Errorf("expected inactive user markus to be filtered out, got:\n%s", out)
	}
}

func TestUsersListFilterActiveFalse(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.User, error) {
			return []*poweradmin.User{
				{ID: 1, Username: "max", Active: true},
				{ID: 2, Username: "markus", Active: false},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewListCmd(), []string{"--active=false"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if strings.Contains(out, "max") {
		t.Errorf("expected active user max to be filtered out, got:\n%s", out)
	}
	if !strings.Contains(out, "markus") {
		t.Errorf("expected inactive user markus in output, got:\n%s", out)
	}
}

func TestUsersListSortByUsername(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.User, error) {
			return []*poweradmin.User{
				{ID: 1, Username: "zack"},
				{ID: 2, Username: "amy"},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewListCmd(), []string{"--sort", "username"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	amyIdx := strings.Index(out, "amy")
	zackIdx := strings.Index(out, "zack")
	if amyIdx == -1 || zackIdx == -1 || amyIdx > zackIdx {
		t.Errorf("expected amy to sort before zack, got:\n%s", out)
	}
}
