// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package permission_templates

import (
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewDeleteCmd returns a new "permission-templates delete" command instance.
func NewDeleteCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a permission template",
		Long:  `Delete a Poweradmin permission template by name or numeric ID.`,
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

			var tmplID int
			var tmplName string

			if idStr != "" {
				var err error
				tmplID, err = strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
				tmplName = idStr
			} else {
				t, _, err := client.PermissionTemplate.GetByName(cmd.Context(), name)
				if err != nil {
					return fmt.Errorf("failed to resolve permission template: %w", err)
				}
				tmplID = t.ID
				tmplName = t.Name
			}

			if !base.Confirm(cmd, fmt.Sprintf("Delete permission template %s (id %d)? [y/N] ", tmplName, tmplID)) {
				return nil
			}

			_, err = client.PermissionTemplate.Delete(cmd.Context(), tmplID)
			if err != nil {
				return fmt.Errorf("failed to delete permission template: %w", err)
			}

			if base.IsQuiet(cmd) {
				return nil
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, map[string]any{
					"id":   tmplID,
					"name": tmplName,
				})
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted permission template %s (id %d)\n", tmplName, tmplID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Template name to identify the template")
	cmd.Flags().String("id", "", "Template ID to identify the template")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output after deletion")
	cmd.RegisterFlagCompletionFunc("name", base.PermissionTemplateNameCompletion(s))
	return cmd
}
