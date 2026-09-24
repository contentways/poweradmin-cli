// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package permission_templates

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
	"github.com/spf13/cobra"
)

// NewCreateCmd returns a new "permission-templates create" command instance.
func NewCreateCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a permission template",
		Long:  `Create a new Poweradmin permission template.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			interactive, _ := cmd.Flags().GetBool("interactive")

			name, _ := cmd.Flags().GetString("name")
			descr, _ := cmd.Flags().GetString("description")
			tmplType, _ := cmd.Flags().GetString("type")
			permsStr, _ := cmd.Flags().GetStringSlice("permissions")

			if interactive {
				var err error

				if name == "" {
					name, err = base.PromptString("Template name", "", true)
					if err != nil {
						return err
					}
				}

				if !cmd.Flags().Changed("description") {
					descr, err = base.PromptString("Description", "optional", false)
					if err != nil {
						return err
					}
				}

				if !cmd.Flags().Changed("type") {
					tmplType, err = base.PromptSelect(
						"Template type",
						[]string{"user", "group"},
						tmplType,
					)
					if err != nil {
						return err
					}
				}

				if !cmd.Flags().Changed("permissions") {
					permsStr, err = base.PromptStringSlice(
						"Permission IDs",
						"Comma-separated, e.g. 1,2,3 (optional)",
					)
					if err != nil {
						return err
					}
				}
			}

			if name == "" {
				return fmt.Errorf("--name is required")
			}

			// Parse permission IDs.
			var permissions []int
			for _, p := range permsStr {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				id, err := strconv.Atoi(p)
				if err != nil {
					return fmt.Errorf("invalid permission ID %q: %w", p, err)
				}
				permissions = append(permissions, id)
			}

			if interactive {
				summary := fmt.Sprintf(
					"Create permission template:\n  Name:         %s\n  Description:  %s\n  Type:         %s\n  Permissions:  %v\n\nProceed? [y/N] ",
					name, descr, tmplType, permissions,
				)
				if !base.Confirm(cmd, summary) {
					return nil
				}
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			tmpl, _, err := client.PermissionTemplate.Create(cmd.Context(), poweradmin.PermissionTemplateOpts{
				Name:         name,
				Descr:        descr,
				TemplateType: tmplType,
				Permissions:  permissions,
			})
			if err != nil {
				return fmt.Errorf("failed to create permission template: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, tmpl)
			}

			if base.IsQuiet(cmd) {
				fmt.Fprintln(cmd.OutOrStdout(), tmpl.ID)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "created permission template %s (id %d)\n", tmpl.Name, tmpl.ID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Template name (required)")
	cmd.Flags().String("description", "", "Template description")
	cmd.Flags().String("type", "user", "Template type. One of: user|group")
	cmd.Flags().StringSlice("permissions", []string{}, "Permission IDs to assign (comma-separated or multiple flags)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("quiet", "q", false, "Only print the ID of the created template")
	cmd.Flags().BoolP("interactive", "i", false, "Prompt interactively for missing values (feature preview)")
	return cmd
}
