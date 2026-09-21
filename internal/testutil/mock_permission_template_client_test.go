// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil_test

import (
	"context"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestMockPermissionTemplateClientUnsetFieldsReturnZeroValues(t *testing.T) {
	m := &testutil.MockPermissionTemplateClient{}
	ctx := context.Background()

	if p, resp, err := m.GetByName(ctx, "x"); p != nil || resp != nil || err != nil {
		t.Errorf("GetByName: expected all nil, got %v, %v, %v", p, resp, err)
	}
	if p, resp, err := m.GetByID(ctx, 1); p != nil || resp != nil || err != nil {
		t.Errorf("GetByID: expected all nil, got %v, %v, %v", p, resp, err)
	}
	if ps, resp, err := m.List(ctx); ps != nil || resp != nil || err != nil {
		t.Errorf("List: expected all nil, got %v, %v, %v", ps, resp, err)
	}
	if p, resp, err := m.Create(ctx, poweradmin.PermissionTemplateOpts{}); p != nil || resp != nil || err != nil {
		t.Errorf("Create: expected all nil, got %v, %v, %v", p, resp, err)
	}
	if p, resp, err := m.Update(ctx, 1, poweradmin.PermissionTemplateOpts{}); p != nil || resp != nil || err != nil {
		t.Errorf("Update: expected all nil, got %v, %v, %v", p, resp, err)
	}
	if resp, err := m.Delete(ctx, 1); resp != nil || err != nil {
		t.Errorf("Delete: expected all nil, got %v, %v", resp, err)
	}
}

func TestMockPermissionTemplateClientSetFieldsDelegate(t *testing.T) {
	var calls []string
	m := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			calls = append(calls, "GetByName")
			return &poweradmin.PermissionTemplate{Name: name}, nil, nil
		},
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			calls = append(calls, "GetByID")
			return &poweradmin.PermissionTemplate{ID: id}, nil, nil
		},
		ListFn: func(ctx context.Context) ([]*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			calls = append(calls, "List")
			return []*poweradmin.PermissionTemplate{{ID: 1}}, nil, nil
		},
		CreateFn: func(ctx context.Context, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			calls = append(calls, "Create")
			return &poweradmin.PermissionTemplate{ID: 42}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			calls = append(calls, "Update")
			return &poweradmin.PermissionTemplate{ID: id}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			calls = append(calls, "Delete")
			return nil, nil
		},
	}
	ctx := context.Background()

	_, _, _ = m.GetByName(ctx, "x")
	_, _, _ = m.GetByID(ctx, 1)
	_, _, _ = m.List(ctx)
	_, _, _ = m.Create(ctx, poweradmin.PermissionTemplateOpts{})
	_, _, _ = m.Update(ctx, 1, poweradmin.PermissionTemplateOpts{})
	_, _ = m.Delete(ctx, 1)

	want := []string{"GetByName", "GetByID", "List", "Create", "Update", "Delete"}
	if len(calls) != len(want) {
		t.Fatalf("expected %d calls, got %d: %v", len(want), len(calls), calls)
	}
	for i, name := range want {
		if calls[i] != name {
			t.Errorf("call %d: expected %s, got %s", i, name, calls[i])
		}
	}
}
