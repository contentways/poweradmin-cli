// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewZonesCmd returns a new "groups zones" command instance.
func NewZonesCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zones",
		Short: "List zones of a group",
		Long:  `List all zones associated with a Poweradmin group.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			name, _ := cmd.Flags().GetString("name")
			idStr, _ := cmd.Flags().GetString("id")

			if name == "" && idStr == "" {
				return fmt.Errorf("either --name or --id is required")
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			var groupID int
			if idStr != "" {
				groupID, err = strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
			} else {
				group, _, err := client.Group.GetByName(cmd.Context(), name)
				if err != nil {
					return fmt.Errorf("failed to resolve group: %w", err)
				}
				groupID = group.ID
			}

			zones, _, err := client.Group.Zones(cmd.Context(), groupID)
			if err != nil {
				return fmt.Errorf("failed to list zones: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(map[string]any{
					"zones": zones,
					"count": len(zones),
				}, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			t := output.New(cmd.OutOrStdout())
			t.AddHeader("ZONE ID", "ZONE NAME", "TYPE")
			noHeader, _ := cmd.Flags().GetBool("no-header")
			t.SetNoHeader(noHeader)
			for _, z := range zones {
				t.AddRow(strconv.Itoa(z.ZoneID), z.ZoneName, z.ZoneType)
			}
			t.Flush()

			return nil
		},
	}

	cmd.Flags().String("name", "", "Group name")
	cmd.Flags().String("id", "", "Group ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")

	return cmd
}
