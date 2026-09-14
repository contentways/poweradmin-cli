// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"fmt"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewCreateCmd returns a new "records create" command instance.
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a DNS record",
		Long:  `Create a new DNS record in a zone.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			zoneName, _ := cmd.Flags().GetString("zone-name")
			zoneIDStr, _ := cmd.Flags().GetString("zone-id")
			name, _ := cmd.Flags().GetString("name")
			recordType, _ := cmd.Flags().GetString("type")
			content, _ := cmd.Flags().GetString("content")
			ttl, _ := cmd.Flags().GetInt("ttl")
			priority, _ := cmd.Flags().GetInt("priority")

			if zoneName == "" && zoneIDStr == "" {
				return fmt.Errorf("either --zone-name or --zone-id is required")
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			zoneID, err := base.ResolveZoneID(cmd, client)
			if err != nil {
				return err
			}

			id, _, err := client.Record.Create(cmd.Context(), zoneID, poweradmin.RecordCreateOpts{
				Name:     name,
				Type:     recordType,
				Content:  content,
				TTL:      ttl,
				Priority: priority,
			})
			if err != nil {
				return fmt.Errorf("failed to create record: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, map[string]any{
					"id":      id,
					"name":    name,
					"type":    recordType,
					"content": content,
					"ttl":     ttl,
					"zone_id": zoneID,
				})
			}

			if base.IsQuiet(cmd) {
				fmt.Fprintln(cmd.OutOrStdout(), id)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "created record %s %s %s (id %s)\n", name, recordType, content, id)
			return nil
		},
	}

	cmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("zone-id", "", "Zone ID")
	cmd.Flags().String("name", "", "Record name (e.g. www.example.com)")
	cmd.Flags().String("type", "", "Record type. One of: A|AAAA|CNAME|MX|TXT|NS|SRV|...")
	cmd.Flags().String("content", "", "Record content (e.g. 1.2.3.4 for A records)")
	cmd.Flags().Int("ttl", 3600, "Time to live in seconds (default: 3600)")
	cmd.Flags().Int("priority", 0, "Record priority, used for MX records (default: 0)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().BoolP("quiet", "q", false, "Only print the ID of the created record")
	return cmd
}
