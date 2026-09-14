// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil

import (
	"context"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
)

// MockPermissionTemplateClient implements poweradmin.IPermissionTemplateClient for testing.
// Each method can be overridden by setting the corresponding function field.
// Unset fields return zero values and no error by default.
type MockPermissionTemplateClient struct {
	GetByNameFn func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error)
	GetByIDFn   func(ctx context.Context, id int) (*poweradmin.PermissionTemplate, *poweradmin.Response, error)
	ListFn      func(ctx context.Context) ([]*poweradmin.PermissionTemplate, *poweradmin.Response, error)
	CreateFn    func(ctx context.Context, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error)
	UpdateFn    func(ctx context.Context, id int, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error)
	DeleteFn    func(ctx context.Context, id int) (*poweradmin.Response, error)
}

func (m *MockPermissionTemplateClient) GetByName(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
	if m.GetByNameFn != nil {
		return m.GetByNameFn(ctx, name)
	}
	return nil, nil, nil
}

func (m *MockPermissionTemplateClient) GetByID(ctx context.Context, id int) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil, nil
}

func (m *MockPermissionTemplateClient) List(ctx context.Context) ([]*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx)
	}
	return nil, nil, nil
}

func (m *MockPermissionTemplateClient) Create(ctx context.Context, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, opts)
	}
	return nil, nil, nil
}

func (m *MockPermissionTemplateClient) Update(ctx context.Context, id int, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, opts)
	}
	return nil, nil, nil
}

func (m *MockPermissionTemplateClient) Delete(ctx context.Context, id int) (*poweradmin.Response, error) {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil, nil
}
