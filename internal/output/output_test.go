// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package output_test

import (
	"bytes"
	"strings"
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

func TestTableNoHeader(t *testing.T) {
	var buf bytes.Buffer
	tbl := output.New(&buf)
	tbl.SetNoHeader(true)
	tbl.AddHeader("ID", "NAME")
	tbl.AddRow("1", "example.com")
	tbl.Flush()

	out := buf.String()
	if strings.Contains(out, "ID") {
		t.Errorf("expected header to be suppressed, got:\n%s", out)
	}
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected data row present, got:\n%s", out)
	}
}

func TestTableColoredRowNoTTYIgnoresColor(t *testing.T) {
	var buf bytes.Buffer
	tbl := output.New(&buf)
	tbl.AddHeader("ID", "TYPE")
	tbl.AddColoredRow(output.PlainCell("1"), output.Cell("NATIVE", output.CyanCode()))
	tbl.Flush()

	out := buf.String()
	if !strings.Contains(out, "NATIVE") {
		t.Errorf("expected output to contain NATIVE, got:\n%s", out)
	}
	// Since go test doesn't run in a TTY, no ANSI escape codes should appear.
	if strings.Contains(out, "\033[") {
		t.Errorf("expected no ANSI escape codes outside a TTY, got:\n%s", out)
	}
}

func TestTableEmptyNoRows(t *testing.T) {
	var buf bytes.Buffer
	tbl := output.New(&buf)
	tbl.AddHeader("ID", "NAME")
	tbl.Flush()

	out := buf.String()
	if !strings.Contains(out, "ID") {
		t.Errorf("expected header present even with no rows, got:\n%s", out)
	}
}

func TestTableMultipleColoredRows(t *testing.T) {
	var buf bytes.Buffer
	tbl := output.New(&buf)
	tbl.AddHeader("ID", "NAME", "TYPE")
	tbl.AddColoredRow(output.PlainCell("1"), output.PlainCell("example.com"), output.Cell("NATIVE", output.CyanCode()))
	tbl.AddColoredRow(output.PlainCell("2"), output.PlainCell("example.org"), output.Cell("MASTER", output.CyanCode()))
	tbl.Flush()

	out := buf.String()
	if !strings.Contains(out, "example.com") || !strings.Contains(out, "example.org") {
		t.Errorf("expected both rows present, got:\n%s", out)
	}
}
