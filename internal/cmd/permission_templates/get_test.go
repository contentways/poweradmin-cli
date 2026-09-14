// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package permission_templates_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/permission_templates"
	"github.com/contentways/poweradmin-cli/internal/testutil"
)

func TestPermissionTemplatesGet(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: 1, Name: name, Descr: "Full rights", TemplateType: "user"}, nil, nil
		},
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: id, Name: "Administrator", Descr: "Full rights", TemplateType: "user"}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewGetCmd(nil), []string{"--name", "Administrator"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "Administrator") {
		t.Errorf("expected Administrator in output, got:\n%s", fx.Stdout.String())
	}
}

func TestPermissionTemplatesGetJSON(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: 1, Name: name, TemplateType: "user"}, nil, nil
		},
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: id, Name: "Administrator", TemplateType: "user"}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewGetCmd(nil), []string{"--name", "Administrator", "-o", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "Administrator") {
		t.Errorf("expected JSON output")
	}
}

func TestPermissionTemplatesGetMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, &testutil.MockPermissionTemplateClient{})
	err := fx.Run(permission_templates.NewGetCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestPermissionTemplatesGetError(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewGetCmd(nil), []string{"--name", "Administrator"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
