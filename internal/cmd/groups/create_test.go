// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/groups"
	"github.com/contentways/poweradmin-cli/internal/testutil"
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
