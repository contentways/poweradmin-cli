// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package records provides CLI commands for managing DNS records in Poweradmin.
// Available subcommands: list, create, delete.
package records

import (
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewRecordsCommand builds and returns the "records" subcommand group.
// All record-related commands are registered here and made available
// under the "poweradmin records" namespace.
func NewRecordsCommand(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "records",
		Short: "Manage DNS records",
		Long:  `Manage DNS records in Poweradmin — list, create and delete records.`,
	}

	cmd.AddCommand(NewListCmd(s))
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewDeleteCmd(s))
	cmd.AddCommand(NewUpdateCmd(s))
	cmd.AddCommand(NewGetCmd(s))

	return cmd
}
