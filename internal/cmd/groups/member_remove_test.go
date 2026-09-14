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

func TestGroupsMemberRemove(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		RemoveMemberFn: func(ctx context.Context, groupID, userID int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewMemberRemoveCmd(nil), []string{"--group-id", "1", "--user-id", "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "removed user 2 from group 1") {
		t.Errorf("expected confirmation, got:\n%s", out)
	}
}

func TestGroupsMemberRemoveMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, &testutil.MockGroupClient{}, nil)
	err := fx.Run(groups.NewMemberRemoveCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestGroupsMemberRemoveError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		RemoveMemberFn: func(ctx context.Context, groupID, userID int) (*poweradmin.Response, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewMemberRemoveCmd(nil), []string{"--group-id", "1", "--user-id", "2"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
