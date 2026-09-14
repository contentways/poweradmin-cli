// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package dnssec

import (
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewEnableCmd returns a new "zones dnssec enable" command instance.
func NewEnableCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable DNSSEC for a zone",
		Long:  `Enable DNSSEC for a Poweradmin zone.`,
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

			dnssec, _, err := client.Zone.SetDNSSEC(cmd.Context(), zoneID, true)
			if err != nil {
				return fmt.Errorf("failed to enable DNSSEC: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, dnssec)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "DNSSEC enabled for zone %s\n", name)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Zone name")
	cmd.Flags().String("id", "", "Zone ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
