// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewCreateCmd returns a new "zones create" command instance.
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a DNS zone",
		Long:  `Create a new DNS zone in Poweradmin.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			zoneType, _ := cmd.Flags().GetString("type")
			ttl, _ := cmd.Flags().GetInt("ttl")

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			id, _, err := client.Zone.Create(cmd.Context(), poweradmin.ZoneCreateOpts{
				Name: args[0],
				Type: poweradmin.ZoneType(zoneType),
			})
			if err != nil {
				return fmt.Errorf("failed to create zone: %w", err)
			}

			nameservers, _ := cmd.Flags().GetStringSlice("nameserver")
			for _, ns := range nameservers {
				_, _, err := client.Record.Create(cmd.Context(), id, poweradmin.RecordCreateOpts{
					Name:    args[0],
					Type:    "NS",
					Content: ns,
					TTL:     ttl,
				})
				if err != nil {
					return fmt.Errorf("failed to create NS record for %s: %w", ns, err)
				}
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, map[string]any{
					"id":          id,
					"name":        args[0],
					"type":        zoneType,
					"nameservers": nameservers,
					"ttl":         ttl,
				})
			}

			if base.IsQuiet(cmd) {
				fmt.Fprintln(cmd.OutOrStdout(), id)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "created zone %s (id %d)\n", args[0], id)
			return nil
		},
	}

	cmd.Flags().String("type", "NATIVE", "Zone type. One of: NATIVE|MASTER|SLAVE")
	cmd.Flags().StringSlice("nameserver", []string{}, "Nameserver to add (comma-separated or multiple flags)")
	cmd.Flags().Int("ttl", 3600, "TTL for the created NS records")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().BoolP("quiet", "q", false, "Only print the ID of the created zone")
	return cmd
}
