// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package base_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/spf13/cobra"
)

func TestDeprecatedAlias(t *testing.T) {
	cmd := base.DeprecatedAlias(&cobra.Command{Use: "get"}, "metadata-get", "zones metadata get")

	if cmd.Use != "metadata-get" {
		t.Errorf("Use = %q, want metadata-get", cmd.Use)
	}
	if cmd.Deprecated != `use "zones metadata get" instead` {
		t.Errorf("Deprecated = %q", cmd.Deprecated)
	}
	if cmd.IsAvailableCommand() {
		t.Error("deprecated alias must not be listed in help")
	}
}
