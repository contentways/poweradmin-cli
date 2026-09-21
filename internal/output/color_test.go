// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package output_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/output"
)

// go test does not run with stdout attached to a real terminal, so
// output.IsTTY() is false for the duration of these tests. That means
// every color-wrapping function below is exercised on its non-TTY path,
// where it must return the input unchanged.

func TestIsTTYFalseInTestEnvironment(t *testing.T) {
	if output.IsTTY() {
		t.Skip("stdout is unexpectedly a TTY in this test environment; the non-TTY assertions below don't apply")
	}
}

func TestBoldNoTTYReturnsUnchanged(t *testing.T) {
	if got := output.Bold("hello"); got != "hello" {
		t.Errorf("Bold(%q) = %q, want unchanged", "hello", got)
	}
}

func TestCyanNoTTYReturnsUnchanged(t *testing.T) {
	if got := output.Cyan("A"); got != "A" {
		t.Errorf("Cyan(%q) = %q, want unchanged", "A", got)
	}
}

func TestGreenNoTTYReturnsUnchanged(t *testing.T) {
	if got := output.Green("active"); got != "active" {
		t.Errorf("Green(%q) = %q, want unchanged", "active", got)
	}
}

func TestRedNoTTYReturnsUnchanged(t *testing.T) {
	if got := output.Red("inactive"); got != "inactive" {
		t.Errorf("Red(%q) = %q, want unchanged", "inactive", got)
	}
}

func TestYellowNoTTYReturnsUnchanged(t *testing.T) {
	if got := output.Yellow("warning"); got != "warning" {
		t.Errorf("Yellow(%q) = %q, want unchanged", "warning", got)
	}
}

func TestGrayNoTTYReturnsUnchanged(t *testing.T) {
	if got := output.Gray("info"); got != "info" {
		t.Errorf("Gray(%q) = %q, want unchanged", "info", got)
	}
}

func TestColorCodesAreDistinctAndNonEmpty(t *testing.T) {
	codes := map[string]string{
		"cyan":  output.CyanCode(),
		"green": output.GreenCode(),
		"red":   output.RedCode(),
	}
	seen := make(map[string]string)
	for name, code := range codes {
		if code == "" {
			t.Errorf("%s code is empty", name)
		}
		if other, ok := seen[code]; ok {
			t.Errorf("%s and %s share the same code %q", name, other, code)
		}
		seen[code] = name
	}
}
