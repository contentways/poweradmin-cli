// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"
	"strings"

	"github.com/contentways/poweradmin-cli/v2/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v2/internal/output"
	"github.com/contentways/poweradmin-cli/v2/internal/schema"
	"github.com/contentways/poweradmin-cli/v2/internal/state"
	"github.com/spf13/cobra"
)

// NewMetadataCmd returns a new "zones metadata" command instance.
func NewMetadataCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "List metadata entries for a DNS zone",
		Long:  `List all metadata entries (e.g. ALLOW-AXFR-FROM) for a Poweradmin zone.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			idStr, _ := cmd.Flags().GetString("id")
			if name == "" && idStr == "" {
				return fmt.Errorf("either --name or --id is required")
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

			metadata, _, err := client.Zone.ListMetadata(cmd.Context(), zone.ID)
			if err != nil {
				return fmt.Errorf("failed to list metadata: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.ZoneMetadataListFromSDK(metadata))
			}

			t := base.NewTable(cmd)
			t.AddHeader("KIND", "VALUES")
			for _, m := range metadata {
				t.AddRow(m.Kind, strings.Join(m.Values, ", "))
			}
			t.Flush()

			return nil
		},
	}

	cmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("id", "", "Zone ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
