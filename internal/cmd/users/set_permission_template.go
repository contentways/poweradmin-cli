// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

import (
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewSetPermissionTemplateCmd returns a new "users set-permission-template" command instance.
func NewSetPermissionTemplateCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-permission-template",
		Short: "Assign a permission template to a user",
		Long:  `Assign a permission template to a Poweradmin user by username or numeric ID.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			name, _ := cmd.Flags().GetString("name")
			idStr, _ := cmd.Flags().GetString("id")
			templateIDStr, _ := cmd.Flags().GetString("template-id")

			if name == "" && idStr == "" {
				return fmt.Errorf("either --name or --id is required")
			}
			if templateIDStr == "" {
				return fmt.Errorf("--template-id is required")
			}

			templateID, err := strconv.Atoi(templateIDStr)
			if err != nil {
				return fmt.Errorf("invalid template-id: %w", err)
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			user, err := base.ResolveUser(cmd, client)
			if err != nil {
				return err
			}

			_, err = client.User.SetPermissionTemplate(cmd.Context(), user.ID, templateID)
			if err != nil {
				return fmt.Errorf("failed to set permission template: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, map[string]any{
					"user_id":     user.ID,
					"username":    user.Username,
					"template_id": templateID,
				})
			}

			fmt.Fprintf(cmd.OutOrStdout(), "assigned permission template %d to user %s (id %d)\n", templateID, user.Username, user.ID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Username to identify the user")
	cmd.Flags().String("id", "", "User ID to identify the user")
	cmd.Flags().String("template-id", "", "Permission template ID to assign (required)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
	return cmd
}
