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

func TestGroupsZoneAdd(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		AddZoneFn: func(ctx context.Context, groupID, zoneID int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewZoneAddCmd(nil), []string{"--group-id", "1", "--zone-id", "78"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "added zone 78 to group 1") {
		t.Errorf("expected confirmation, got:\n%s", out)
	}
}

func TestGroupsZoneAddMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, &testutil.MockGroupClient{}, nil)
	err := fx.Run(groups.NewZoneAddCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestGroupsZoneAddError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		AddZoneFn: func(ctx context.Context, groupID, zoneID int) (*poweradmin.Response, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewZoneAddCmd(nil), []string{"--group-id", "1", "--zone-id", "78"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
