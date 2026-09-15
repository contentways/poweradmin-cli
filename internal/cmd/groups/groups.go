// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package groups provides CLI commands for managing groups in Poweradmin.
// Available subcommands: list, get, create, update, delete,
// members, member-add, member-remove, zones, zone-add, zone-remove.
package groups

import (
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
	cmd.AddCommand(NewMembersCmd(s))
	cmd.AddCommand(NewMemberAddCmd(s))
	cmd.AddCommand(NewMemberRemoveCmd(s))
	cmd.AddCommand(NewZonesCmd(s))
	cmd.AddCommand(NewZoneAddCmd(s))
	cmd.AddCommand(NewZoneRemoveCmd(s))

	return cmd
}
