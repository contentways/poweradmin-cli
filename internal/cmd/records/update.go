// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"encoding/json"
	"fmt"
	"strconv"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/schema"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewUpdateCmd returns a new "records update" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// The zone can be identified by name (--zone-name) or numeric ID (--zone-id).
// The record is identified by its opaque string ID (--id).
// Only flags that are explicitly set are sent to the API — unset flags are omitted.
// Output can be a human-readable confirmation (default) or JSON.
func NewUpdateCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a DNS record",
		Long:  `Update an existing DNS record in a zone.`,
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

			// Resolve zone ID — either parse the numeric flag directly,
			// or look up the zone by name to obtain its ID.
			var zoneID int
			if zoneIDStr != "" {
				zoneID, err = strconv.Atoi(zoneIDStr)
				if err != nil {
					return fmt.Errorf("invalid zone-id: %w", err)
				}
			} else {
				zone, _, err := client.Zone.GetByName(cmd.Context(), zoneName)
				if err != nil {
					return fmt.Errorf("failed to resolve zone: %w", err)
				}
				zoneID = zone.ID
			}

			// Build update opts — only include fields that were explicitly set.
			opts := poweradmin.RecordUpdateOpts{}
			if cmd.Flags().Changed("name") {
				opts.Name, _ = cmd.Flags().GetString("name")
			}
			if cmd.Flags().Changed("type") {
				opts.Type, _ = cmd.Flags().GetString("type")
			}
			if cmd.Flags().Changed("content") {
				opts.Content, _ = cmd.Flags().GetString("content")
			}
			if cmd.Flags().Changed("ttl") {
				ttl, _ := cmd.Flags().GetInt("ttl")
				opts.TTL = &ttl
			}
			if cmd.Flags().Changed("priority") {
				priority, _ := cmd.Flags().GetInt("priority")
				opts.Priority = &priority
			}
			if cmd.Flags().Changed("disabled") {
				disabled, _ := cmd.Flags().GetBool("disabled")
				opts.Disabled = &disabled
			}

			record, _, err := client.Record.Update(cmd.Context(), zoneID, recordID, opts)
			if err != nil {
				return fmt.Errorf("failed to update record: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — return the updated record.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(schema.RecordFromSDK(record), "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Default output — human-readable confirmation.
			fmt.Fprintf(cmd.OutOrStdout(), "updated record %s %s %s (id %s)\n", record.Name, record.Type, record.Content, record.ID)
			return nil
		},
	}

	cmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("zone-id", "", "Zone ID")
	cmd.Flags().String("id", "", "Record ID (required)")
	cmd.Flags().String("name", "", "New record name")
	cmd.Flags().String("type", "", "New record type. One of: A|AAAA|CNAME|MX|TXT|NS|SRV|...")
	cmd.Flags().String("content", "", "New record content")
	cmd.Flags().Int("ttl", 0, "New TTL in seconds")
	cmd.Flags().Int("priority", 0, "New priority (for MX records)")
	cmd.Flags().Bool("disabled", false, "Disable the record")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
