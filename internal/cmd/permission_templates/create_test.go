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

// TestPermissionTemplatesCreateInteractiveSkipsPromptsWhenAllFlagsSet
// verifies that when --interactive is combined with all required values
// already supplied via flags (--description, --type, --permissions
// explicitly set so those huh prompts don't fire), no huh prompt runs at
// all — only the final confirmation, answered via a piped stdin.
func TestPermissionTemplatesCreateInteractiveSkipsPromptsWhenAllFlagsSet(t *testing.T) {
	var createdName, createdDescr, createdType string
	var createdPerms []int

	mock := &testutil.MockPermissionTemplateClient{
		CreateFn: func(ctx context.Context, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			createdName = opts.Name
			createdDescr = opts.Descr
			createdType = opts.TemplateType
			createdPerms = opts.Permissions
			return &poweradmin.PermissionTemplate{ID: 1, Name: opts.Name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	testutil.WithStdin(t, "y\n")

	err := fx.Run(permission_templates.NewCreateCmd(nil), []string{
		"--name", "MyTemplate",
		"--description", "Can edit zone records",
		"--type", "group",
		"--permissions", "1,2,3",
		"--interactive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if createdName != "MyTemplate" {
		t.Errorf("expected name MyTemplate, got %q", createdName)
	}
	if createdDescr != "Can edit zone records" {
		t.Errorf("expected description 'Can edit zone records', got %q", createdDescr)
	}
	if createdType != "group" {
		t.Errorf("expected type group, got %q", createdType)
	}
	if len(createdPerms) != 3 {
		t.Errorf("expected 3 permissions, got %v", createdPerms)
	}
}

// TestPermissionTemplatesCreateInteractiveDeclineAbortsWithoutCreating
// verifies that declining the final confirmation prevents the API call
// entirely.
func TestPermissionTemplatesCreateInteractiveDeclineAbortsWithoutCreating(t *testing.T) {
	created := false
	mock := &testutil.MockPermissionTemplateClient{
		CreateFn: func(ctx context.Context, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			created = true
			return &poweradmin.PermissionTemplate{ID: 1, Name: opts.Name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	testutil.WithStdin(t, "n\n")

	err := fx.Run(permission_templates.NewCreateCmd(nil), []string{
		"--name", "MyTemplate",
		"--description", "Can edit zone records",
		"--type", "group",
		"--permissions", "1,2,3",
		"--interactive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created {
		t.Error("expected template NOT to be created after declining confirmation")
	}
}

func TestPermissionTemplatesCreateQuiet(t *testing.T) {
	mock := &testutil.MockPermissionTemplateClient{
		CreateFn: func(ctx context.Context, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return &poweradmin.PermissionTemplate{ID: 1, Name: opts.Name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)
	err := fx.Run(permission_templates.NewCreateCmd(nil), []string{"--name", "MyTemplate", "-q"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(fx.Stdout.String())
	if out != "1" {
		t.Errorf("expected quiet output to be just the id, got: %q", out)
	}
}

// TestPermissionTemplatesCreateInteractiveFullFlow drives every prompt in
// "permission-templates create --interactive" via huh's accessible mode.
func TestPermissionTemplatesCreateInteractiveFullFlow(t *testing.T) {
	testutil.WithAccessiblePrompts(t)

	var createdName, createdDescr, createdType string
	var createdPerms []int
	mock := &testutil.MockPermissionTemplateClient{
		CreateFn: func(ctx context.Context, opts poweradmin.PermissionTemplateOpts) (*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			createdName = opts.Name
			createdDescr = opts.Descr
			createdType = opts.TemplateType
			createdPerms = opts.Permissions
			return &poweradmin.PermissionTemplate{ID: 1, Name: opts.Name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mock)

	// name, description, type (1 = user, skips nothing since both options
	// are shown), permissions, then the final confirmation.
	testutil.WithDelayedStdin(t, "MyTemplate\n", "Can edit zones\n", "1\n", "1,2,3\n", "y\n")

	err := fx.Run(permission_templates.NewCreateCmd(nil), []string{"--interactive"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if createdName != "MyTemplate" {
		t.Errorf("expected name MyTemplate, got %q", createdName)
	}
	if createdDescr != "Can edit zones" {
		t.Errorf("expected description 'Can edit zones', got %q", createdDescr)
	}
	if createdType != "user" {
		t.Errorf("expected type user, got %q", createdType)
	}
	if len(createdPerms) != 3 {
		t.Errorf("expected 3 permissions, got %v", createdPerms)
	}
}
