// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/users"
	"github.com/contentways/poweradmin-cli/internal/testutil"
)

func TestUsersSetPermissionTemplate(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		SetPermissionTemplateFn: func(ctx context.Context, id, permTemplID int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewSetPermissionTemplateCmd(nil), []string{
		"--name", "max",
		"--template-id", "5",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "max") {
		t.Errorf("expected output to contain max, got:\n%s", out)
	}
	if !strings.Contains(out, "5") {
		t.Errorf("expected output to contain template id 5, got:\n%s", out)
	}
}

func TestUsersSetPermissionTemplateJSON(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		SetPermissionTemplateFn: func(ctx context.Context, id, permTemplID int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewSetPermissionTemplateCmd(nil), []string{
		"--name", "max",
		"--template-id", "5",
		"--output", "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"template_id": 5`) {
		t.Errorf("expected JSON to contain template_id 5, got:\n%s", out)
	}
	if !strings.Contains(out, `"username": "max"`) {
		t.Errorf("expected JSON to contain max, got:\n%s", out)
	}
}

func TestUsersSetPermissionTemplateMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, &testutil.MockUserClient{}, nil, nil)

	err := fx.Run(users.NewSetPermissionTemplateCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestUsersSetPermissionTemplateError(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 42, Username: username}, nil, nil
		},
		SetPermissionTemplateFn: func(ctx context.Context, id, permTemplID int) (*poweradmin.Response, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	err := fx.Run(users.NewSetPermissionTemplateCmd(nil), []string{
		"--name", "max",
		"--template-id", "5",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
