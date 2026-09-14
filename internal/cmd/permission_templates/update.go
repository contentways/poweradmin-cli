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

// NewUpdateCmd returns a new "permission-templates update" command instance.
func NewUpdateCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a permission template",
		Long:  `Update an existing Poweradmin permission template by name or numeric ID.`,
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

			// Resolve ID.
			var tmplID int
			if idStr != "" {
				tmplID, err = strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
			} else {
				t, _, err := client.PermissionTemplate.GetByName(cmd.Context(), name)
				if err != nil {
					return fmt.Errorf("failed to resolve permission template: %w", err)
				}
				tmplID = t.ID
			}

			// Fetch current state to preserve unchanged fields.
			current, _, err := client.PermissionTemplate.GetByID(cmd.Context(), tmplID)
			if err != nil {
				return fmt.Errorf("failed to get permission template: %w", err)
			}

			opts := poweradmin.PermissionTemplateOpts{
				Name:         current.Name,
				Descr:        current.Descr,
				TemplateType: current.TemplateType,
			}
			for _, p := range current.Permissions {
				opts.Permissions = append(opts.Permissions, p.ID)
			}

			if cmd.Flags().Changed("new-name") {
				opts.Name, _ = cmd.Flags().GetString("new-name")
			}
			if cmd.Flags().Changed("description") {
				opts.Descr, _ = cmd.Flags().GetString("description")
			}
			if cmd.Flags().Changed("type") {
				opts.TemplateType, _ = cmd.Flags().GetString("type")
			}
			if cmd.Flags().Changed("permissions") {
				permsStr, _ := cmd.Flags().GetStringSlice("permissions")
				opts.Permissions = nil
				for _, p := range permsStr {
					p = strings.TrimSpace(p)
					if p == "" {
						continue
					}
					id, err := strconv.Atoi(p)
					if err != nil {
						return fmt.Errorf("invalid permission ID %q: %w", p, err)
					}
					opts.Permissions = append(opts.Permissions, id)
				}
			}

			updated, _, err := client.PermissionTemplate.Update(cmd.Context(), tmplID, opts)
			if err != nil {
				return fmt.Errorf("failed to update permission template: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, updated)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "updated permission template %s (id %d)\n", updated.Name, updated.ID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Template name to identify the template")
	cmd.Flags().String("id", "", "Template ID to identify the template")
	cmd.Flags().String("new-name", "", "New template name")
	cmd.Flags().String("description", "", "New description")
	cmd.Flags().String("type", "", "New template type. One of: user|group")
	cmd.Flags().StringSlice("permissions", []string{}, "New permission IDs (replaces existing)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.RegisterFlagCompletionFunc("name", base.PermissionTemplateNameCompletion(s))
	return cmd
}
