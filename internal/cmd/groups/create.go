// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"fmt"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewCreateCmd returns a new "groups create" command instance.
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new group",
		Long:  `Create a new group in Poweradmin.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			name, _ := cmd.Flags().GetString("name")
			description, _ := cmd.Flags().GetString("description")
			permTemplID, _ := cmd.Flags().GetInt("perm-template-id")

			if name == "" {
				return fmt.Errorf("--name is required")
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			id, _, err := client.Group.Create(cmd.Context(), poweradmin.GroupCreateOpts{
				Name:        name,
				Description: description,
				PermTemplID: permTemplID,
			})
			if err != nil {
				return fmt.Errorf("failed to create group: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, map[string]any{
					"id":   id,
					"name": name,
				})
			}

			if base.IsQuiet(cmd) {
				fmt.Fprintln(cmd.OutOrStdout(), id)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "created group %s (id %d)\n", name, id)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Group name (required)")
	cmd.Flags().String("description", "", "Group description")
	cmd.Flags().Int("perm-template-id", 0, "Permission template ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().BoolP("quiet", "q", false, "Only print the ID of the created group")
	return cmd
}
