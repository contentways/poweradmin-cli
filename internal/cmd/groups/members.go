// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewMembersCmd returns a new "groups members" command instance.
func NewMembersCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "members",
		Short: "List members of a group",
		Long:  `List all members of a Poweradmin group.`,
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

			var groupID int
			if idStr != "" {
				groupID, err = strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
			} else {
				group, _, err := client.Group.GetByName(cmd.Context(), name)
				if err != nil {
					return fmt.Errorf("failed to resolve group: %w", err)
				}
				groupID = group.ID
			}

			members, _, err := client.Group.Members(cmd.Context(), groupID)
			if err != nil {
				return fmt.Errorf("failed to list members: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(map[string]any{
					"members": members,
					"count":   len(members),
				}, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			t := output.New(cmd.OutOrStdout())
			t.AddHeader("USER ID", "USERNAME", "FULLNAME")
			noHeader, _ := cmd.Flags().GetBool("no-header")
			t.SetNoHeader(noHeader)
			for _, m := range members {
				t.AddRow(strconv.Itoa(m.UserID), m.Username)
			}
			t.Flush()

			return nil
		},
	}

	cmd.Flags().String("name", "", "Group name")
	cmd.Flags().String("id", "", "Group ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")
	// Register shell completion for --name flag.
	cmd.RegisterFlagCompletionFunc("name", base.GroupNameCompletion(s))

	return cmd
}
