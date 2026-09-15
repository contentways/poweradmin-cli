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

// NewMetadataSetCmd returns a new "zones metadata-set" command instance.
func NewMetadataSetCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata-set",
		Short: "Set (replace) zone metadata for a kind",
		Long:  `Create or replace all values for a metadata kind (e.g. ALLOW-AXFR-FROM) on a Poweradmin zone.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			idStr, _ := cmd.Flags().GetString("id")
			kind, _ := cmd.Flags().GetString("kind")
			values, _ := cmd.Flags().GetStringSlice("values")
			if name == "" && idStr == "" {
				return fmt.Errorf("either --name or --id is required")
			}
			if kind == "" {
				return fmt.Errorf("--kind is required")
			}
			if len(values) == 0 {
				return fmt.Errorf("--values is required")
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

			if _, err := client.Zone.SetMetadata(cmd.Context(), zone.ID, kind, values); err != nil {
				return fmt.Errorf("failed to set metadata: %w", err)
			}

			if base.IsQuiet(cmd) {
				return nil
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.ZoneMetadata{Kind: kind, Values: values})
			}

			fmt.Fprintf(cmd.OutOrStdout(), "set metadata %s (%d value(s)) on zone %s\n", kind, len(values), zone.Name)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("id", "", "Zone ID")
	cmd.Flags().String("kind", "", "Metadata kind (e.g. ALLOW-AXFR-FROM) (required)")
	cmd.Flags().StringSlice("values", []string{}, "Metadata values (comma-separated or multiple flags) (required)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
