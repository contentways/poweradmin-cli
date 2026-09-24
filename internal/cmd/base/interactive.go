// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

package base

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// accessibleEnvVar, when set to a non-empty value, switches every prompt in
// this package to huh's accessible mode: plain line-based stdin/stdout
// prompts instead of the interactive TUI. This exists so tests can drive
// prompts via piped stdin without a real terminal; it is not intended as a
// user-facing accessibility toggle (unlike huh's own ACCESSIBLE convention,
// this one is deliberately unadvertised to end users).
const accessibleEnvVar = "POWERADMIN_TEST_ACCESSIBLE"

// isAccessible reports whether prompts should run in huh's accessible mode.
func isAccessible() bool {
	return os.Getenv(accessibleEnvVar) != ""
}

// runForm runs a single-field huh form, honoring accessible mode.
func runForm(field huh.Field) error {
	form := huh.NewForm(huh.NewGroup(field))
	if isAccessible() {
		// huh's accessible prompts wrap their input in a fresh bufio.Scanner
		// per prompt. Reading os.Stdin directly lets one prompt swallow the
		// lines meant for the following prompts whenever several lines are
		// already waiting in the pipe. lineReader hands out one line per Read
		// so every prompt consumes exactly its own line.
		form = form.WithAccessible(true).WithInput(lineReader{r: os.Stdin})
	}
	return form.Run()
}

// lineReader returns at most one line per Read call. It reads byte by byte
// and keeps no buffer, so nothing is lost between successive readers of the
// same underlying stream.
type lineReader struct {
	r io.Reader
}

func (l lineReader) Read(p []byte) (int, error) {
	var b [1]byte
	n := 0
	for n < len(p) {
		m, err := l.r.Read(b[:])
		if m == 1 {
			p[n] = b[0]
			n++
			if b[0] == '\n' {
				return n, nil
			}
		}
		if err != nil {
			if n > 0 {
				return n, nil
			}
			return 0, err
		}
	}
	return n, nil
}

// PromptString interactively asks the user for a single string value.
// If required is true, empty input is rejected.
func PromptString(title, description string, required bool) (string, error) {
	var value string
	field := huh.NewInput().
		Title(title).
		Description(description).
		Value(&value)

	if required {
		field = field.Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("must not be empty")
			}
			return nil
		})
	}

	if err := runForm(field); err != nil {
		return "", fmt.Errorf("prompt aborted: %w", err)
	}
	return value, nil
}

// PromptSelect interactively asks the user to choose one of options.
// selected is pre-filled with the current default.
func PromptSelect(title string, options []string, selected string) (string, error) {
	value := selected
	field := huh.NewSelect[string]().
		Title(title).
		Options(huh.NewOptions(options...)...).
		Value(&value)

	if err := runForm(field); err != nil {
		return "", fmt.Errorf("prompt aborted: %w", err)
	}
	return value, nil
}

// PromptInt interactively asks the user for an integer value.
// The input is pre-filled with the current default.
func PromptInt(title, description string, defaultValue int) (int, error) {
	value := fmt.Sprintf("%d", defaultValue)
	field := huh.NewInput().
		Title(title).
		Description(description).
		Value(&value).
		Validate(func(s string) error {
			var n int
			if _, err := fmt.Sscanf(s, "%d", &n); err != nil {
				return fmt.Errorf("must be a number")
			}
			return nil
		})

	if err := runForm(field); err != nil {
		return 0, fmt.Errorf("prompt aborted: %w", err)
	}

	var n int
	fmt.Sscanf(value, "%d", &n)
	return n, nil
}

// PromptStringSlice interactively asks for a comma-separated list of values.
// Returns nil if the input is empty.
func PromptStringSlice(title, description string) ([]string, error) {
	raw, err := PromptString(title, description, false)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			result = append(result, v)
		}
	}
	return result, nil
}

// PromptBool interactively asks the user a yes/no question.
// defaultValue pre-selects the corresponding option.
func PromptBool(title string, defaultValue bool) (bool, error) {
	value := defaultValue
	field := huh.NewConfirm().
		Title(title).
		Affirmative("Yes").
		Negative("No").
		Value(&value)

	if err := runForm(field); err != nil {
		return false, fmt.Errorf("prompt aborted: %w", err)
	}
	return value, nil
}

// PrintPreviewNotice writes a short notice to stderr indicating that
// --interactive is a feature preview and its behavior may change.
func PrintPreviewNotice(cmd *cobra.Command) {
	fmt.Fprintln(cmd.ErrOrStderr(), "⚠ --interactive is a feature preview and may change in a future release.")
}

// PromptMultiSelect interactively asks the user to choose zero or more of
// options using arrow keys and space to toggle selection.
func PromptMultiSelect(title string, options []string) ([]string, error) {
	var selected []string
	field := huh.NewMultiSelect[string]().
		Title(title).
		Options(huh.NewOptions(options...)...).
		Value(&selected)

	if err := runForm(field); err != nil {
		return nil, fmt.Errorf("prompt aborted: %w", err)
	}
	return selected, nil
}
