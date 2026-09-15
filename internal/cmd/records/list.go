// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/schema"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewListCmd returns a new "records list" command instance.
func NewListCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all records in a zone",
		Long:  `List all DNS records in a zone by name or ID.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			zoneName, _ := cmd.Flags().GetString("zone-name")
			zoneIDStr, _ := cmd.Flags().GetString("zone-id")

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

			records, err := client.Record.All(cmd.Context(), zoneID)
			if err != nil {
				return fmt.Errorf("failed to list records: %w", err)
			}

			sortBy, _ := cmd.Flags().GetString("sort")
			if sortBy != "" {
				sort.Slice(records, func(i, j int) bool {
					switch sortBy {
					case "type":
						return records[i].Type < records[j].Type
					case "ttl":
						return records[i].TTL < records[j].TTL
					default: // name
						return records[i].Name < records[j].Name
					}
				})
			}

			typeFilter, _ := cmd.Flags().GetString("type")
			if typeFilter != "" {
				filtered := records[:0]
				for _, r := range records {
					if strings.EqualFold(r.Type, typeFilter) {
						filtered = append(filtered, r)
					}
				}
				records = filtered
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.RecordListFromSDK(records))
			}

			t := base.NewTable(cmd)
			t.AddHeader("NAME", "TYPE", "CONTENT", "TTL")
			for _, r := range records {
				content := r.Content
				if outputFmt == output.FormatTable {
					content = output.Truncate(content, 50)
				}
				t.AddColoredRow(
					output.PlainCell(r.Name),
					output.Cell(r.Type, output.CyanCode()),
					output.PlainCell(content),
					output.PlainCell(strconv.Itoa(r.TTL)),
				)
			}
			t.Flush()
			return nil
		},
	}

	cmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("zone-id", "", "Zone ID")
	cmd.Flags().String("type", "", "Filter by record type (e.g. A, AAAA, MX, TXT)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|full|json")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")
	cmd.Flags().String("sort", "", "Sort by field. One of: name|type|ttl")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
