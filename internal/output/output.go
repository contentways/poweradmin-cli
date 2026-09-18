// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package output provides shared output formatting utilities for the CLI.
// It supports four output formats: table (default), full, json and yaml.
package output

// Format represents the output format requested by the user via --output flag.
type Format string

const (
	// FormatTable renders a human-readable aligned table with truncated content.
	// This is the default format.
	FormatTable Format = "table"

	// FormatFull renders a human-readable aligned table without truncation.
	// Long content values (e.g. TXT records) are shown in full.
	FormatFull Format = "full"

	// FormatJSON renders the output as indented JSON.
	// Useful for scripting and piping into tools like jq.
	FormatJSON Format = "json"

	// FormatYAML renders the output as YAML.
	// Useful for scripting and piping into tools like yq.
	FormatYAML Format = "yaml"
)

// ParseFormat converts a raw string flag value into a Format constant.
// Any unrecognised value falls back to FormatTable.
func ParseFormat(s string) Format {
	switch s {
	case "full":
		return FormatFull
	case "json":
		return FormatJSON
	case "yaml", "yml":
		return FormatYAML
	default:
		return FormatTable
	}
}

// IsStructured reports whether the format is a machine-readable, structured
// format (JSON or YAML) as opposed to a human-readable table.
func (f Format) IsStructured() bool {
	return f == FormatJSON || f == FormatYAML
}
