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
	"github.com/contentways/poweradmin-go/v3/poweradmin"
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

// TestUsersCreateInteractiveSkipsPromptsWhenAllFlagsSet verifies that when
// --interactive is combined with all required values already supplied via
// flags (including --password, so the masked term.ReadPassword prompt never
// fires, and --fullname/--active so those huh prompts don't fire either),
// no huh prompt runs at all — only the final confirmation, answered via a
// piped stdin.
func TestUsersCreateInteractiveSkipsPromptsWhenAllFlagsSet(t *testing.T) {
	var createdUsername, createdEmail string

	mockUser := &testutil.MockUserClient{
		CreateFn: func(ctx context.Context, opts poweradmin.UserCreateOpts) (int, *poweradmin.Response, error) {
			createdUsername = opts.Username
			createdEmail = opts.Email
			return 42, nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)
	testutil.WithStdin(t, "y\n")

	err := fx.Run(users.NewCreateCmd(), []string{
		"--username", "max",
		"--password", "secret123",
		"--email", "max@example.com",
		"--fullname", "Max Mustermann",
		"--active",
		"--interactive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if createdUsername != "max" {
		t.Errorf("expected username max, got %q", createdUsername)
	}
	if createdEmail != "max@example.com" {
		t.Errorf("expected email max@example.com, got %q", createdEmail)
	}
}

// TestUsersCreateInteractiveDeclineAbortsWithoutCreating verifies that
// declining the final confirmation prevents the API call entirely.
func TestUsersCreateInteractiveDeclineAbortsWithoutCreating(t *testing.T) {
	created := false
	mockUser := &testutil.MockUserClient{
		CreateFn: func(ctx context.Context, opts poweradmin.UserCreateOpts) (int, *poweradmin.Response, error) {
			created = true
			return 42, nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)
	testutil.WithStdin(t, "n\n")

	err := fx.Run(users.NewCreateCmd(), []string{
		"--username", "max",
		"--password", "secret123",
		"--email", "max@example.com",
		"--fullname", "Max Mustermann",
		"--active",
		"--interactive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created {
		t.Error("expected user NOT to be created after declining confirmation")
	}
}

func TestUsersCreateQuiet(t *testing.T) {
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
		"--quiet",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(fx.Stdout.String())
	if out != "42" {
		t.Errorf("expected quiet output to be just the id, got: %q", out)
	}
}

func TestUsersCreateExplicitlyInactive(t *testing.T) {
	var createdActive bool
	mockUser := &testutil.MockUserClient{
		CreateFn: func(ctx context.Context, opts poweradmin.UserCreateOpts) (int, *poweradmin.Response, error) {
			createdActive = opts.Active
			return 42, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewCreateCmd(), []string{
		"--username", "max",
		"--password", "secret123",
		"--email", "max@example.com",
		"--active=false",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if createdActive {
		t.Error("expected active to be false")
	}
}

// TestUsersCreateInteractiveFullFlow drives the huh prompts for username,
// email, fullname, and active via accessible mode — --password is still
// supplied via flag, since the masked password prompt uses term.ReadPassword
// directly against a real terminal file descriptor and cannot be driven
// through a piped stdin.
func TestUsersCreateInteractiveFullFlow(t *testing.T) {
	testutil.WithAccessiblePrompts(t)

	var createdUsername, createdEmail, createdFullname string
	var createdActive bool
	mockUser := &testutil.MockUserClient{
		CreateFn: func(ctx context.Context, opts poweradmin.UserCreateOpts) (int, *poweradmin.Response, error) {
			createdUsername = opts.Username
			createdEmail = opts.Email
			createdFullname = opts.Fullname
			createdActive = opts.Active
			return 42, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	// username, email, fullname, active (y), then the final confirmation.
	testutil.WithDelayedStdin(t,
		"max\n",
		"max@example.com\n",
		"Max Mustermann\n",
		"y\n",
		"y\n",
	)

	err := fx.Run(users.NewCreateCmd(), []string{
		"--password", "secret123",
		"--interactive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if createdUsername != "max" {
		t.Errorf("expected username max, got %q", createdUsername)
	}
	if createdEmail != "max@example.com" {
		t.Errorf("expected email max@example.com, got %q", createdEmail)
	}
	if createdFullname != "Max Mustermann" {
		t.Errorf("expected fullname Max Mustermann, got %q", createdFullname)
	}
	if !createdActive {
		t.Error("expected active to be true")
	}
}
