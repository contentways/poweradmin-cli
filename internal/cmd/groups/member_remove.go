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

// NewMemberRemoveCmd returns a new "groups member-remove" command instance.
func NewMemberRemoveCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member-remove",
		Short: "Remove a user from a group",
		Long:  `Remove a user from a Poweradmin group by group and user ID.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			groupIDStr, _ := cmd.Flags().GetString("group-id")
			userIDStr, _ := cmd.Flags().GetString("user-id")

			if groupIDStr == "" {
				return fmt.Errorf("--group-id is required")
			}
			if userIDStr == "" {
				return fmt.Errorf("--user-id is required")
			}

			groupID, err := strconv.Atoi(groupIDStr)
			if err != nil {
				return fmt.Errorf("invalid group-id: %w", err)
			}
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				return fmt.Errorf("invalid user-id: %w", err)
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			_, err = client.Group.RemoveMember(cmd.Context(), groupID, userID)
			if err != nil {
				return fmt.Errorf("failed to remove member: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "removed user %d from group %d\n", userID, groupID)
			return nil
		},
	}

	cmd.Flags().String("group-id", "", "Group ID (required)")
	cmd.Flags().String("user-id", "", "User ID to remove (required)")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.GroupNameCompletion(s))
	return cmd
}
