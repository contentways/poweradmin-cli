// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups_test

import (
	"context"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/groups"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestGroupsMembersAddNested(t *testing.T) {
	var gotGroup, gotUser int
	mockGroup := &testutil.MockGroupClient{
		AddMemberFn: func(ctx context.Context, groupID, userID int) (*poweradmin.Response, error) {
			gotGroup, gotUser = groupID, userID
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)

	err := fx.Run(groups.NewGroupsCommand(nil), []string{"members", "add", "--group-id", "1", "--user-id", "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotGroup != 1 || gotUser != 2 {
		t.Errorf("AddMember(%d, %d), want (1, 2)", gotGroup, gotUser)
	}
}

func TestGroupsZonesRemoveNested(t *testing.T) {
	called := false
	mockGroup := &testutil.MockGroupClient{
		RemoveZoneFn: func(ctx context.Context, groupID, zoneID int) (*poweradmin.Response, error) {
			called = groupID == 1 && zoneID == 78
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)

	err := fx.Run(groups.NewGroupsCommand(nil), []string{"zones", "remove", "--group-id", "1", "--zone-id", "78"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected RemoveZone(1, 78)")
	}
}

func TestGroupsMemberAddDeprecatedAlias(t *testing.T) {
	called := false
	mockGroup := &testutil.MockGroupClient{
		AddMemberFn: func(ctx context.Context, groupID, userID int) (*poweradmin.Response, error) {
			called = true
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)

	err := fx.Run(groups.NewGroupsCommand(nil), []string{"member-add", "--group-id", "1", "--user-id", "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("alias must still add the member")
	}
	if !strings.Contains(fx.Stdout.String()+fx.Stderr.String(), `use "groups members add" instead`) {
		t.Errorf("expected deprecation note, got stdout:\n%s\nstderr:\n%s", fx.Stdout.String(), fx.Stderr.String())
	}
}
