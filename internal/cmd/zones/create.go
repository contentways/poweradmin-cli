// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
	"github.com/spf13/cobra"
)

// NewCreateCmd returns a new "zones create" command instance.
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a DNS zone",
		Long:  `Create a new DNS zone in Poweradmin.`,
		// Name can be given positionally, or supplied via --interactive prompt.
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			interactive, _ := cmd.Flags().GetBool("interactive")

			var name string
			if len(args) == 1 {
				name = args[0]
			}

			zoneType, _ := cmd.Flags().GetString("type")
			ttl, _ := cmd.Flags().GetInt("ttl")
			nameservers, _ := cmd.Flags().GetStringSlice("nameserver")

			if interactive {
				var err error

				if name == "" {
					name, err = base.PromptString("Zone name", "e.g. example.com", true)
					if err != nil {
						return err
					}
				}

				if !cmd.Flags().Changed("type") {
					zoneType, err = base.PromptSelect(
						"Zone type",
						[]string{"NATIVE", "MASTER", "SLAVE"},
						zoneType,
					)
					if err != nil {
						return err
					}
				}

				if !cmd.Flags().Changed("nameserver") {
					nameservers, err = base.PromptStringSlice(
						"Nameservers",
						"Comma-separated, e.g. ns1.example.com,ns2.example.com (optional)",
					)
					if err != nil {
						return err
					}
				}

				if len(nameservers) > 0 && !cmd.Flags().Changed("ttl") {
					ttl, err = base.PromptInt("TTL for NS records", "in seconds", ttl)
					if err != nil {
						return err
					}
				}

				summary := fmt.Sprintf(
					"Create zone:\n  Name:        %s\n  Type:        %s\n  Nameservers: %v\n  TTL:         %d\n\nProceed? [y/N] ",
					name, zoneType, nameservers, ttl,
				)
				if !base.Confirm(cmd, summary) {
					return nil
				}
			}

			if name == "" {
				return fmt.Errorf("zone name is required (pass it as an argument or use --interactive)")
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			id, _, err := client.Zone.Create(cmd.Context(), poweradmin.ZoneCreateOpts{
				Name: name,
				Type: poweradmin.ZoneType(zoneType),
			})
			if err != nil {
				return fmt.Errorf("failed to create zone: %w", err)
			}

			for _, ns := range nameservers {
				_, _, err := client.Record.Create(cmd.Context(), id, poweradmin.RecordCreateOpts{
					Name:    name,
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

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, map[string]any{
					"id":          id,
					"name":        name,
					"type":        zoneType,
					"nameservers": nameservers,
					"ttl":         ttl,
				})
			}

			if base.IsQuiet(cmd) {
				fmt.Fprintln(cmd.OutOrStdout(), id)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "created zone %s (id %d)\n", name, id)
			return nil
		},
	}

	cmd.Flags().String("type", "NATIVE", "Zone type. One of: NATIVE|MASTER|SLAVE")
	cmd.Flags().StringSlice("nameserver", []string{}, "Nameserver to add (comma-separated or multiple flags)")
	cmd.Flags().Int("ttl", 3600, "TTL for the created NS records")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("quiet", "q", false, "Only print the ID of the created zone")
	cmd.Flags().BoolP("interactive", "i", false, "Prompt interactively for missing values (feature preview)")
	return cmd
}
