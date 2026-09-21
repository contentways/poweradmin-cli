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
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestGroupsCreate(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		CreateFn: func(ctx context.Context, opts poweradmin.GroupCreateOpts) (int, *poweradmin.Response, error) {
			return 42, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewCreateCmd(), []string{"--name", "TestGroup"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "TestGroup") {
		t.Errorf("expected output to contain TestGroup, got:\n%s", out)
	}
}

func TestGroupsCreateJSON(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		CreateFn: func(ctx context.Context, opts poweradmin.GroupCreateOpts) (int, *poweradmin.Response, error) {
			return 42, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewCreateCmd(), []string{"--name", "TestGroup", "-o", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, `"id": 42`) {
		t.Errorf("expected JSON to contain id 42, got:\n%s", out)
	}
}

func TestGroupsCreateMissingName(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, &testutil.MockGroupClient{}, nil)
	err := fx.Run(groups.NewCreateCmd(), []string{})
	if err == nil {
		t.Fatal("expected error when name is missing")
	}
}

func TestGroupsCreateError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		CreateFn: func(ctx context.Context, opts poweradmin.GroupCreateOpts) (int, *poweradmin.Response, error) {
			return 0, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewCreateCmd(), []string{"--name", "TestGroup"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestGroupsCreateInteractiveSkipsPromptsWhenAllFlagsSet verifies that when
// --interactive is combined with all required values already supplied via
// flags (--description and --perm-template-id explicitly set so those huh
// prompts don't fire), no huh prompt runs at all — only the final
// confirmation, answered via a piped stdin.
func TestGroupsCreateInteractiveSkipsPromptsWhenAllFlagsSet(t *testing.T) {
	var createdName, createdDescription string
	var createdPermTemplID int

	mockGroup := &testutil.MockGroupClient{
		CreateFn: func(ctx context.Context, opts poweradmin.GroupCreateOpts) (int, *poweradmin.Response, error) {
			createdName = opts.Name
			createdDescription = opts.Description
			createdPermTemplID = opts.PermTemplID
			return 42, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	testutil.WithStdin(t, "y\n")

	err := fx.Run(groups.NewCreateCmd(), []string{
		"--name", "TestGroup",
		"--description", "A test group",
		"--perm-template-id", "3",
		"--interactive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if createdName != "TestGroup" {
		t.Errorf("expected name TestGroup, got %q", createdName)
	}
	if createdDescription != "A test group" {
		t.Errorf("expected description 'A test group', got %q", createdDescription)
	}
	if createdPermTemplID != 3 {
		t.Errorf("expected perm-template-id 3, got %d", createdPermTemplID)
	}
}

// TestGroupsCreateInteractiveDeclineAbortsWithoutCreating verifies that
// declining the final confirmation prevents the API call entirely.
func TestGroupsCreateInteractiveDeclineAbortsWithoutCreating(t *testing.T) {
	created := false
	mockGroup := &testutil.MockGroupClient{
		CreateFn: func(ctx context.Context, opts poweradmin.GroupCreateOpts) (int, *poweradmin.Response, error) {
			created = true
			return 42, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	testutil.WithStdin(t, "n\n")

	err := fx.Run(groups.NewCreateCmd(), []string{
		"--name", "TestGroup",
		"--description", "A test group",
		"--perm-template-id", "3",
		"--interactive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created {
		t.Error("expected group NOT to be created after declining confirmation")
	}
}
