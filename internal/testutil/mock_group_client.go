// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil

import (
	"context"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
)

// MockGroupClient implements poweradmin.IGroupClient for testing.
// Each method can be overridden by setting the corresponding function field.
// Unset fields return zero values and no error by default.
type MockGroupClient struct {
	GetByNameFn    func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error)
	GetByIDFn      func(ctx context.Context, id int) (*poweradmin.Group, *poweradmin.Response, error)
	ListFn         func(ctx context.Context, opts poweradmin.ListOpts) ([]*poweradmin.Group, *poweradmin.Response, error)
	AllFn          func(ctx context.Context) ([]*poweradmin.Group, error)
	CreateFn       func(ctx context.Context, opts poweradmin.GroupCreateOpts) (int, *poweradmin.Response, error)
	UpdateFn       func(ctx context.Context, id int, opts poweradmin.GroupUpdateOpts) (*poweradmin.Group, *poweradmin.Response, error)
	DeleteFn       func(ctx context.Context, id int) (*poweradmin.Response, error)
	MembersFn      func(ctx context.Context, id int) ([]*poweradmin.GroupMember, *poweradmin.Response, error)
	AddMemberFn    func(ctx context.Context, groupID, userID int) (*poweradmin.Response, error)
	RemoveMemberFn func(ctx context.Context, groupID, userID int) (*poweradmin.Response, error)
	ZonesFn        func(ctx context.Context, id int) ([]*poweradmin.GroupZone, *poweradmin.Response, error)
	AddZoneFn      func(ctx context.Context, groupID, zoneID int) (*poweradmin.Response, error)
	RemoveZoneFn   func(ctx context.Context, groupID, zoneID int) (*poweradmin.Response, error)
}

func (m *MockGroupClient) GetByName(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
	if m.GetByNameFn != nil {
		return m.GetByNameFn(ctx, name)
	}
	return nil, nil, nil
}

func (m *MockGroupClient) GetByID(ctx context.Context, id int) (*poweradmin.Group, *poweradmin.Response, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil, nil
}

func (m *MockGroupClient) List(ctx context.Context, opts poweradmin.ListOpts) ([]*poweradmin.Group, *poweradmin.Response, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, opts)
	}
	return nil, nil, nil
}

func (m *MockGroupClient) All(ctx context.Context) ([]*poweradmin.Group, error) {
	if m.AllFn != nil {
		return m.AllFn(ctx)
	}
	return nil, nil
}

func (m *MockGroupClient) Create(ctx context.Context, opts poweradmin.GroupCreateOpts) (int, *poweradmin.Response, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, opts)
	}
	return 0, nil, nil
}

func (m *MockGroupClient) Update(ctx context.Context, id int, opts poweradmin.GroupUpdateOpts) (*poweradmin.Group, *poweradmin.Response, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, opts)
	}
	return nil, nil, nil
}

func (m *MockGroupClient) Delete(ctx context.Context, id int) (*poweradmin.Response, error) {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil, nil
}

func (m *MockGroupClient) Members(ctx context.Context, id int) ([]*poweradmin.GroupMember, *poweradmin.Response, error) {
	if m.MembersFn != nil {
		return m.MembersFn(ctx, id)
	}
	return nil, nil, nil
}

func (m *MockGroupClient) AddMember(ctx context.Context, groupID, userID int) (*poweradmin.Response, error) {
	if m.AddMemberFn != nil {
		return m.AddMemberFn(ctx, groupID, userID)
	}
	return nil, nil
}

func (m *MockGroupClient) RemoveMember(ctx context.Context, groupID, userID int) (*poweradmin.Response, error) {
	if m.RemoveMemberFn != nil {
		return m.RemoveMemberFn(ctx, groupID, userID)
	}
	return nil, nil
}

func (m *MockGroupClient) Zones(ctx context.Context, id int) ([]*poweradmin.GroupZone, *poweradmin.Response, error) {
	if m.ZonesFn != nil {
		return m.ZonesFn(ctx, id)
	}
	return nil, nil, nil
}

func (m *MockGroupClient) AddZone(ctx context.Context, groupID, zoneID int) (*poweradmin.Response, error) {
	if m.AddZoneFn != nil {
		return m.AddZoneFn(ctx, groupID, zoneID)
	}
	return nil, nil
}

func (m *MockGroupClient) RemoveZone(ctx context.Context, groupID, zoneID int) (*poweradmin.Response, error) {
	if m.RemoveZoneFn != nil {
		return m.RemoveZoneFn(ctx, groupID, zoneID)
	}
	return nil, nil
}
