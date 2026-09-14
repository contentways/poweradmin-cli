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

func TestPermissionTemplatesList(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		ListFn: func(ctx context.Context) ([]*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return []*poweradmin.PermissionTemplate{
				{ID: 1, Name: "Administrator", Descr: "Full rights", TemplateType: "user"},
				{ID: 2, Name: "Editor", Descr: "Edit only", TemplateType: "user"},
			}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewListCmd(nil), []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "Administrator") {
		t.Errorf("expected Administrator in output, got:\n%s", fx.Stdout.String())
	}
}

func TestPermissionTemplatesListJSON(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		ListFn: func(ctx context.Context) ([]*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return []*poweradmin.PermissionTemplate{
				{ID: 1, Name: "Administrator", Descr: "Full rights", TemplateType: "user"},
			}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewListCmd(nil), []string{"-o", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), `"name": "Administrator"`) {
		t.Errorf("expected JSON output with name")
	}
}

func TestPermissionTemplatesListError(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		ListFn: func(ctx context.Context) ([]*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewListCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
