// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil_test

import (
	"context"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestMockGroupClientUnsetFieldsReturnZeroValues(t *testing.T) {
	m := &testutil.MockGroupClient{}
	ctx := context.Background()

	if g, resp, err := m.GetByName(ctx, "x"); g != nil || resp != nil || err != nil {
		t.Errorf("GetByName: expected all nil, got %v, %v, %v", g, resp, err)
	}
	if g, resp, err := m.GetByID(ctx, 1); g != nil || resp != nil || err != nil {
		t.Errorf("GetByID: expected all nil, got %v, %v, %v", g, resp, err)
	}
	if gs, resp, err := m.List(ctx, poweradmin.ListOpts{}); gs != nil || resp != nil || err != nil {
		t.Errorf("List: expected all nil, got %v, %v, %v", gs, resp, err)
	}
	if gs, err := m.All(ctx); gs != nil || err != nil {
		t.Errorf("All: expected all nil, got %v, %v", gs, err)
	}
	if id, resp, err := m.Create(ctx, poweradmin.GroupCreateOpts{}); id != 0 || resp != nil || err != nil {
		t.Errorf("Create: expected zero values, got %v, %v, %v", id, resp, err)
	}
	if g, resp, err := m.Update(ctx, 1, poweradmin.GroupUpdateOpts{}); g != nil || resp != nil || err != nil {
		t.Errorf("Update: expected all nil, got %v, %v, %v", g, resp, err)
	}
	if resp, err := m.Delete(ctx, 1); resp != nil || err != nil {
		t.Errorf("Delete: expected all nil, got %v, %v", resp, err)
	}
	if ms, resp, err := m.Members(ctx, 1); ms != nil || resp != nil || err != nil {
		t.Errorf("Members: expected all nil, got %v, %v, %v", ms, resp, err)
	}
	if resp, err := m.AddMember(ctx, 1, 2); resp != nil || err != nil {
		t.Errorf("AddMember: expected all nil, got %v, %v", resp, err)
	}
	if resp, err := m.RemoveMember(ctx, 1, 2); resp != nil || err != nil {
		t.Errorf("RemoveMember: expected all nil, got %v, %v", resp, err)
	}
	if zs, resp, err := m.Zones(ctx, 1); zs != nil || resp != nil || err != nil {
		t.Errorf("Zones: expected all nil, got %v, %v, %v", zs, resp, err)
	}
	if resp, err := m.AddZone(ctx, 1, 78); resp != nil || err != nil {
		t.Errorf("AddZone: expected all nil, got %v, %v", resp, err)
	}
	if resp, err := m.RemoveZone(ctx, 1, 78); resp != nil || err != nil {
		t.Errorf("RemoveZone: expected all nil, got %v, %v", resp, err)
	}
}

func TestMockGroupClientSetFieldsDelegate(t *testing.T) {
	var calls []string
	m := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			calls = append(calls, "GetByName")
			return &poweradmin.Group{Name: name}, nil, nil
		},
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.Group, *poweradmin.Response, error) {
			calls = append(calls, "GetByID")
			return &poweradmin.Group{ID: id}, nil, nil
		},
		ListFn: func(ctx context.Context, opts poweradmin.ListOpts) ([]*poweradmin.Group, *poweradmin.Response, error) {
			calls = append(calls, "List")
			return []*poweradmin.Group{{ID: 1}}, nil, nil
		},
		AllFn: func(ctx context.Context) ([]*poweradmin.Group, error) {
			calls = append(calls, "All")
			return []*poweradmin.Group{{ID: 1}}, nil
		},
		CreateFn: func(ctx context.Context, opts poweradmin.GroupCreateOpts) (int, *poweradmin.Response, error) {
			calls = append(calls, "Create")
			return 42, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.GroupUpdateOpts) (*poweradmin.Group, *poweradmin.Response, error) {
			calls = append(calls, "Update")
			return &poweradmin.Group{ID: id}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			calls = append(calls, "Delete")
			return nil, nil
		},
		MembersFn: func(ctx context.Context, id int) ([]*poweradmin.GroupMember, *poweradmin.Response, error) {
			calls = append(calls, "Members")
			return []*poweradmin.GroupMember{{UserID: 1}}, nil, nil
		},
		AddMemberFn: func(ctx context.Context, groupID, userID int) (*poweradmin.Response, error) {
			calls = append(calls, "AddMember")
			return nil, nil
		},
		RemoveMemberFn: func(ctx context.Context, groupID, userID int) (*poweradmin.Response, error) {
			calls = append(calls, "RemoveMember")
			return nil, nil
		},
		ZonesFn: func(ctx context.Context, id int) ([]*poweradmin.GroupZone, *poweradmin.Response, error) {
			calls = append(calls, "Zones")
			return []*poweradmin.GroupZone{{ZoneID: 78}}, nil, nil
		},
		AddZoneFn: func(ctx context.Context, groupID, zoneID int) (*poweradmin.Response, error) {
			calls = append(calls, "AddZone")
			return nil, nil
		},
		RemoveZoneFn: func(ctx context.Context, groupID, zoneID int) (*poweradmin.Response, error) {
			calls = append(calls, "RemoveZone")
			return nil, nil
		},
	}
	ctx := context.Background()

	_, _, _ = m.GetByName(ctx, "x")
	_, _, _ = m.GetByID(ctx, 1)
	_, _, _ = m.List(ctx, poweradmin.ListOpts{})
	_, _ = m.All(ctx)
	_, _, _ = m.Create(ctx, poweradmin.GroupCreateOpts{})
	_, _, _ = m.Update(ctx, 1, poweradmin.GroupUpdateOpts{})
	_, _ = m.Delete(ctx, 1)
	_, _, _ = m.Members(ctx, 1)
	_, _ = m.AddMember(ctx, 1, 2)
	_, _ = m.RemoveMember(ctx, 1, 2)
	_, _, _ = m.Zones(ctx, 1)
	_, _ = m.AddZone(ctx, 1, 78)
	_, _ = m.RemoveZone(ctx, 1, 78)

	want := []string{
		"GetByName", "GetByID", "List", "All", "Create", "Update", "Delete",
		"Members", "AddMember", "RemoveMember", "Zones", "AddZone", "RemoveZone",
	}
	if len(calls) != len(want) {
		t.Fatalf("expected %d calls, got %d: %v", len(want), len(calls), calls)
	}
	for i, name := range want {
		if calls[i] != name {
			t.Errorf("call %d: expected %s, got %s", i, name, calls[i])
		}
	}
}
