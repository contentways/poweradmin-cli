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
