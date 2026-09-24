// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"errors"
	"fmt"
	"strings"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
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
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			if interactive {
				base.PrintPreviewNotice(cmd)
				return runInteractiveDelete(cmd, client, dryRun)
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

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would delete group %s (id %d). No changes made.\n", group.Name, group.ID)
				return nil
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
	cmd.Flags().Bool("dry-run", false, "Show what would be deleted without making changes")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.GroupNameCompletion(s))
	return cmd
}

// runInteractiveDelete lists all groups, lets the user pick zero or more via
// a multi-select prompt, then hands off to deleteSelectedGroups.
func runInteractiveDelete(cmd *cobra.Command, client *poweradmin.Client, dryRun bool) error {
	grps, err := client.Group.All(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to list groups: %w", err)
	}
	if len(grps) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no groups found")
		return nil
	}

	labels := make([]string, len(grps))
	byLabel := make(map[string]*poweradmin.Group, len(grps))
	for i, g := range grps {
		label := fmt.Sprintf("%s (id %d)", g.Name, g.ID)
		labels[i] = label
		byLabel[label] = g
	}

	selectedLabels, err := base.PromptMultiSelect("Select groups to delete", labels)
	if err != nil {
		return err
	}

	selected := make([]*poweradmin.Group, 0, len(selectedLabels))
	for _, label := range selectedLabels {
		selected = append(selected, byLabel[label])
	}

	return deleteSelectedGroups(cmd, client, selected, dryRun)
}

// deleteSelectedGroups confirms and deletes the given groups. In dry-run
// mode it prints what would be deleted and returns without confirming or
// making any API calls. Failures on individual groups do not stop the
// remaining deletions; all errors are collected and returned together at
// the end. Split out from runInteractiveDelete so the deletion/
// confirmation/error-collection logic can be exercised directly in tests
// without going through the multi-select prompt.
func deleteSelectedGroups(cmd *cobra.Command, client *poweradmin.Client, selected []*poweradmin.Group, dryRun bool) error {
	if len(selected) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no groups selected, nothing to do")
		return nil
	}

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "Would delete %d group(s):\n", len(selected))
		for _, g := range selected {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s (id %d)\n", g.Name, g.ID)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "No changes made.")
		return nil
	}

	var summary strings.Builder
	fmt.Fprintf(&summary, "The following %d group(s) will be deleted:\n", len(selected))
	for _, g := range selected {
		fmt.Fprintf(&summary, "  - %s (id %d)\n", g.Name, g.ID)
	}
	summary.WriteString("\nProceed? [y/N] ")
	if !base.Confirm(cmd, summary.String()) {
		return nil
	}

	var errs []error
	for _, g := range selected {
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
