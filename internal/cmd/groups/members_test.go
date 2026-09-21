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

func TestGroupsMembers(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 1, Name: name}, nil, nil
		},
		MembersFn: func(ctx context.Context, id int) ([]*poweradmin.GroupMember, *poweradmin.Response, error) {
			return []*poweradmin.GroupMember{
				{UserID: 1, Username: "patrick"},
			}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewMembersCmd(nil), []string{"--name", "Administrators"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "patrick") {
		t.Errorf("expected output to contain patrick, got:\n%s", out)
	}
}

func TestGroupsMembersMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, &testutil.MockGroupClient{}, nil)
	err := fx.Run(groups.NewMembersCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestGroupsMembersError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 1, Name: name}, nil, nil
		},
		MembersFn: func(ctx context.Context, id int) ([]*poweradmin.GroupMember, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewMembersCmd(nil), []string{"--name", "Administrators"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGroupsMembersByID(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		MembersFn: func(ctx context.Context, id int) ([]*poweradmin.GroupMember, *poweradmin.Response, error) {
			return []*poweradmin.GroupMember{{UserID: 1, Username: "patrick"}}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewMembersCmd(nil), []string{"--id", "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGroupsMembersJSON(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 1, Name: name}, nil, nil
		},
		MembersFn: func(ctx context.Context, id int) ([]*poweradmin.GroupMember, *poweradmin.Response, error) {
			return []*poweradmin.GroupMember{{UserID: 1, Username: "patrick"}}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewMembersCmd(nil), []string{"--name", "Administrators", "-o", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "patrick") {
		t.Errorf("expected JSON to contain patrick, got:\n%s", out)
	}
}

func TestGroupsMembersNilWithoutError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return nil, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewMembersCmd(nil), []string{"--name", "Ghosts"})
	if err == nil {
		t.Fatal("expected error for nil group with nil error, got nil")
	}
}
