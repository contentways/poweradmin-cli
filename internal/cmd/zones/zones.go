// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package zones provides CLI commands for managing DNS zones in Poweradmin.
// Available subcommands: list, get, create, delete, export, import, metadata.
package zones

import (
	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
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
		Long:  `Manage DNS zones in Poweradmin — list, get, create, delete, export and import zones, and manage their metadata.`,
	}

	cmd.AddCommand(NewListCmd(s))
	cmd.AddCommand(NewGetCmd(s))
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewDeleteCmd(s))
	cmd.AddCommand(NewExportCmd(s))
	cmd.AddCommand(NewImportCmd(s))
	cmd.AddCommand(NewMetadataGroupCmd(s))

	// Flat names from before the metadata commands were grouped.
	cmd.AddCommand(base.DeprecatedAlias(NewMetadataGetCmd(s), "metadata-get", "zones metadata get"))
	cmd.AddCommand(base.DeprecatedAlias(NewMetadataSetCmd(s), "metadata-set", "zones metadata set"))
	cmd.AddCommand(base.DeprecatedAlias(NewMetadataDeleteCmd(s), "metadata-delete", "zones metadata delete"))

	return cmd
}
