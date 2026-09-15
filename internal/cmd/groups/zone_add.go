// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewZoneAddCmd returns a new "groups zone-add" command instance.
func NewZoneAddCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone-add",
		Short: "Add a zone to a group",
		Long:  `Associate a zone with a Poweradmin group.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			groupIDStr, _ := cmd.Flags().GetString("group-id")
			zoneIDStr, _ := cmd.Flags().GetString("zone-id")

			if groupIDStr == "" {
				return fmt.Errorf("--group-id is required")
			}
			if zoneIDStr == "" {
				return fmt.Errorf("--zone-id is required")
			}

			groupID, err := strconv.Atoi(groupIDStr)
			if err != nil {
				return fmt.Errorf("invalid group-id: %w", err)
			}
			zoneID, err := strconv.Atoi(zoneIDStr)
			if err != nil {
				return fmt.Errorf("invalid zone-id: %w", err)
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			_, err = client.Group.AddZone(cmd.Context(), groupID, zoneID)
			if err != nil {
				return fmt.Errorf("failed to add zone to group: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "added zone %d to group %d\n", zoneID, groupID)
			return nil
		},
	}

	cmd.Flags().String("group-id", "", "Group ID (required)")
	cmd.Flags().String("zone-id", "", "Zone ID to add (required)")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.GroupNameCompletion(s))
	return cmd
}
