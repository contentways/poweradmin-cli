// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package users provides CLI commands for managing users in Poweradmin.
// Available subcommands: list, get, create, update, delete, set-permission-template.
package users

import (
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewUsersCommand builds and returns the "users" subcommand group.
// All user-related commands are registered here and made available
// under the "poweradmin users" namespace.
func NewUsersCommand(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage Poweradmin users",
		Long:  `Manage Poweradmin users — list, get, create, update, delete and assign permission templates.`,
	}

	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewGetCmd(s))
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewUpdateCmd(s))
	cmd.AddCommand(NewDeleteCmd(s))
	cmd.AddCommand(NewSetPermissionTemplateCmd(s))

	return cmd
}
