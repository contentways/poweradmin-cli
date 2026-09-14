// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil

import (
	"context"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
)

// MockUserClient implements poweradmin.IUserClient for testing.
// Each method can be overridden by setting the corresponding function field.
// Unset fields return zero values and no error by default.
type MockUserClient struct {
	GetByNameFn             func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error)
	GetByIDFn               func(ctx context.Context, id int) (*poweradmin.User, *poweradmin.Response, error)
	ListFn                  func(ctx context.Context, opts poweradmin.ListOpts) ([]*poweradmin.User, *poweradmin.Response, error)
	AllFn                   func(ctx context.Context) ([]*poweradmin.User, error)
	CreateFn                func(ctx context.Context, opts poweradmin.UserCreateOpts) (int, *poweradmin.Response, error)
	UpdateFn                func(ctx context.Context, id int, opts poweradmin.UserUpdateOpts) (*poweradmin.User, *poweradmin.Response, error)
	DeleteFn                func(ctx context.Context, id int) (*poweradmin.Response, error)
	SetPermissionTemplateFn func(ctx context.Context, id, permTemplID int) (*poweradmin.Response, error)
}

func (m *MockUserClient) GetByName(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
	if m.GetByNameFn != nil {
		return m.GetByNameFn(ctx, username)
	}
	return nil, nil, nil
}

func (m *MockUserClient) GetByID(ctx context.Context, id int) (*poweradmin.User, *poweradmin.Response, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil, nil
}

func (m *MockUserClient) List(ctx context.Context, opts poweradmin.ListOpts) ([]*poweradmin.User, *poweradmin.Response, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, opts)
	}
	return nil, nil, nil
}

func (m *MockUserClient) All(ctx context.Context) ([]*poweradmin.User, error) {
	if m.AllFn != nil {
		return m.AllFn(ctx)
	}
	return nil, nil
}

func (m *MockUserClient) Create(ctx context.Context, opts poweradmin.UserCreateOpts) (int, *poweradmin.Response, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, opts)
	}
	return 0, nil, nil
}

func (m *MockUserClient) Update(ctx context.Context, id int, opts poweradmin.UserUpdateOpts) (*poweradmin.User, *poweradmin.Response, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, opts)
	}
	return nil, nil, nil
}

func (m *MockUserClient) Delete(ctx context.Context, id int) (*poweradmin.Response, error) {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil, nil
}

func (m *MockUserClient) SetPermissionTemplate(ctx context.Context, id, permTemplID int) (*poweradmin.Response, error) {
	if m.SetPermissionTemplateFn != nil {
		return m.SetPermissionTemplateFn(ctx, id, permTemplID)
	}
	return nil, nil
}
