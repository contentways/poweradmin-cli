// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

import (
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewDeleteCmd returns a new "users delete" command instance.
func NewDeleteCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a user",
		Long:  `Delete a Poweradmin user by username or numeric ID.`,
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

			user, err := base.ResolveUser(cmd, client)
			if err != nil {
				return err
			}

			if !base.Confirm(cmd, fmt.Sprintf("Delete user %s (id %d)? [y/N] ", user.Username, user.ID)) {
				return nil
			}

			_, err = client.User.Delete(cmd.Context(), user.ID)
			if err != nil {
				return fmt.Errorf("failed to delete user: %w", err)
			}

			if base.IsQuiet(cmd) {
				return nil
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, map[string]any{
					"id":       user.ID,
					"username": user.Username,
				})
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted user %s (id %d)\n", user.Username, user.ID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Username to identify the user")
	cmd.Flags().String("id", "", "User ID to identify the user")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output after deletion")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
