// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package output_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/v2/internal/output"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected output.Format
	}{
		{"table", output.FormatTable},
		{"full", output.FormatFull},
		{"json", output.FormatJSON},
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
