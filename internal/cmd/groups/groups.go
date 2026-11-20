// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package groups provides CLI commands for managing groups in Poweradmin.
// Available subcommands: list, get, create, update, delete,
// members (list, add, remove), zones (list, add, remove).
package groups

import (
	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewGroupsCommand builds and returns the "groups" subcommand group.
// All group-related commands are registered here and made available
// under the "poweradmin groups" namespace.
func NewGroupsCommand(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "Manage Poweradmin groups",
		Long:  `Manage Poweradmin groups — list, get, create, update, delete and manage members and zones.`,
	}

	cmd.AddCommand(NewListCmd(s))
	cmd.AddCommand(NewGetCmd(s))
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewUpdateCmd(s))
	cmd.AddCommand(NewDeleteCmd(s))
	cmd.AddCommand(NewMembersGroupCmd(s))
	cmd.AddCommand(NewZonesGroupCmd(s))

	// Flat names from before the member and zone commands were grouped.
	cmd.AddCommand(base.DeprecatedAlias(NewMemberAddCmd(s), "member-add", "groups members add"))
	cmd.AddCommand(base.DeprecatedAlias(NewMemberRemoveCmd(s), "member-remove", "groups members remove"))
	cmd.AddCommand(base.DeprecatedAlias(NewZoneAddCmd(s), "zone-add", "groups zones add"))
	cmd.AddCommand(base.DeprecatedAlias(NewZoneRemoveCmd(s), "zone-remove", "groups zones remove"))

	return cmd
}
