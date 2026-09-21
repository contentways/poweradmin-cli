// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package permission_templates

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
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

			interactive, _ := cmd.Flags().GetBool("interactive")

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			if interactive {
				base.PrintPreviewNotice(cmd)
				return runInteractiveDelete(cmd, client)
			}

			name, _ := cmd.Flags().GetString("name")
			idStr, _ := cmd.Flags().GetString("id")

			if name == "" && idStr == "" {
				return fmt.Errorf("either --name or --id is required")
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
				if t == nil {
					return fmt.Errorf("permission template %q not found", name)
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

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, map[string]any{
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
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output after deletion")
	cmd.Flags().BoolP("interactive", "i", false, "Interactively select templates to delete (feature preview)")
	cmd.RegisterFlagCompletionFunc("name", base.PermissionTemplateNameCompletion(s))
	return cmd
}

// runInteractiveDelete lists all permission templates, lets the user pick
// zero or more via a multi-select prompt, then hands off to
// deleteSelectedTemplates.
func runInteractiveDelete(cmd *cobra.Command, client *poweradmin.Client) error {
	templates, _, err := client.PermissionTemplate.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to list permission templates: %w", err)
	}
	if len(templates) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no permission templates found")
		return nil
	}

	labels := make([]string, len(templates))
	byLabel := make(map[string]*poweradmin.PermissionTemplate, len(templates))
	for i, t := range templates {
		label := fmt.Sprintf("%s (id %d)", t.Name, t.ID)
		labels[i] = label
		byLabel[label] = t
	}

	selectedLabels, err := base.PromptMultiSelect("Select permission templates to delete", labels)
	if err != nil {
		return err
	}

	selected := make([]*poweradmin.PermissionTemplate, 0, len(selectedLabels))
	for _, label := range selectedLabels {
		selected = append(selected, byLabel[label])
	}

	return deleteSelectedTemplates(cmd, client, selected)
}

// deleteSelectedTemplates confirms and deletes the given permission
// templates. Failures on individual templates do not stop the remaining
// deletions; all errors are collected and returned together at the end.
// Split out from runInteractiveDelete so the deletion/confirmation/
// error-collection logic can be exercised directly in tests without going
// through the multi-select prompt.
func deleteSelectedTemplates(cmd *cobra.Command, client *poweradmin.Client, selected []*poweradmin.PermissionTemplate) error {
	if len(selected) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no permission templates selected, nothing to do")
		return nil
	}

	summary := fmt.Sprintf("The following %d permission template(s) will be deleted:\n", len(selected))
	for _, t := range selected {
		summary += fmt.Sprintf("  - %s (id %d)\n", t.Name, t.ID)
	}
	summary += "\nProceed? [y/N] "
	if !base.Confirm(cmd, summary) {
		return nil
	}

	var errs []error
	for _, t := range selected {
		if _, err := client.PermissionTemplate.Delete(cmd.Context(), t.ID); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "failed to delete permission template %s (id %d): %s\n", t.Name, t.ID, err)
			errs = append(errs, fmt.Errorf("template %s: %w", t.Name, err))
			continue
		}
		fmt.Fprintf(cmd.OutOrStdout(), "deleted permission template %s (id %d)\n", t.Name, t.ID)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to delete %d of %d permission template(s): %w", len(errs), len(selected), errors.Join(errs...))
	}
	return nil
}
