// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"errors"
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
	"github.com/spf13/cobra"
)

// NewDeleteCmd returns a new "records delete" command instance.
func NewDeleteCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a DNS record",
		Long:  `Delete a DNS record by ID from a zone.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			interactive, _ := cmd.Flags().GetBool("interactive")

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			if interactive {
				base.PrintPreviewNotice(cmd)

				zoneName, _ := cmd.Flags().GetString("zone-name")
				zoneIDStr, _ := cmd.Flags().GetString("zone-id")

				if zoneName == "" && zoneIDStr == "" {
					zoneName, err = selectZoneName(cmd, client)
					if err != nil {
						return err
					}
					if err := cmd.Flags().Set("zone-name", zoneName); err != nil {
						return fmt.Errorf("failed to set zone-name: %w", err)
					}
				}

				zoneID, err := base.ResolveZoneID(cmd, client)
				if err != nil {
					return err
				}

				return runInteractiveDelete(cmd, client, zoneID)
			}

			zoneName, _ := cmd.Flags().GetString("zone-name")
			zoneIDStr, _ := cmd.Flags().GetString("zone-id")
			recordID, _ := cmd.Flags().GetString("id")

			if zoneName == "" && zoneIDStr == "" {
				return fmt.Errorf("either --zone-name or --zone-id is required")
			}
			if recordID == "" {
				return fmt.Errorf("--id is required")
			}

			zoneID, err := base.ResolveZoneID(cmd, client)
			if err != nil {
				return err
			}

			if !base.Confirm(cmd, fmt.Sprintf("Delete record (id %s)? [y/N] ", recordID)) {
				return nil
			}

			_, err = client.Record.Delete(cmd.Context(), zoneID, recordID)
			if err != nil {
				return fmt.Errorf("failed to delete record: %w", err)
			}

			if base.IsQuiet(cmd) {
				return nil
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, map[string]any{
					"id":      recordID,
					"zone_id": zoneID,
				})
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted record (id %s) from zone (id %d)\n", recordID, zoneID)
			return nil
		},
	}

	cmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("zone-id", "", "Zone ID")
	cmd.Flags().String("id", "", "Record ID (opaque string returned by the API)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output after deletion")
	cmd.Flags().BoolP("interactive", "i", false, "Interactively select records to delete (feature preview)")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}

// selectZoneName lists all zones and lets the user pick one interactively.
func selectZoneName(cmd *cobra.Command, client *poweradmin.Client) (string, error) {
	zones, err := client.Zone.All(cmd.Context())
	if err != nil {
		return "", fmt.Errorf("failed to list zones: %w", err)
	}
	if len(zones) == 0 {
		return "", fmt.Errorf("no zones found")
	}

	names := make([]string, len(zones))
	for i, z := range zones {
		names[i] = z.Name
	}

	return base.PromptSelect("Zone", names, "")
}

// runInteractiveDelete lists all records in the given zone, lets the user
// pick zero or more via a multi-select prompt, then hands off to
// deleteSelectedRecords.
func runInteractiveDelete(cmd *cobra.Command, client *poweradmin.Client, zoneID int) error {
	records, _, err := client.Record.List(cmd.Context(), zoneID, poweradmin.RecordListOpts{})
	if err != nil {
		return fmt.Errorf("failed to list records: %w", err)
	}
	if len(records) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no records found")
		return nil
	}

	labels := make([]string, len(records))
	byLabel := make(map[string]*poweradmin.Record, len(records))
	for i, r := range records {
		label := fmt.Sprintf("%s %s %s (id %s)", r.Name, r.Type, r.Content, r.ID)
		labels[i] = label
		byLabel[label] = r
	}

	selectedLabels, err := base.PromptMultiSelect("Select records to delete", labels)
	if err != nil {
		return err
	}

	selected := make([]*poweradmin.Record, 0, len(selectedLabels))
	for _, label := range selectedLabels {
		selected = append(selected, byLabel[label])
	}

	return deleteSelectedRecords(cmd, client, zoneID, selected)
}

// deleteSelectedRecords confirms and deletes the given records from the
// given zone. Failures on individual records do not stop the remaining
// deletions; all errors are collected and returned together at the end.
// Split out from runInteractiveDelete so the deletion/confirmation/
// error-collection logic can be exercised directly in tests without going
// through the multi-select prompt.
func deleteSelectedRecords(cmd *cobra.Command, client *poweradmin.Client, zoneID int, selected []*poweradmin.Record) error {
	if len(selected) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no records selected, nothing to do")
		return nil
	}

	summary := fmt.Sprintf("The following %d record(s) will be deleted:\n", len(selected))
	for _, r := range selected {
		summary += fmt.Sprintf("  - %s %s %s (id %s)\n", r.Name, r.Type, r.Content, r.ID)
	}
	summary += "\nProceed? [y/N] "
	if !base.Confirm(cmd, summary) {
		return nil
	}

	var errs []error
	for _, r := range selected {
		if _, err := client.Record.Delete(cmd.Context(), zoneID, r.ID); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "failed to delete record %s %s (id %s): %s\n", r.Name, r.Type, r.ID, err)
			errs = append(errs, fmt.Errorf("record %s: %w", r.ID, err))
			continue
		}
		fmt.Fprintf(cmd.OutOrStdout(), "deleted record %s %s (id %s)\n", r.Name, r.Type, r.ID)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to delete %d of %d record(s): %w", len(errs), len(selected), errors.Join(errs...))
	}
	return nil
}
