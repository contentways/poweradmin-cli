// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"fmt"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/schema"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewGetCmd returns a new "records get" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// The zone can be identified by name (--zone-name) or numeric ID (--zone-id).
// The record is identified by its opaque string ID (--id).
// Note: the Poweradmin API does not reliably support GET by record ID,
// so all records in the zone are fetched and searched client-side.
// Output can be a human-readable key-value summary (default) or JSON.
func NewGetCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a DNS record by ID",
		Long:  `Get a single DNS record by its ID from a zone.`,
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

			// The Poweradmin API does not reliably support GET by record ID.
			// Fetch all records and search client-side.
			records, err := client.Record.All(cmd.Context(), zoneID)
			if err != nil {
				return fmt.Errorf("failed to get records: %w", err)
			}

			var record *poweradmin.Record
			for _, r := range records {
				if r.ID == recordID {
					record = r
					break
				}
			}
			if record == nil {
				return fmt.Errorf("record not found: %s", recordID)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.RecordFromSDK(record))
			}

			// Default output — human-readable key-value summary.
			fmt.Fprintf(cmd.OutOrStdout(), "ID:       %s\n", record.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Name:     %s\n", record.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Type:     %s\n", output.Cyan(record.Type))
			fmt.Fprintf(cmd.OutOrStdout(), "Content:  %s\n", record.Content)
			fmt.Fprintf(cmd.OutOrStdout(), "TTL:      %d\n", record.TTL)
			if record.Priority > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Priority: %d\n", record.Priority)
			}
			return nil
		},
	}

	cmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("zone-id", "", "Zone ID")
	cmd.Flags().String("id", "", "Record ID (required)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
