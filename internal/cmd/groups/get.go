// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/schema"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewGetCmd returns a new "groups get" command instance.
func NewGetCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a group by name or ID",
		Long:  `Get a Poweradmin group by name or numeric ID.`,
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

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.GroupFromSDK(group))
			}

			fmt.Fprintf(cmd.OutOrStdout(), "ID:          %d\n", group.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Name:        %s\n", group.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", group.Description)
			fmt.Fprintf(cmd.OutOrStdout(), "Members:     %d\n", group.MemberCount)
			fmt.Fprintf(cmd.OutOrStdout(), "Zones:       %d\n", group.ZoneCount)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Group name")
	cmd.Flags().String("id", "", "Group ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.GroupNameCompletion(s))
	return cmd
}
