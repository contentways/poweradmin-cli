// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package version_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	cmdversion "github.com/contentways/poweradmin-cli/v2/internal/cmd/version"
	"github.com/contentways/poweradmin-cli/v2/internal/state"
)

func TestVersionCmd(t *testing.T) {
	cmd := cmdversion.NewVersionCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	ctx := state.New("https://test.example.com", "test-key").WithContext(context.Background())
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "poweradmin version") {
		t.Errorf("expected output to contain 'poweradmin version', got:\n%s", out)
	}
}
