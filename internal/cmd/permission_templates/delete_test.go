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
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestPermissionTemplatesDelete(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: 1, Name: name}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewDeleteCmd(nil), []string{"--name", "MyTemplate", "--yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "deleted") {
		t.Errorf("expected deleted in output, got:\n%s", fx.Stdout.String())
	}
}

func TestPermissionTemplatesDeleteJSON(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: 1, Name: name}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewDeleteCmd(nil), []string{"--name", "MyTemplate", "--yes", "-o", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), `"name"`) {
		t.Errorf("expected JSON output")
	}
}

func TestPermissionTemplatesDeleteMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, &testutil.MockPermissionTemplateClient{})
	err := fx.Run(permission_templates.NewDeleteCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestPermissionTemplatesDeleteError(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewDeleteCmd(nil), []string{"--name", "MyTemplate", "--yes"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPermissionTemplatesDeleteByID(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: id, Name: "MyTemplate"}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewDeleteCmd(nil), []string{"--id", "1", "--yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPermissionTemplatesDeleteQuiet(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: 1, Name: name}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewDeleteCmd(nil), []string{"--name", "MyTemplate", "--yes", "-q"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fx.Stdout.String() != "" {
		t.Errorf("expected no output in quiet mode, got:\n%s", fx.Stdout.String())
	}
}

func TestPermissionTemplatesDeleteNilWithoutError(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return nil, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewDeleteCmd(nil), []string{"--name", "Ghost", "--yes"})
	if err == nil {
		t.Fatal("expected error for nil template with nil error, got nil")
	}
}
