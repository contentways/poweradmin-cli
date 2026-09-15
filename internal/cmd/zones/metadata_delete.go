// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewMetadataDeleteCmd returns a new "zones metadata-delete" command instance.
func NewMetadataDeleteCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata-delete",
		Short: "Delete zone metadata for a kind",
		Long:  `Delete all values for a metadata kind (e.g. ALLOW-AXFR-FROM) on a Poweradmin zone.`,
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

			if !base.Confirm(cmd, fmt.Sprintf("Delete metadata %s on zone %s (id %d)? [y/N] ", kind, zone.Name, zone.ID)) {
				return nil
			}

			if _, err := client.Zone.DeleteMetadata(cmd.Context(), zone.ID, kind); err != nil {
				return fmt.Errorf("failed to delete metadata: %w", err)
			}

			if base.IsQuiet(cmd) {
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted metadata %s on zone %s\n", kind, zone.Name)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("id", "", "Zone ID")
	cmd.Flags().String("kind", "", "Metadata kind (e.g. ALLOW-AXFR-FROM) (required)")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output after deletion")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
