// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package permission_templates_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/permission_templates"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestPermissionTemplatesUpdate(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: 1, Name: name, TemplateType: "user"}, nil, nil
		},
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: id, Name: "Administrator", TemplateType: "user"}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: id, Name: opts.Name, TemplateType: opts.TemplateType}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewUpdateCmd(nil), []string{"--name", "Administrator", "--new-name", "SuperAdmin"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "updated") {
		t.Errorf("expected updated in output, got:\n%s", fx.Stdout.String())
	}
}

func TestPermissionTemplatesUpdateMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, &testutil.MockPermissionTemplateClient{})
	err := fx.Run(permission_templates.NewUpdateCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestPermissionTemplatesUpdateError(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewUpdateCmd(nil), []string{"--name", "Administrator", "--new-name", "SuperAdmin"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPermissionTemplatesUpdateByID(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: id, Name: "Administrator", TemplateType: "user"}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: id, Name: opts.Name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewUpdateCmd(nil), []string{"--id", "1", "--new-name", "SuperAdmin"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPermissionTemplatesUpdateJSON(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: 1, Name: name, TemplateType: "user"}, nil, nil
		},
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: id, Name: "SuperAdmin", TemplateType: "user"}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: id, Name: opts.Name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewUpdateCmd(nil), []string{"--name", "Administrator", "--new-name", "SuperAdmin", "-o", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "SuperAdmin") {
		t.Errorf("expected JSON to contain SuperAdmin, got:\n%s", out)
	}
}

func TestPermissionTemplatesUpdatePermissions(t *testing.T) {
	var updatedPerms []int
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: 1, Name: name, TemplateType: "user"}, nil, nil
		},
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: id, Name: "Administrator", TemplateType: "user"}, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			updatedPerms = opts.Permissions
			return &poweradmin.PermissionTemplate{ID: id, Name: opts.Name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewUpdateCmd(nil), []string{"--name", "Administrator", "--permissions", "1,2,3,4"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updatedPerms) != 4 {
		t.Errorf("expected 4 permissions, got %v", updatedPerms)
	}
}

func TestPermissionTemplatesUpdateNilWithoutError(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return nil, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewUpdateCmd(nil), []string{"--id", "1", "--new-name", "SuperAdmin"})
	if err == nil {
		t.Fatal("expected error for nil template with nil error, got nil")
	}
}
