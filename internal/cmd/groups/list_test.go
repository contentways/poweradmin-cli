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

func TestGroupsList(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Group, error) {
			return []*poweradmin.Group{
				{ID: 1, Name: "Administrators", Description: "Full access", MemberCount: 1},
				{ID: 2, Name: "Editors", Description: "Edit access"},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewListCmd(nil), []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "Administrators") {
		t.Errorf("expected output to contain Administrators, got:\n%s", out)
	}
}

func TestGroupsListJSON(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Group, error) {
			return []*poweradmin.Group{
				{ID: 1, Name: "Administrators"},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewListCmd(nil), []string{"-o", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, `"name": "Administrators"`) {
		t.Errorf("expected JSON to contain Administrators, got:\n%s", out)
	}
}

func TestGroupsListError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Group, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewListCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
