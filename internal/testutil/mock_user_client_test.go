// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil_test

import (
	"context"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestMockUserClientUnsetFieldsReturnZeroValues(t *testing.T) {
	m := &testutil.MockUserClient{}
	ctx := context.Background()

	if u, resp, err := m.GetByName(ctx, "x"); u != nil || resp != nil || err != nil {
		t.Errorf("GetByName: expected all nil, got %v, %v, %v", u, resp, err)
	}
	if u, resp, err := m.GetByID(ctx, 1); u != nil || resp != nil || err != nil {
		t.Errorf("GetByID: expected all nil, got %v, %v, %v", u, resp, err)
	}
	if us, resp, err := m.List(ctx, poweradmin.ListOpts{}); us != nil || resp != nil || err != nil {
		t.Errorf("List: expected all nil, got %v, %v, %v", us, resp, err)
	}
	if us, err := m.All(ctx); us != nil || err != nil {
		t.Errorf("All: expected all nil, got %v, %v", us, err)
	}
	if id, resp, err := m.Create(ctx, poweradmin.UserCreateOpts{}); id != 0 || resp != nil || err != nil {
		t.Errorf("Create: expected zero values, got %v, %v, %v", id, resp, err)
	}
	if u, resp, err := m.Update(ctx, 1, poweradmin.UserUpdateOpts{}); u != nil || resp != nil || err != nil {
		t.Errorf("Update: expected all nil, got %v, %v, %v", u, resp, err)
	}
	if n, resp, err := m.Delete(ctx, 1, poweradmin.UserDeleteOpts{}); n != 0 || resp != nil || err != nil {
		t.Errorf("Delete: expected zero values, got %v, %v, %v", n, resp, err)
	}
	if resp, err := m.SetPermissionTemplate(ctx, 1, 2); resp != nil || err != nil {
		t.Errorf("SetPermissionTemplate: expected all nil, got %v, %v", resp, err)
	}
}

func TestMockUserClientSetFieldsDelegate(t *testing.T) {
	var calls []string
	m := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			calls = append(calls, "GetByName")
			return &poweradmin.User{Username: username}, nil, nil
		},
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.User, *poweradmin.Response, error) {
			calls = append(calls, "GetByID")
			return &poweradmin.User{ID: id}, nil, nil
		},
		ListFn: func(ctx context.Context, opts poweradmin.ListOpts) ([]*poweradmin.User, *poweradmin.Response, error) {
			calls = append(calls, "List")
			return []*poweradmin.User{{ID: 1}}, nil, nil
		},
		AllFn: func(ctx context.Context) ([]*poweradmin.User, error) {
			calls = append(calls, "All")
			return []*poweradmin.User{{ID: 1}}, nil
		},
		CreateFn: func(ctx context.Context, opts poweradmin.UserCreateOpts) (int, *poweradmin.Response, error) {
			calls = append(calls, "Create")
			return 42, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.UserUpdateOpts) (*poweradmin.User, *poweradmin.Response, error) {
			calls = append(calls, "Update")
			return &poweradmin.User{ID: id}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int, opts poweradmin.UserDeleteOpts) (int, *poweradmin.Response, error) {
			calls = append(calls, "Delete")
			return 0, nil, nil
		},
		SetPermissionTemplateFn: func(ctx context.Context, id, permTemplID int) (*poweradmin.Response, error) {
			calls = append(calls, "SetPermissionTemplate")
			return nil, nil
		},
	}
	ctx := context.Background()

	_, _, _ = m.GetByName(ctx, "x")
	_, _, _ = m.GetByID(ctx, 1)
	_, _, _ = m.List(ctx, poweradmin.ListOpts{})
	_, _ = m.All(ctx)
	_, _, _ = m.Create(ctx, poweradmin.UserCreateOpts{})
	_, _, _ = m.Update(ctx, 1, poweradmin.UserUpdateOpts{})
	_, _, _ = m.Delete(ctx, 1, poweradmin.UserDeleteOpts{})
	_, _ = m.SetPermissionTemplate(ctx, 1, 2)

	want := []string{"GetByName", "GetByID", "List", "All", "Create", "Update", "Delete", "SetPermissionTemplate"}
	if len(calls) != len(want) {
		t.Fatalf("expected %d calls, got %d: %v", len(want), len(calls), calls)
	}
	for i, name := range want {
		if calls[i] != name {
			t.Errorf("call %d: expected %s, got %s", i, name, calls[i])
		}
	}
}
