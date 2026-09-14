// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"fmt"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/schema"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewUpdateCmd returns a new "groups update" command instance.
func NewUpdateCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a group",
		Long:  `Update an existing Poweradmin group by name or numeric ID.`,
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

			group, err := base.ResolveGroup(cmd, client)
			if err != nil {
				return err
			}

			opts := poweradmin.GroupUpdateOpts{}
			if cmd.Flags().Changed("new-name") {
				opts.Name, _ = cmd.Flags().GetString("new-name")
			}
			if cmd.Flags().Changed("description") {
				desc, _ := cmd.Flags().GetString("description")
				opts.Description = &desc
			}

			updated, _, err := client.Group.Update(cmd.Context(), group.ID, opts)
			if err != nil {
				return fmt.Errorf("failed to update group: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.GroupFromSDK(updated))
			}

			fmt.Fprintf(cmd.OutOrStdout(), "updated group %s (id %d)\n", updated.Name, updated.ID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Group name to identify the group")
	cmd.Flags().String("id", "", "Group ID to identify the group")
	cmd.Flags().String("new-name", "", "New group name")
	cmd.Flags().String("description", "", "New description")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.GroupNameCompletion(s))
	return cmd
}
