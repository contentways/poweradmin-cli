// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package permission_templates provides CLI commands for managing Poweradmin
// permission templates via the REST API.
package permission_templates

import (
	"github.com/contentways/poweradmin-cli/v2/internal/state"
	"github.com/spf13/cobra"
)

// NewPermissionTemplatesCommand returns the "permission-templates" subcommand
// with all child commands registered.
func NewPermissionTemplatesCommand(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "permission-templates",
		Aliases: []string{"pt"},
		Short:   "Manage Poweradmin permission templates",
		Long:    `Create, list, update and delete Poweradmin permission templates.`,
	}

	cmd.AddCommand(NewListCmd(s))
	cmd.AddCommand(NewGetCmd(s))
	cmd.AddCommand(NewCreateCmd(s))
	cmd.AddCommand(NewUpdateCmd(s))
	cmd.AddCommand(NewDeleteCmd(s))

	return cmd
}
