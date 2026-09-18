// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
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

			zoneName, _ := cmd.Flags().GetString("zone-name")
			zoneIDStr, _ := cmd.Flags().GetString("zone-id")
			recordID, _ := cmd.Flags().GetString("id")

			if zoneName == "" && zoneIDStr == "" {
				return fmt.Errorf("either --zone-name or --zone-id is required")
			}
			if recordID == "" {
				return fmt.Errorf("--id is required")
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
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
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
