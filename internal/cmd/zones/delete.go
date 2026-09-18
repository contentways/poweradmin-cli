// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
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

			name, _ := cmd.Flags().GetString("name")
			idStr, _ := cmd.Flags().GetString("id")

			if name == "" && idStr == "" {
				return fmt.Errorf("either --name or --id is required")
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
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
				zoneID = zone.ID
				name = zone.Name
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
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
