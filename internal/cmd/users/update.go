// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

import (
	"fmt"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/schema"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewUpdateCmd returns a new "users update" command instance.
func NewUpdateCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a user",
		Long:  `Update an existing Poweradmin user by username or numeric ID.`,
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

			opts := poweradmin.UserUpdateOpts{}
			if cmd.Flags().Changed("email") {
				opts.Email, _ = cmd.Flags().GetString("email")
			}
			if cmd.Flags().Changed("fullname") {
				opts.Fullname, _ = cmd.Flags().GetString("fullname")
			}
			if cmd.Flags().Changed("password") {
				opts.Password, _ = cmd.Flags().GetString("password")
			}
			if cmd.Flags().Changed("active") {
				active, _ := cmd.Flags().GetBool("active")
				opts.Active = &active
			}

			updated, _, err := client.User.Update(cmd.Context(), user.ID, opts)
			if err != nil {
				return fmt.Errorf("failed to update user: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.UserFromSDK(updated))
			}

			fmt.Fprintf(cmd.OutOrStdout(), "updated user %s (id %d)\n", updated.Username, updated.ID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Username to identify the user")
	cmd.Flags().String("id", "", "User ID to identify the user")
	cmd.Flags().String("email", "", "New email address")
	cmd.Flags().String("fullname", "", "New full name")
	cmd.Flags().String("password", "", "New password")
	cmd.Flags().Bool("active", true, "Whether the user is active")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
