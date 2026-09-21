// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

import (
	"errors"
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
	"github.com/spf13/cobra"
)

// NewDeleteCmd returns a new "users delete" command instance.
func NewDeleteCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a user",
		Long:  `Delete a Poweradmin user by username or numeric ID.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			interactive, _ := cmd.Flags().GetBool("interactive")

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			if interactive {
				base.PrintPreviewNotice(cmd)
				return runInteractiveDelete(cmd, client)
			}

			name, _ := cmd.Flags().GetString("name")
			idStr, _ := cmd.Flags().GetString("id")

			if name == "" && idStr == "" {
				return fmt.Errorf("either --name or --id is required")
			}

			user, err := base.ResolveUser(cmd, client)
			if err != nil {
				return err
			}

			if !base.Confirm(cmd, fmt.Sprintf("Delete user %s (id %d)? [y/N] ", user.Username, user.ID)) {
				return nil
			}

			_, err = client.User.Delete(cmd.Context(), user.ID)
			if err != nil {
				return fmt.Errorf("failed to delete user: %w", err)
			}

			if base.IsQuiet(cmd) {
				return nil
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, map[string]any{
					"id":       user.ID,
					"username": user.Username,
				})
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted user %s (id %d)\n", user.Username, user.ID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Username to identify the user")
	cmd.Flags().String("id", "", "User ID to identify the user")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output after deletion")
	cmd.Flags().BoolP("interactive", "i", false, "Interactively select users to delete (feature preview)")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.UserNameCompletion(s))
	return cmd
}

// runInteractiveDelete lists all users, lets the user pick zero or more via
// a multi-select prompt, then hands off to deleteSelectedUsers.
func runInteractiveDelete(cmd *cobra.Command, client *poweradmin.Client) error {
	usrs, err := client.User.All(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}
	if len(usrs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no users found")
		return nil
	}

	labels := make([]string, len(usrs))
	byLabel := make(map[string]*poweradmin.User, len(usrs))
	for i, u := range usrs {
		label := fmt.Sprintf("%s (id %d)", u.Username, u.ID)
		labels[i] = label
		byLabel[label] = u
	}

	selectedLabels, err := base.PromptMultiSelect("Select users to delete", labels)
	if err != nil {
		return err
	}

	selected := make([]*poweradmin.User, 0, len(selectedLabels))
	for _, label := range selectedLabels {
		selected = append(selected, byLabel[label])
	}

	return deleteSelectedUsers(cmd, client, selected)
}

// deleteSelectedUsers confirms and deletes the given users. Failures on
// individual users do not stop the remaining deletions; all errors are
// collected and returned together at the end. Split out from
// runInteractiveDelete so the deletion/confirmation/error-collection logic
// can be exercised directly in tests without going through the multi-select
// prompt.
func deleteSelectedUsers(cmd *cobra.Command, client *poweradmin.Client, selected []*poweradmin.User) error {
	if len(selected) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no users selected, nothing to do")
		return nil
	}

	summary := fmt.Sprintf("The following %d user(s) will be deleted:\n", len(selected))
	for _, u := range selected {
		summary += fmt.Sprintf("  - %s (id %d)\n", u.Username, u.ID)
	}
	summary += "\nProceed? [y/N] "
	if !base.Confirm(cmd, summary) {
		return nil
	}

	var errs []error
	for _, u := range selected {
		if _, err := client.User.Delete(cmd.Context(), u.ID); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "failed to delete user %s (id %d): %s\n", u.Username, u.ID, err)
			errs = append(errs, fmt.Errorf("user %s: %w", u.Username, err))
			continue
		}
		fmt.Fprintf(cmd.OutOrStdout(), "deleted user %s (id %d)\n", u.Username, u.ID)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to delete %d of %d user(s): %w", len(errs), len(selected), errors.Join(errs...))
	}
	return nil
}
