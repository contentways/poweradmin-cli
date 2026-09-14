// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package permission_templates

import (
	"fmt"
	"strconv"
	"strings"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewGetCmd returns a new "permission-templates get" command instance.
// Fetches the full permission list by calling GetByID after resolving by name.
func NewGetCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a permission template by name or ID",
		Long:  `Get a Poweradmin permission template by name or numeric ID.`,
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

			var tmpl *poweradmin.PermissionTemplate

			if idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
				tmpl, _, err = client.PermissionTemplate.GetByID(cmd.Context(), id)
				if err != nil {
					return fmt.Errorf("failed to get permission template: %w", err)
				}
			} else {
				// GetByName returns a list result without full permissions.
				// Call GetByID to get the full permission list.
				t, _, err := client.PermissionTemplate.GetByName(cmd.Context(), name)
				if err != nil {
					return fmt.Errorf("failed to get permission template: %w", err)
				}
				tmpl, _, err = client.PermissionTemplate.GetByID(cmd.Context(), t.ID)
				if err != nil {
					return fmt.Errorf("failed to get permission template details: %w", err)
				}
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, tmpl)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "ID:          %d\n", tmpl.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Name:        %s\n", tmpl.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", tmpl.Descr)
			fmt.Fprintf(cmd.OutOrStdout(), "Type:        %s\n", tmpl.TemplateType)
			if len(tmpl.Permissions) > 0 {
				perms := make([]string, len(tmpl.Permissions))
				for i, p := range tmpl.Permissions {
					perms[i] = fmt.Sprintf("%d (%s)", p.ID, p.Name)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Permissions: %s\n", strings.Join(perms, ", "))
			}
			return nil
		},
	}

	cmd.Flags().String("name", "", "Permission template name")
	cmd.Flags().String("id", "", "Permission template ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.RegisterFlagCompletionFunc("name", base.PermissionTemplateNameCompletion(s))
	return cmd
}
