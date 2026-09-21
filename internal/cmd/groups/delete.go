// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"errors"
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
	"github.com/spf13/cobra"
)

// NewDeleteCmd returns a new "groups delete" command instance.
func NewDeleteCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a group",
		Long:  `Delete a Poweradmin group by name or numeric ID.`,
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

			group, err := base.ResolveGroup(cmd, client)
			if err != nil {
				return err
			}

			if !base.Confirm(cmd, fmt.Sprintf("Delete group %s (id %d)? [y/N] ", group.Name, group.ID)) {
				return nil
			}

			_, err = client.Group.Delete(cmd.Context(), group.ID)
			if err != nil {
				return fmt.Errorf("failed to delete group: %w", err)
			}

			if base.IsQuiet(cmd) {
				return nil
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, map[string]any{
					"id":   group.ID,
					"name": group.Name,
				})
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted group %s (id %d)\n", group.Name, group.ID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Group name to identify the group")
	cmd.Flags().String("id", "", "Group ID to identify the group")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output after deletion")
	cmd.Flags().BoolP("interactive", "i", false, "Interactively select groups to delete (feature preview)")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.GroupNameCompletion(s))
	return cmd
}

// runInteractiveDelete lists all groups, lets the user pick zero or more via
// a multi-select prompt, confirms, and deletes each selected group. Failures
// on individual groups do not stop the remaining deletions; all errors are
// collected and returned together at the end.
func runInteractiveDelete(cmd *cobra.Command, client *poweradmin.Client) error {
	groups, err := client.Group.All(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to list groups: %w", err)
	}
	if len(groups) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no groups found")
		return nil
	}

	labels := make([]string, len(groups))
	byLabel := make(map[string]*poweradmin.Group, len(groups))
	for i, g := range groups {
		label := fmt.Sprintf("%s (id %d)", g.Name, g.ID)
		labels[i] = label
		byLabel[label] = g
	}

	selected, err := base.PromptMultiSelect("Select groups to delete", labels)
	if err != nil {
		return err
	}
	if len(selected) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no groups selected, nothing to do")
		return nil
	}

	summary := fmt.Sprintf("The following %d group(s) will be deleted:\n", len(selected))
	for _, label := range selected {
		summary += fmt.Sprintf("  - %s\n", label)
	}
	summary += "\nProceed? [y/N] "
	if !base.Confirm(cmd, summary) {
		return nil
	}

	var errs []error
	for _, label := range selected {
		g := byLabel[label]
		if _, err := client.Group.Delete(cmd.Context(), g.ID); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "failed to delete group %s (id %d): %s\n", g.Name, g.ID, err)
			errs = append(errs, fmt.Errorf("group %s: %w", g.Name, err))
			continue
		}
		fmt.Fprintf(cmd.OutOrStdout(), "deleted group %s (id %d)\n", g.Name, g.ID)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to delete %d of %d group(s): %w", len(errs), len(selected), errors.Join(errs...))
	}
	return nil
}
