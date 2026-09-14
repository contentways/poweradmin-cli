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

func TestPermissionTemplatesCreate(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		CreateFn: func(ctx context.Context, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: 1, Name: opts.Name, TemplateType: opts.TemplateType}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewCreateCmd(nil), []string{"--name", "MyTemplate"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "MyTemplate") {
		t.Errorf("expected MyTemplate in output, got:\n%s", fx.Stdout.String())
	}
}

func TestPermissionTemplatesCreateJSON(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		CreateFn: func(ctx context.Context, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: 1, Name: opts.Name, TemplateType: "user"}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewCreateCmd(nil), []string{"--name", "MyTemplate", "-o", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "MyTemplate") {
		t.Errorf("expected JSON output")
	}
}

func TestPermissionTemplatesCreateMissingName(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, &testutil.MockPermissionTemplateClient{})
	err := fx.Run(permission_templates.NewCreateCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when --name missing")
	}
}

func TestPermissionTemplatesCreateError(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		CreateFn: func(ctx context.Context, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewCreateCmd(nil), []string{"--name", "MyTemplate"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
