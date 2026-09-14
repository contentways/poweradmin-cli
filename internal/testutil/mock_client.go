// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil

import (
	"context"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
)

// MockZoneClient implements poweradmin.IZoneClient for testing.
// Each method can be overridden by setting the corresponding function field.
// Unset fields return zero values and no error by default.
type MockZoneClient struct {
	GetByIDFn     func(ctx context.Context, id int) (*poweradmin.Zone, *poweradmin.Response, error)
	GetByNameFn   func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error)
	ListFn        func(ctx context.Context, opts poweradmin.ListOpts) ([]*poweradmin.Zone, *poweradmin.Response, error)
	AllFn         func(ctx context.Context) ([]*poweradmin.Zone, error)
	CreateFn      func(ctx context.Context, opts poweradmin.ZoneCreateOpts) (int, *poweradmin.Response, error)
	UpdateFn      func(ctx context.Context, id int, opts poweradmin.ZoneUpdateOpts) (*poweradmin.Zone, *poweradmin.Response, error)
	DeleteFn      func(ctx context.Context, id int) (*poweradmin.Response, error)
	OwnersFn      func(ctx context.Context, zoneID int) ([]*poweradmin.ZoneOwner, *poweradmin.Response, error)
	AddOwnerFn    func(ctx context.Context, zoneID, userID int) (*poweradmin.Response, error)
	AddOwnersFn   func(ctx context.Context, zoneID int, userIDs []int) (*poweradmin.Response, error)
	RemoveOwnerFn func(ctx context.Context, zoneID, userID int) (*poweradmin.Response, error)
}

func (m *MockZoneClient) GetByID(ctx context.Context, id int) (*poweradmin.Zone, *poweradmin.Response, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil, nil
}

func (m *MockZoneClient) GetByName(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
	if m.GetByNameFn != nil {
		return m.GetByNameFn(ctx, name)
	}
	return nil, nil, nil
}

func (m *MockZoneClient) List(ctx context.Context, opts poweradmin.ListOpts) ([]*poweradmin.Zone, *poweradmin.Response, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, opts)
	}
	return nil, nil, nil
}

func (m *MockZoneClient) All(ctx context.Context) ([]*poweradmin.Zone, error) {
	if m.AllFn != nil {
		return m.AllFn(ctx)
	}
	return nil, nil
}

func (m *MockZoneClient) Create(ctx context.Context, opts poweradmin.ZoneCreateOpts) (int, *poweradmin.Response, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, opts)
	}
	return 0, nil, nil
}

func (m *MockZoneClient) Update(ctx context.Context, id int, opts poweradmin.ZoneUpdateOpts) (*poweradmin.Zone, *poweradmin.Response, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, opts)
	}
	return nil, nil, nil
}

func (m *MockZoneClient) Delete(ctx context.Context, id int) (*poweradmin.Response, error) {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil, nil
}

func (m *MockZoneClient) Owners(ctx context.Context, zoneID int) ([]*poweradmin.ZoneOwner, *poweradmin.Response, error) {
	if m.OwnersFn != nil {
		return m.OwnersFn(ctx, zoneID)
	}
	return nil, nil, nil
}

func (m *MockZoneClient) AddOwner(ctx context.Context, zoneID, userID int) (*poweradmin.Response, error) {
	if m.AddOwnerFn != nil {
		return m.AddOwnerFn(ctx, zoneID, userID)
	}
	return nil, nil
}

func (m *MockZoneClient) AddOwners(ctx context.Context, zoneID int, userIDs []int) (*poweradmin.Response, error) {
	if m.AddOwnersFn != nil {
		return m.AddOwnersFn(ctx, zoneID, userIDs)
	}
	return nil, nil
}

func (m *MockZoneClient) RemoveOwner(ctx context.Context, zoneID, userID int) (*poweradmin.Response, error) {
	if m.RemoveOwnerFn != nil {
		return m.RemoveOwnerFn(ctx, zoneID, userID)
	}
	return nil, nil
}
