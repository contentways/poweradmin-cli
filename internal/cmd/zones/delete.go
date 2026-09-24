// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
	"github.com/spf13/cobra"
)

// NewDeleteCmd returns a new "zones delete" command instance.
func NewDeleteCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a DNS zone",
		Long:  `Delete a DNS zone from Poweradmin by name or ID.`,
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

			var zoneID int
			if idStr != "" {
				zoneID, err = strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
			} else {
				zone, _, err := client.Zone.GetByName(cmd.Context(), name)
				if err != nil {
					return fmt.Errorf("failed to resolve zone: %w", err)
				}
				if zone == nil {
					return fmt.Errorf("zone %q not found", name)
				}
				zoneID = zone.ID
				name = zone.Name
			}

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would delete zone %s (id %d). No changes made.\n", name, zoneID)
				return nil
			}

			if !base.Confirm(cmd, fmt.Sprintf("Delete zone %s (id %d)? [y/N] ", name, zoneID)) {
				return nil
			}

			_, err = client.Zone.Delete(cmd.Context(), zoneID)
			if err != nil {
				return fmt.Errorf("failed to delete zone: %w", err)
			}

			if base.IsQuiet(cmd) {
				return nil
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, map[string]any{
					"id":   zoneID,
					"name": name,
				})
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted zone %s (id %d)\n", name, zoneID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("id", "", "Zone ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output after deletion")
	cmd.Flags().BoolP("interactive", "i", false, "Interactively select zones to delete (feature preview)")
	cmd.Flags().Bool("dry-run", false, "Show what would be deleted without making changes")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}

// runInteractiveDelete lists all zones, lets the user pick zero or more via
// a multi-select prompt, then hands off to deleteSelectedZones.
func runInteractiveDelete(cmd *cobra.Command, client *poweradmin.Client, dryRun bool) error {
	zones, err := client.Zone.All(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to list zones: %w", err)
	}
	if len(zones) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no zones found")
		return nil
	}

	labels := make([]string, len(zones))
	byLabel := make(map[string]*poweradmin.Zone, len(zones))
	for i, z := range zones {
		label := fmt.Sprintf("%s (id %d)", z.Name, z.ID)
		labels[i] = label
		byLabel[label] = z
	}

	selectedLabels, err := base.PromptMultiSelect("Select zones to delete", labels)
	if err != nil {
		return err
	}

	selected := make([]*poweradmin.Zone, 0, len(selectedLabels))
	for _, label := range selectedLabels {
		selected = append(selected, byLabel[label])
	}

	return deleteSelectedZones(cmd, client, selected, dryRun)
}

// deleteSelectedZones confirms and deletes the given zones. In dry-run mode
// it prints what would be deleted and returns without confirming or making
// any API calls. Failures on individual zones do not stop the remaining
// deletions; all errors are collected and returned together at the end.
// Split out from runInteractiveDelete so the deletion/confirmation/
// error-collection logic can be exercised directly in tests without going
// through the multi-select prompt.
func deleteSelectedZones(cmd *cobra.Command, client *poweradmin.Client, selected []*poweradmin.Zone, dryRun bool) error {
	if len(selected) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no zones selected, nothing to do")
		return nil
	}

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "Would delete %d zone(s):\n", len(selected))
		for _, z := range selected {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s (id %d)\n", z.Name, z.ID)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "No changes made.")
		return nil
	}

	var summary strings.Builder
	fmt.Fprintf(&summary, "The following %d zone(s) will be deleted:\n", len(selected))
	for _, z := range selected {
		fmt.Fprintf(&summary, "  - %s (id %d)\n", z.Name, z.ID)
	}
	summary.WriteString("\nProceed? [y/N] ")
	if !base.Confirm(cmd, summary.String()) {
		return nil
	}

	var errs []error
	for _, z := range selected {
		if _, err := client.Zone.Delete(cmd.Context(), z.ID); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "failed to delete zone %s (id %d): %s\n", z.Name, z.ID, err)
			errs = append(errs, fmt.Errorf("zone %s: %w", z.Name, err))
			continue
		}
		fmt.Fprintf(cmd.OutOrStdout(), "deleted zone %s (id %d)\n", z.Name, z.ID)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to delete %d of %d zone(s): %w", len(errs), len(selected), errors.Join(errs...))
	}
	return nil
}
