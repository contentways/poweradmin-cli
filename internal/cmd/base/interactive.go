// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

package base

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

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
				return fmt.Errorf("darf nicht leer sein")
			}
			return nil
		})
	}

	if err := huh.NewForm(huh.NewGroup(field)).Run(); err != nil {
		return "", fmt.Errorf("prompt abgebrochen: %w", err)
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

	if err := huh.NewForm(huh.NewGroup(field)).Run(); err != nil {
		return "", fmt.Errorf("prompt abgebrochen: %w", err)
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
				return fmt.Errorf("muss eine Zahl sein")
			}
			return nil
		})

	if err := huh.NewForm(huh.NewGroup(field)).Run(); err != nil {
		return 0, fmt.Errorf("prompt abgebrochen: %w", err)
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

	if err := huh.NewForm(huh.NewGroup(field)).Run(); err != nil {
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

	if err := huh.NewForm(huh.NewGroup(field)).Run(); err != nil {
		return nil, fmt.Errorf("prompt aborted: %w", err)
	}
	return selected, nil
}
