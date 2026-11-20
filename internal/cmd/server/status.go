// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package server

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/schema"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
	"github.com/spf13/cobra"
)

// ErrPowerDNSUnreachable is returned when Poweradmin answers but PowerDNS does not.
var ErrPowerDNSUnreachable = errors.New("PowerDNS is not reachable")

// NewStatusCmd returns a new "server status" command instance.
func NewStatusCmd(_ *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the status of the PowerDNS server",
		Long: `Show whether the PowerDNS server behind Poweradmin is running, its version,
uptime and optionally its metrics and the reachability of the autoprimary servers.

Requires the server_status_view permission (administrators have it implicitly);
--include-slaves additionally requires supermaster_view. The command exits with
a non-zero status when PowerDNS is not reachable, so it can be used in
monitoring checks.`,
		Example: `  poweradmin server status
  poweradmin server status --metrics uptime,udp-queries
  poweradmin server status --show-metrics -o json
  poweradmin server status --include-slaves`,
		RunE: func(cmd *cobra.Command, args []string) error {
			metrics, _ := cmd.Flags().GetStringSlice("metrics")
			includeSlaves, _ := cmd.Flags().GetBool("include-slaves")
			showMetrics, _ := cmd.Flags().GetBool("show-metrics")

			s := state.FromContext(cmd.Context())
			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			status, _, err := client.Server.Status(cmd.Context(), poweradmin.ServerStatusOpts{
				Metrics:       metrics,
				IncludeSlaves: includeSlaves,
			})
			if poweradmin.IsServiceUnavailable(err) {
				return ErrPowerDNSUnreachable
			}
			if err != nil {
				return fmt.Errorf("failed to get server status: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)
			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, schema.ServerStatusFromSDK(status))
			}

			printStatus(cmd, status)
			if showMetrics || len(metrics) > 0 {
				printMetrics(cmd, status)
			}
			if includeSlaves {
				printSlaves(cmd, status)
			}
			return nil
		},
	}

	cmd.Flags().StringSlice("metrics", nil, "Only return these metrics (comma-separated, e.g. uptime,udp-queries)")
	cmd.Flags().Bool("show-metrics", false, "Show all metrics in table output (implied by --metrics)")
	cmd.Flags().Bool("include-slaves", false, "Also probe the autoprimary servers (slower, requires supermaster_view)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")
	return cmd
}

func printStatus(cmd *cobra.Command, s *poweradmin.ServerStatus) {
	out := cmd.OutOrStdout()
	running := "no"
	if s.Running {
		running = "yes"
	}
	uptime := "unknown"
	if s.UptimeSeconds != nil {
		uptime = formatUptime(*s.UptimeSeconds)
	}
	fmt.Fprintf(out, "Running:   %s\n", running)
	fmt.Fprintf(out, "Daemon:    %s\n", s.DaemonType)
	fmt.Fprintf(out, "Version:   %s\n", s.Version)
	fmt.Fprintf(out, "Server ID: %s\n", s.ServerID)
	fmt.Fprintf(out, "Uptime:    %s\n", uptime)
}

func printMetrics(cmd *cobra.Command, s *poweradmin.ServerStatus) {
	fmt.Fprintln(cmd.OutOrStdout())
	t := base.NewTable(cmd)
	t.AddHeader("METRIC", "VALUE")
	for _, name := range slices.Sorted(maps.Keys(s.Metrics)) {
		t.AddRow(name, s.Metrics[name])
	}
	t.Flush()
}

func printSlaves(cmd *cobra.Command, s *poweradmin.ServerStatus) {
	fmt.Fprintln(cmd.OutOrStdout())
	if len(s.Slaves) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No autoprimary servers configured.")
		return
	}
	t := base.NewTable(cmd)
	t.AddHeader("SLAVE", "STATUS", "LAST CHECKED", "ERROR")
	for _, sl := range s.Slaves {
		t.AddRow(sl.IP, sl.Status, deref(sl.LastChecked), deref(sl.Error))
	}
	t.Flush()
}

// formatUptime renders seconds as e.g. "3d 4h 5m".
func formatUptime(seconds int) string {
	d, h, m := seconds/86400, seconds%86400/3600, seconds%3600/60
	switch {
	case d > 0:
		return fmt.Sprintf("%dd %dh %dm", d, h, m)
	case h > 0:
		return fmt.Sprintf("%dh %dm", h, m)
	default:
		return fmt.Sprintf("%dm %ds", m, seconds%60)
	}
}

func deref(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}
