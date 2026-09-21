// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil

import (
	"os"
	"testing"
	"time"
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

// WithAccessiblePrompts sets POWERADMIN_TEST_ACCESSIBLE for the duration of
// the test, switching every huh-based prompt in internal/cmd/base to plain
// line-based stdin/stdout mode instead of the interactive TUI. Combine with
// WithStdin to drive prompts programmatically in tests.
func WithAccessiblePrompts(t *testing.T) {
	t.Helper()
	t.Setenv("POWERADMIN_TEST_ACCESSIBLE", "1")
}

// WithDelayedStdin behaves like WithStdin, but writes each line separately
// with a short delay in between. huh's accessible-mode MultiSelect prompt
// appears to need a brief pause between rapid successive inputs to register
// each toggle correctly — piping all input at once can cause intermediate
// selections to be dropped. Use this instead of WithStdin when driving a
// PromptMultiSelect in tests.
func WithDelayedStdin(t *testing.T, lines ...string) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}

	original := os.Stdin
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = original
		r.Close()
	})

	go func() {
		defer w.Close()
		for _, line := range lines {
			if _, err := w.WriteString(line); err != nil {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
}
