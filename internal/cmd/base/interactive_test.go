// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package base_test

import (
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/spf13/cobra"
)

func TestPromptStringAccessibleMode(t *testing.T) {
	testutil.WithAccessiblePrompts(t)
	testutil.WithStdin(t, "example.com\n")

	value, err := base.PromptString("Zone name", "e.g. example.com", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "example.com" {
		t.Errorf("expected value example.com, got %q", value)
	}
}

func TestPromptStringEmptyOptional(t *testing.T) {
	testutil.WithAccessiblePrompts(t)
	testutil.WithStdin(t, "\n")

	value, err := base.PromptString("Description", "optional", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "" {
		t.Errorf("expected empty value, got %q", value)
	}
}

func TestPromptSelectAccessibleMode(t *testing.T) {
	testutil.WithAccessiblePrompts(t)
	testutil.WithStdin(t, "2\n")

	value, err := base.PromptSelect("Zone type", []string{"NATIVE", "MASTER", "SLAVE"}, "NATIVE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "MASTER" {
		t.Errorf("expected MASTER, got %q", value)
	}
}

func TestPromptIntAccessibleMode(t *testing.T) {
	testutil.WithAccessiblePrompts(t)
	testutil.WithStdin(t, "7200\n")

	value, err := base.PromptInt("TTL", "in seconds", 3600)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != 7200 {
		t.Errorf("expected 7200, got %d", value)
	}
}

func TestPromptIntInvalidThenValid(t *testing.T) {
	testutil.WithAccessiblePrompts(t)
	// huh's accessible mode re-prompts on validation failure, so feed an
	// invalid value first, then a valid one.
	testutil.WithStdin(t, "abc\n7200\n")

	value, err := base.PromptInt("TTL", "in seconds", 3600)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != 7200 {
		t.Errorf("expected 7200, got %d", value)
	}
}

func TestPromptStringSliceWithValues(t *testing.T) {
	testutil.WithAccessiblePrompts(t)
	testutil.WithStdin(t, "ns1.example.com,ns2.example.com\n")

	value, err := base.PromptStringSlice("Nameservers", "comma-separated")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(value) != 2 || value[0] != "ns1.example.com" || value[1] != "ns2.example.com" {
		t.Errorf("expected [ns1.example.com ns2.example.com], got %v", value)
	}
}

func TestPromptStringSliceEmpty(t *testing.T) {
	testutil.WithAccessiblePrompts(t)
	testutil.WithStdin(t, "\n")

	value, err := base.PromptStringSlice("Nameservers", "optional")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != nil {
		t.Errorf("expected nil, got %v", value)
	}
}

func TestPromptBoolYes(t *testing.T) {
	testutil.WithAccessiblePrompts(t)
	testutil.WithStdin(t, "y\n")

	value, err := base.PromptBool("Active?", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !value {
		t.Error("expected true")
	}
}

func TestPromptBoolDefaultOnEmptyInput(t *testing.T) {
	testutil.WithAccessiblePrompts(t)
	testutil.WithStdin(t, "\n")

	value, err := base.PromptBool("Active?", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !value {
		t.Errorf("expected default true to be kept on empty input, got %v", value)
	}
}

func TestPromptMultiSelectTogglesAndConfirms(t *testing.T) {
	testutil.WithAccessiblePrompts(t)
	// Toggle option 1, toggle option 3, then 0 confirms the selection.
	// huh's accessible-mode MultiSelect needs a brief pause between rapid
	// successive inputs to register each toggle correctly, hence
	// WithDelayedStdin instead of WithStdin here.
	testutil.WithDelayedStdin(t, "1\n", "3\n", "0\n")

	value, err := base.PromptMultiSelect("Select zones", []string{"a.com", "b.com", "c.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(value) != 2 {
		t.Fatalf("expected 2 selected items, got %d: %v", len(value), value)
	}
	got := strings.Join(value, ",")
	if !strings.Contains(got, "a.com") || !strings.Contains(got, "c.com") {
		t.Errorf("expected a.com and c.com selected, got %v", value)
	}
}

func TestPromptMultiSelectEmpty(t *testing.T) {
	testutil.WithAccessiblePrompts(t)
	// 0 immediately confirms with nothing toggled.
	testutil.WithStdin(t, "0\n")

	value, err := base.PromptMultiSelect("Select zones", []string{"a.com", "b.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(value) != 0 {
		t.Errorf("expected empty selection, got %v", value)
	}
}

func TestPrintPreviewNotice(t *testing.T) {
	var buf strings.Builder
	cmd := &cobra.Command{Use: "test"}
	cmd.SetErr(&buf)

	base.PrintPreviewNotice(cmd)

	if !strings.Contains(buf.String(), "feature preview") {
		t.Errorf("expected preview notice on stderr, got:\n%s", buf.String())
	}
}
