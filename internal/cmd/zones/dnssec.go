// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/schema"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
	"github.com/spf13/cobra"
)

// NewDNSSECCmd returns the "zones dnssec" command group.
func NewDNSSECCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dnssec",
		Short: "Manage DNSSEC of a zone",
		Long: `Sign and unsign zones, show their DS records, manage their DNSSEC keys and
rectify them.

Key management and rectify require Poweradmin 4.5 or newer. Changing DNSSEC
requires the zone_dnssec_manage_own permission for the zone.`,
	}

	cmd.AddCommand(NewDNSSECStatusCmd(s))
	cmd.AddCommand(NewDNSSECEnableCmd(s))
	cmd.AddCommand(NewDNSSECDisableCmd(s))
	cmd.AddCommand(NewDNSSECRectifyCmd(s))
	cmd.AddCommand(NewDNSSECKeysCmd(s))

	return cmd
}

// addZoneFlags registers the --name/--id flags used to select the zone.
func addZoneFlags(cmd *cobra.Command, s *state.State) {
	cmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("id", "", "Zone ID")
	cmd.RegisterFlagCompletionFunc("name", base.ZoneNameCompletion(s))
}

// resolveDNSSECZone checks the zone flags, creates the client and resolves the zone.
func resolveDNSSECZone(cmd *cobra.Command) (*poweradmin.Client, *poweradmin.Zone, error) {
	name, _ := cmd.Flags().GetString("name")
	idStr, _ := cmd.Flags().GetString("id")
	if name == "" && idStr == "" {
		return nil, nil, fmt.Errorf("either --name or --id is required")
	}

	s := state.FromContext(cmd.Context())
	client, err := s.Client()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client: %w", err)
	}

	zone, err := base.ResolveZone(cmd, client)
	if err != nil {
		return nil, nil, err
	}
	return client, zone, nil
}

// printDNSSEC writes the DNSSEC status in the format chosen with --output.
func printDNSSEC(cmd *cobra.Command, zone *poweradmin.Zone, d *poweradmin.ZoneDNSSEC) error {
	outputStr, _ := cmd.Flags().GetString("output")
	outputFmt := output.ParseFormat(outputStr)
	if outputFmt.IsStructured() {
		return base.PrintFormatted(cmd, outputFmt, schema.ZoneDNSSECFromSDK(zone.Name, d))
	}

	out := cmd.OutOrStdout()
	signing := "unsigned"
	if d.Enabled {
		signing = "signed"
	}
	if d.Presigned {
		signing += " (presigned, managed at the primary)"
	}
	fmt.Fprintf(out, "Zone:   %s\n", zone.Name)
	fmt.Fprintf(out, "DNSSEC: %s\n", signing)
	if len(d.DSRecords) > 0 {
		fmt.Fprintln(out, "DS records:")
		for _, ds := range d.DSRecords {
			fmt.Fprintf(out, "  %s. IN DS %d %d %d %s\n", zone.Name, ds.KeyTag, ds.Algorithm, ds.DigestType, ds.Digest)
		}
	}
	return nil
}

// NewDNSSECStatusCmd returns a new "zones dnssec status" command instance.
func NewDNSSECStatusCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the DNSSEC status and DS records of a zone",
		Example: `  poweradmin zones dnssec status --name example.com
  poweradmin zones dnssec status --name example.com -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, zone, err := resolveDNSSECZone(cmd)
			if err != nil {
				return err
			}
			d, _, err := client.Zone.GetDNSSEC(cmd.Context(), zone.ID)
			if err != nil {
				return fmt.Errorf("failed to get DNSSEC status: %w", err)
			}
			return printDNSSEC(cmd, zone, d)
		},
	}

	addZoneFlags(cmd, s)
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	return cmd
}

// NewDNSSECEnableCmd returns a new "zones dnssec enable" command instance.
func NewDNSSECEnableCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Sign a zone",
		Long: `Sign a zone with DNSSEC and show its DS records.

The zone is only validated by resolvers once the DS records are published in
the parent zone (usually at the registrar).`,
		Example: `  poweradmin zones dnssec enable --name example.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, zone, err := resolveDNSSECZone(cmd)
			if err != nil {
				return err
			}
			d, _, err := client.Zone.SetDNSSEC(cmd.Context(), zone.ID, true)
			if err != nil {
				return fmt.Errorf("failed to enable DNSSEC: %w", err)
			}
			return printDNSSEC(cmd, zone, d)
		},
	}

	addZoneFlags(cmd, s)
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	return cmd
}

// NewDNSSECDisableCmd returns a new "zones dnssec disable" command instance.
func NewDNSSECDisableCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Unsign a zone",
		Long: `Remove DNSSEC from a zone and delete its keys.

Remove the DS records from the parent zone first; otherwise validating
resolvers will reject the zone.`,
		Example: `  poweradmin zones dnssec disable --name example.com --yes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, zone, err := resolveDNSSECZone(cmd)
			if err != nil {
				return err
			}
			if !base.Confirm(cmd, fmt.Sprintf("Unsign zone %s (id %d) and delete its DNSSEC keys? [y/N] ", zone.Name, zone.ID)) {
				return nil
			}
			d, _, err := client.Zone.SetDNSSEC(cmd.Context(), zone.ID, false)
			if err != nil {
				return fmt.Errorf("failed to disable DNSSEC: %w", err)
			}
			return printDNSSEC(cmd, zone, d)
		},
	}

	addZoneFlags(cmd, s)
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	return cmd
}

// NewDNSSECRectifyCmd returns a new "zones dnssec rectify" command instance.
func NewDNSSECRectifyCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rectify",
		Short: "Rectify a signed zone",
		Long: `Recalculate the DNSSEC ordering and auth fields of a signed zone
(like pdnsutil rectify-zone). Not possible for unsigned, presigned or
secondary zones. Requires Poweradmin 4.5 or newer.`,
		Example: `  poweradmin zones dnssec rectify --name example.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, zone, err := resolveDNSSECZone(cmd)
			if err != nil {
				return err
			}
			if _, err := client.DNSSEC.Rectify(cmd.Context(), zone.ID); err != nil {
				return fmt.Errorf("failed to rectify zone: %w", err)
			}
			if base.IsQuiet(cmd) {
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "rectified zone %s\n", zone.Name)
			return nil
		},
	}

	addZoneFlags(cmd, s)
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output")
	return cmd
}
