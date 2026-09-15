// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package zones provides CLI commands for managing DNS zones in Poweradmin.
// Available subcommands: list, get, create, delete.
package zones

import (
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewZonesCommand builds and returns the "zones" subcommand group.
// All zone-related commands are registered here and made available
// under the "poweradmin zones" namespace.
func NewZonesCommand(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zones",
		Short: "Manage DNS zones",
		Long:  `Manage DNS zones in Poweradmin — list, get, create and delete zones.`,
	}

	cmd.AddCommand(NewListCmd(s))
	cmd.AddCommand(NewGetCmd(s))
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewDeleteCmd(s))
	cmd.AddCommand(NewExportCmd(s))
	cmd.AddCommand(NewImportCmd(s))
	cmd.AddCommand(NewMetadataCmd(s))
	cmd.AddCommand(NewMetadataGetCmd(s))
	cmd.AddCommand(NewMetadataSetCmd(s))
	cmd.AddCommand(NewMetadataDeleteCmd(s))

	return cmd
}
