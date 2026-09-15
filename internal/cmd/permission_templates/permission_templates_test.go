// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package permission_templates_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/v2/internal/cmd/permission_templates"
)

func TestNewPermissionTemplatesCommand(t *testing.T) {
	cmd := permission_templates.NewPermissionTemplatesCommand(nil)
	if cmd == nil {
		t.Fatal("expected command, got nil")
	}
	if got := len(cmd.Commands()); got != 5 {
		t.Errorf("got %d commands, want 5", got)
	}
}
