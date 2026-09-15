// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

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

// NewListCmd returns a new "zones list" command instance.
func NewListCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all DNS zones",
		Long:  `List all DNS zones in Poweradmin.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			zones, err := client.Zone.All(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to list zones: %w", err)
			}

			sortBy, _ := cmd.Flags().GetString("sort")
			if sortBy != "" {
				sort.Slice(zones, func(i, j int) bool {
					switch sortBy {
					case "name":
						return zones[i].Name < zones[j].Name
					case "type":
						return string(zones[i].Type) < string(zones[j].Type)
					default: // id
						return zones[i].ID < zones[j].ID
					}
				})
			}

			typeFilter, _ := cmd.Flags().GetString("type")
			nameFilter, _ := cmd.Flags().GetString("name-filter")

			if typeFilter != "" || nameFilter != "" {
				filtered := zones[:0]
				for _, z := range zones {
					if typeFilter != "" && !strings.EqualFold(string(z.Type), typeFilter) {
						continue
					}
					if nameFilter != "" && !strings.Contains(z.Name, nameFilter) {
						continue
					}
					filtered = append(filtered, z)
				}
				zones = filtered
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.ZoneListFromSDK(zones))
			}

			t := base.NewTable(cmd)
			t.AddHeader("ID", "NAME", "TYPE")
			for _, z := range zones {
				t.AddColoredRow(
					output.PlainCell(strconv.Itoa(z.ID)),
					output.PlainCell(z.Name),
					output.Cell(string(z.Type), output.CyanCode()),
				)
			}
			t.Flush()
			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|full")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")
	cmd.Flags().String("type", "", "Filter by zone type. One of: NATIVE|MASTER|SLAVE")
	cmd.Flags().String("name-filter", "", "Filter by zone name (substring match)")
	cmd.Flags().String("sort", "", "Sort by field. One of: id|name|type")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
