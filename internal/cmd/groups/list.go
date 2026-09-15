// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/schema"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewListCmd returns a new "groups list" command instance.
func NewListCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		Long:  `List all groups in Poweradmin.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			groups, err := client.Group.All(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to list groups: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.GroupListFromSDK(groups))
			}

			t := base.NewTable(cmd)
			t.AddHeader("ID", "NAME", "DESCRIPTION", "MEMBERS", "ZONES")
			for _, g := range groups {
				memberColor := ""
				if g.MemberCount > 0 {
					memberColor = output.GreenCode()
				}
				zoneColor := ""
				if g.ZoneCount > 0 {
					zoneColor = output.GreenCode()
				}
				t.AddColoredRow(
					output.PlainCell(strconv.Itoa(g.ID)),
					output.PlainCell(g.Name),
					output.PlainCell(g.Description),
					output.Cell(strconv.Itoa(g.MemberCount), memberColor),
					output.Cell(strconv.Itoa(g.ZoneCount), zoneColor),
				)
			}
			t.Flush()
			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.GroupNameCompletion(s))
	return cmd
}
