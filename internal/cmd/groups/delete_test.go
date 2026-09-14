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
