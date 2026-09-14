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

func TestGroupsZones(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 1, Name: name}, nil, nil
		},
		ZonesFn: func(ctx context.Context, id int) ([]*poweradmin.GroupZone, *poweradmin.Response, error) {
			return []*poweradmin.GroupZone{
				{ZoneID: 78, ZoneName: "contentways.org", ZoneType: "NATIVE"},
			}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewZonesCmd(nil), []string{"--name", "Administrators"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "contentways.org") {
		t.Errorf("expected output to contain contentways.org, got:\n%s", out)
	}
}

func TestGroupsZonesMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, &testutil.MockGroupClient{}, nil)
	err := fx.Run(groups.NewZonesCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestGroupsZonesError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 1, Name: name}, nil, nil
		},
		ZonesFn: func(ctx context.Context, id int) ([]*poweradmin.GroupZone, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)
	err := fx.Run(groups.NewZonesCmd(nil), []string{"--name", "Administrators"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
