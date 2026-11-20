// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package server provides CLI commands for the PowerDNS server behind Poweradmin.
// Available subcommands: status.
package server

import (
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewServerCommand builds and returns the "server" subcommand group.
func NewServerCommand(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Inspect the PowerDNS server",
		Long:  `Inspect the PowerDNS server behind Poweradmin (requires Poweradmin 4.5 or newer).`,
	}

	cmd.AddCommand(NewStatusCmd(s))

	return cmd
}
