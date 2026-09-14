// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package dnssec provides CLI commands for managing DNSSEC on Poweradmin zones.
package dnssec

import (
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewDNSSECCommand returns the "zones dnssec" subcommand with all child commands.
func NewDNSSECCommand(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dnssec",
		Short: "Manage DNSSEC for a DNS zone",
		Long:  `Get DNSSEC status, enable or disable DNSSEC for a Poweradmin zone.`,
	}

	cmd.AddCommand(NewGetCmd(s))
	cmd.AddCommand(NewEnableCmd(s))
	cmd.AddCommand(NewDisableCmd(s))

	return cmd
}
