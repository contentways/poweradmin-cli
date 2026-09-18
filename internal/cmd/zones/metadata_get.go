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

// NewMetadataGetCmd returns a new "zones metadata-get" command instance.
func NewMetadataGetCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata-get",
		Short: "Get zone metadata by kind",
		Long:  `Get the values stored under a specific metadata kind (e.g. ALLOW-AXFR-FROM) for a Poweradmin zone.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			idStr, _ := cmd.Flags().GetString("id")
			kind, _ := cmd.Flags().GetString("kind")
			if name == "" && idStr == "" {
				return fmt.Errorf("either --name or --id is required")
			}
			if kind == "" {
				return fmt.Errorf("--kind is required")
			}

			s := state.FromContext(cmd.Context())
			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			zone, err := base.ResolveZone(cmd, client)
			if err != nil {
				return err
			}

			m, _, err := client.Zone.GetMetadata(cmd.Context(), zone.ID, kind)
			if err != nil {
				return fmt.Errorf("failed to get metadata: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, schema.ZoneMetadataFromSDK(m))
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Kind: %s\n", m.Kind)
			fmt.Fprintln(cmd.OutOrStdout(), "Values:")
			for _, v := range m.Values {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", v)
			}
			return nil
		},
	}

	cmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("id", "", "Zone ID")
	cmd.Flags().String("kind", "", "Metadata kind (e.g. ALLOW-AXFR-FROM) (required)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
