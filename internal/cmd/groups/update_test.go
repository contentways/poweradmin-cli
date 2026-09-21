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

func TestGroupsUpdate(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 42, Name: name}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.GroupUpdateOpts) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 42, Name: "NewName"}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewUpdateCmd(nil), []string{"--name", "TestGroup", "--new-name", "NewName"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "NewName") {
		t.Errorf("expected output to contain NewName, got:\n%s", out)
	}
}

func TestGroupsUpdateJSON(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 42, Name: name}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.GroupUpdateOpts) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 42, Name: "NewName"}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewUpdateCmd(nil), []string{"--name", "TestGroup", "--new-name", "NewName", "-o", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, `"name": "NewName"`) {
		t.Errorf("expected JSON to contain NewName, got:\n%s", out)
	}
}

func TestGroupsUpdateMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, &testutil.MockGroupClient{}, nil)
	err := fx.Run(groups.NewUpdateCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestGroupsUpdateError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 42, Name: name}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.GroupUpdateOpts) (*poweradmin.Group, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewUpdateCmd(nil), []string{"--name", "TestGroup", "--new-name", "NewName"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGroupsUpdateByID(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: id, Name: "TestGroup"}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.GroupUpdateOpts) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: id, Name: "NewName"}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewUpdateCmd(nil), []string{"--id", "42", "--new-name", "NewName"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
