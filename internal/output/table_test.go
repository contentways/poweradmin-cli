// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/output"
)

func TestTableOutput(t *testing.T) {
	var buf bytes.Buffer
	tbl := output.New(&buf)
	tbl.AddHeader("ID", "NAME", "TYPE")
	tbl.AddRow("1", "example.com", "NATIVE")
	tbl.AddRow("2", "example.org", "NATIVE")
	tbl.Flush()

	out := buf.String()
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected output to contain example.com, got:\n%s", out)
	}
	if !strings.Contains(out, "ID") {
		t.Errorf("expected output to contain ID header, got:\n%s", out)
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 8, "hello..."},
		{"hi", 2, "hi"},
		{"hello", 5, "hello"},
	}

	for _, tt := range tests {
		got := output.Truncate(tt.input, tt.max)
		if got != tt.expected {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.expected)
		}
	}
}
