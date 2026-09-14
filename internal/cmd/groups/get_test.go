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

func TestGroupsGetByName(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 1, Name: name, Description: "Full access"}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewGetCmd(nil), []string{"--name", "Administrators"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "Administrators") {
		t.Errorf("expected output to contain Administrators, got:\n%s", out)
	}
}

func TestGroupsGetJSON(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 1, Name: name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewGetCmd(nil), []string{"--name", "Administrators", "-o", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, `"name": "Administrators"`) {
		t.Errorf("expected JSON to contain Administrators, got:\n%s", out)
	}
}

func TestGroupsGetMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, &testutil.MockGroupClient{}, nil)
	err := fx.Run(groups.NewGetCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestGroupsGetError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewGetCmd(nil), []string{"--name", "Administrators"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
