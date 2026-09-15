// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewDeleteCmd returns a new "groups delete" command instance.
func NewDeleteCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a group",
		Long:  `Delete a Poweradmin group by name or numeric ID.`,
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

			if !base.Confirm(cmd, fmt.Sprintf("Delete group %s (id %d)? [y/N] ", group.Name, group.ID)) {
				return nil
			}

			_, err = client.Group.Delete(cmd.Context(), group.ID)
			if err != nil {
				return fmt.Errorf("failed to delete group: %w", err)
			}

			if base.IsQuiet(cmd) {
				return nil
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, map[string]any{
					"id":   group.ID,
					"name": group.Name,
				})
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted group %s (id %d)\n", group.Name, group.ID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Group name to identify the group")
	cmd.Flags().String("id", "", "Group ID to identify the group")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output after deletion")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.GroupNameCompletion(s))
	return cmd
}
