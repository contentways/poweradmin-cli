// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/schema"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewGetCmd returns a new "zones get" command instance.
func NewGetCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a DNS zone by name or ID",
		Long:  `Get a DNS zone by name or ID from Poweradmin.`,
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

			zone, err := base.ResolveZone(cmd, client)
			if err != nil {
				return err
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			records, err := client.Record.All(cmd.Context(), zone.ID)
			if err != nil {
				return fmt.Errorf("failed to get records: %w", err)
			}

			var nameservers []string
			for _, r := range records {
				if r.Type == "NS" {
					nameservers = append(nameservers, r.Content)
				}
			}

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, schema.ZoneWithNameservers{
					Zone:        schema.ZoneFromSDK(zone),
					Nameservers: nameservers,
				})
			}

			fmt.Fprintf(cmd.OutOrStdout(), "ID:   %d\n", zone.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Name: %s\n", zone.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Type: %s\n", zone.Type)
			if zone.Masters != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Masters: %s\n", zone.Masters)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Nameservers:")
			for _, r := range records {
				if r.Type == "NS" {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", r.Content)
				}
			}
			return nil
		},
	}

	cmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("id", "", "Zone ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
