// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewMetadataGroupCmd returns the "zones metadata" command group with the
// list, get, set and delete subcommands.
//
// Called without a subcommand, "zones metadata" still lists the metadata of a
// zone, as it did before the subcommands existed.
func NewMetadataGroupCmd(s *state.State) *cobra.Command {
	cmd := NewMetadataCmd(s)
	cmd.Short = "Manage zone metadata"
	cmd.Long = `Manage metadata entries (e.g. ALLOW-AXFR-FROM) of a Poweradmin zone.

Without a subcommand the metadata of the zone is listed, like "zones metadata list".`

	list := NewMetadataCmd(s)
	list.Use = "list"

	cmd.AddCommand(list)
	cmd.AddCommand(NewMetadataGetCmd(s))
	cmd.AddCommand(NewMetadataSetCmd(s))
	cmd.AddCommand(NewMetadataDeleteCmd(s))
	return cmd
}
