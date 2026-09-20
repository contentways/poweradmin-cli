// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
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

			interactive, _ := cmd.Flags().GetBool("interactive")

			zoneName, _ := cmd.Flags().GetString("zone-name")
			zoneIDStr, _ := cmd.Flags().GetString("zone-id")
			name, _ := cmd.Flags().GetString("name")
			recordType, _ := cmd.Flags().GetString("type")
			content, _ := cmd.Flags().GetString("content")
			ttl, _ := cmd.Flags().GetInt("ttl")
			priority, _ := cmd.Flags().GetInt("priority")

			if interactive {
				var err error

				if zoneName == "" && zoneIDStr == "" {
					zoneName, err = base.PromptString("Zone name", "e.g. example.com", true)
					if err != nil {
						return err
					}
					if err := cmd.Flags().Set("zone-name", zoneName); err != nil {
						return fmt.Errorf("failed to set zone-name: %w", err)
					}
				}

				if !cmd.Flags().Changed("name") {
					name, err = base.PromptString("Record name", "e.g. www.example.com", true)
					if err != nil {
						return err
					}
				}

				if !cmd.Flags().Changed("type") {
					recordType, err = base.PromptSelect(
						"Record type",
						[]string{"A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV"},
						recordType,
					)
					if err != nil {
						return err
					}
				}

				if !cmd.Flags().Changed("content") {
					content, err = base.PromptString("Content", "e.g. 1.2.3.4 for A records", true)
					if err != nil {
						return err
					}
				}

				if !cmd.Flags().Changed("ttl") {
					ttl, err = base.PromptInt("TTL", "in seconds", ttl)
					if err != nil {
						return err
					}
				}

				if (recordType == "MX" || recordType == "SRV") && !cmd.Flags().Changed("priority") {
					priority, err = base.PromptInt("Priority", "used for MX/SRV records", priority)
					if err != nil {
						return err
					}
				}

				zoneRef := zoneName
				if zoneRef == "" {
					zoneRef = "id " + zoneIDStr
				}
				summary := fmt.Sprintf(
					"Create record:\n  Zone:      %s\n  Name:      %s\n  Type:      %s\n  Content:   %s\n  TTL:       %d\n  Priority:  %d\n\nProceed? [y/N] ",
					zoneRef, name, recordType, content, ttl, priority,
				)
				if !base.Confirm(cmd, summary) {
					return nil
				}
			}

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

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, map[string]any{
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
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("quiet", "q", false, "Only print the ID of the created record")
	cmd.Flags().BoolP("interactive", "i", false, "Prompt interactively for missing values")
	return cmd
}
