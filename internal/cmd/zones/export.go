// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"
	"strings"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewExportCmd returns a new "zones export" command instance.
// Exports a DNS zone in BIND zone file format to stdout.
// The zone can be identified by name (--name) or numeric ID (--id).
func NewExportCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export a DNS zone in BIND format",
		Long:  `Export a DNS zone as a BIND-compatible zone file to stdout.`,
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

			records, err := client.Record.All(cmd.Context(), zone.ID)
			if err != nil {
				return fmt.Errorf("failed to get records: %w", err)
			}

			w := cmd.OutOrStdout()

			// Write BIND zone file header.
			fmt.Fprintf(w, "$ORIGIN %s.\n", zone.Name)

			// Find and write SOA record first.
			for _, r := range records {
				if r.Type == "SOA" {
					// SOA content: primary admin serial refresh retry expire minimum
					parts := strings.Fields(r.Content)
					if len(parts) == 7 {
						fmt.Fprintf(w, "$TTL %d\n", r.TTL)
						fmt.Fprintf(w, "%-30s IN  SOA  %s %s (\n", "@", parts[0], parts[1])
						fmt.Fprintf(w, "%-30s         %-12s ; serial\n", "", parts[2])
						fmt.Fprintf(w, "%-30s         %-12s ; refresh\n", "", parts[3])
						fmt.Fprintf(w, "%-30s         %-12s ; retry\n", "", parts[4])
						fmt.Fprintf(w, "%-30s         %-12s ; expire\n", "", parts[5])
						fmt.Fprintf(w, "%-30s         %-12s ; minimum\n", "", parts[6])
						fmt.Fprintf(w, "%-30s         )\n\n", "")
					}
					break
				}
			}

			// Write all non-SOA records.
			for _, r := range records {
				if r.Type == "SOA" {
					continue
				}

				// Convert absolute name to relative — strip zone suffix.
				recName := r.Name
				if recName == zone.Name {
					recName = "@"
				} else if strings.HasSuffix(recName, "."+zone.Name) {
					recName = strings.TrimSuffix(recName, "."+zone.Name)
				}

				fmt.Fprintf(w, "%-30s %-6d IN  %-6s %s\n", recName, r.TTL, r.Type, r.Content)
			}

			return nil
		},
	}

	cmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("id", "", "Zone ID")
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
