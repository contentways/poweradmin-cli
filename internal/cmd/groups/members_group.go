// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewMembersGroupCmd returns the "groups members" command group with the
// list, add and remove subcommands.
//
// Called without a subcommand, "groups members" still lists the members of a
// group, as it did before the subcommands existed.
func NewMembersGroupCmd(s *state.State) *cobra.Command {
	cmd := NewMembersCmd(s)
	cmd.Short = "Manage members of a group"
	cmd.Long = `Manage the users of a Poweradmin group.

Without a subcommand the members of the group are listed, like "groups members list".`

	list := NewMembersCmd(s)
	list.Use = "list"

	cmd.AddCommand(list)
	cmd.AddCommand(NewMemberAddCmd(s))
	cmd.AddCommand(NewMemberRemoveCmd(s))
	return cmd
}

// NewZonesGroupCmd returns the "groups zones" command group with the list,
// add and remove subcommands.
//
// Called without a subcommand, "groups zones" still lists the zones of a
// group, as it did before the subcommands existed.
func NewZonesGroupCmd(s *state.State) *cobra.Command {
	cmd := NewZonesCmd(s)
	cmd.Short = "Manage zones of a group"
	cmd.Long = `Manage the zones assigned to a Poweradmin group.

Without a subcommand the zones of the group are listed, like "groups zones list".`

	list := NewZonesCmd(s)
	list.Use = "list"

	cmd.AddCommand(list)
	cmd.AddCommand(NewZoneAddCmd(s))
	cmd.AddCommand(NewZoneRemoveCmd(s))
	return cmd
}
