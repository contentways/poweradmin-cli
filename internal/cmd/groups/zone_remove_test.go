// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/internal/cmd/groups"
	"github.com/contentways/poweradmin-cli/internal/testutil"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestGroupsZoneRemove(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		RemoveZoneFn: func(ctx context.Context, groupID, zoneID int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewZoneRemoveCmd(nil), []string{"--group-id", "1", "--zone-id", "78"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "removed zone 78 from group 1") {
		t.Errorf("expected confirmation, got:\n%s", out)
	}
}

func TestGroupsZoneRemoveMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, &testutil.MockGroupClient{}, nil)
	err := fx.Run(groups.NewZoneRemoveCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestGroupsZoneRemoveError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		RemoveZoneFn: func(ctx context.Context, groupID, zoneID int) (*poweradmin.Response, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewZoneRemoveCmd(nil), []string{"--group-id", "1", "--zone-id", "78"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
