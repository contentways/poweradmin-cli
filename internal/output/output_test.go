// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package output_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/output"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected output.Format
	}{
		{"table", output.FormatTable},
		{"full", output.FormatFull},
		{"json", output.FormatJSON},
		{"yml", output.FormatYAML},
		{"unknown", output.FormatTable},
		{"", output.FormatTable},
	}

	for _, tt := range tests {
		got := output.ParseFormat(tt.input)
		if got != tt.expected {
			t.Errorf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFormatIsStructured(t *testing.T) {
	tests := []struct {
		format   output.Format
		expected bool
	}{
		{output.FormatTable, false},
		{output.FormatFull, false},
		{output.FormatJSON, true},
		{output.FormatYAML, true},
	}

	for _, tt := range tests {
		got := tt.format.IsStructured()
		if got != tt.expected {
			t.Errorf("%q.IsStructured() = %v, want %v", tt.format, got, tt.expected)
		}
	}
}
