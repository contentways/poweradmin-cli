// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil

import (
	"os"
	"testing"
)

// WithStdin temporarily replaces os.Stdin with a pipe pre-filled with input
// for the duration of the test, restoring the original os.Stdin via
// t.Cleanup. Used to feed answers to prompts that read directly from
// os.Stdin, such as base.Confirm.
func WithStdin(t *testing.T, input string) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}
	if _, err := w.WriteString(input); err != nil {
		t.Fatalf("failed to write to stdin pipe: %v", err)
	}
	w.Close()

	original := os.Stdin
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = original
		r.Close()
	})
}
