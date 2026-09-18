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

// NewListCmd returns a new "permission-templates list" command instance.
func NewListCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all permission templates",
		Long:  `List all permission templates in Poweradmin.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			templates, _, err := client.PermissionTemplate.List(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to list permission templates: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt.IsStructured() {
				type templateJSON struct {
					ID           int    `json:"id"`
					Name         string `json:"name"`
					Description  string `json:"description,omitempty"`
					TemplateType string `json:"template_type"`
				}
				type listJSON struct {
					Templates []templateJSON `json:"permission_templates"`
					Count     int            `json:"count"`
				}
				list := listJSON{Count: len(templates)}
				for _, t := range templates {
					list.Templates = append(list.Templates, templateJSON{
						ID:           t.ID,
						Name:         t.Name,
						Description:  t.Descr,
						TemplateType: t.TemplateType,
					})
				}
				return base.PrintFormatted(cmd, outputFmt, list)
			}

			t := base.NewTable(cmd)
			t.AddHeader("ID", "NAME", "DESCRIPTION", "TYPE")
			for _, tmpl := range templates {
				t.AddColoredRow(
					output.PlainCell(strconv.Itoa(tmpl.ID)),
					output.PlainCell(tmpl.Name),
					output.PlainCell(tmpl.Descr),
					output.Cell(tmpl.TemplateType, output.CyanCode()),
				)
			}
			t.Flush()
			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")
	return cmd
}
