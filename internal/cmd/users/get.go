// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

import (
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/schema"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewGetCmd returns a new "users get" command instance.
func NewGetCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a user by name or ID",
		Long:  `Get a Poweradmin user by username or numeric ID.`,
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

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, schema.UserFromSDK(user))
			}

			active := output.Red("no")
			if user.Active {
				active = output.Green("yes")
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ID:       %d\n", user.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Username: %s\n", user.Username)
			fmt.Fprintf(cmd.OutOrStdout(), "Email:    %s\n", user.Email)
			fmt.Fprintf(cmd.OutOrStdout(), "Active:   %s\n", active)
			if user.Fullname != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Fullname: %s\n", user.Fullname)
			}
			return nil
		},
	}

	cmd.Flags().String("name", "", "Username")
	cmd.Flags().String("id", "", "User ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
