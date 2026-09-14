// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package dnssec

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewGetCmd returns a new "zones dnssec get" command instance.
func NewGetCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get DNSSEC status for a zone",
		Long:  `Get the DNSSEC status, DS records and DNSKey for a Poweradmin zone.`,
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
			}

			dnssec, _, err := client.Zone.GetDNSSEC(cmd.Context(), zoneID)
			if err != nil {
				return fmt.Errorf("failed to get DNSSEC status: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, dnssec)
			}

			enabled := "disabled"
			if dnssec.Enabled {
				enabled = "enabled"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "DNSSEC: %s\n", enabled)

			if dnssec.DNSKey != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "DNSKey: %s\n", *dnssec.DNSKey)
			}

			if len(dnssec.DSRecords) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "DS Records:\n")
				for _, r := range dnssec.DSRecords {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s %d %d %d %s\n",
						name,
						r.KeyTag,
						r.Algorithm,
						r.DigestType,
						strings.ToUpper(r.Digest),
					)
				}
			}
			return nil
		},
	}

	cmd.Flags().String("name", "", "Zone name")
	cmd.Flags().String("id", "", "Zone ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
