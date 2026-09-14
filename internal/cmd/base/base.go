// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package base provides shared helpers for CLI commands.
// It reduces boilerplate in command implementations by centralising
// common patterns: JSON output, table creation, delete confirmation,
// zone/user/group resolution by name or ID, quiet mode and shell completion.
package base

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// PrintJSON marshals v to indented JSON and writes it to cmd's stdout.
func PrintJSON(cmd *cobra.Command, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}

// NewTable creates a new output.Table writing to cmd's stdout.
// It reads the --no-header flag and configures the table accordingly.
func NewTable(cmd *cobra.Command) *output.Table {
	t := output.New(cmd.OutOrStdout())
	noHeader, _ := cmd.Flags().GetBool("no-header")
	t.SetNoHeader(noHeader)
	return t
}

// IsQuiet returns true if the --quiet flag is set.
func IsQuiet(cmd *cobra.Command) bool {
	q, _ := cmd.Flags().GetBool("quiet")
	return q
}

// Confirm prints a confirmation prompt and reads the user's response.
// Returns true if the user confirmed with "y" or "Y".
// Returns true immediately if the --yes flag is set.
func Confirm(cmd *cobra.Command, msg string) bool {
	yes, _ := cmd.Flags().GetBool("yes")
	if yes {
		return true
	}
	fmt.Fprint(cmd.OutOrStdout(), msg)
	var confirm string
	fmt.Fscan(os.Stdin, &confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
		return false
	}
	return true
}

// ResolveZone returns a Zone identified by --name or --id flag.
// If --id is provided it takes precedence over --name.
func ResolveZone(cmd *cobra.Command, client *poweradmin.Client) (*poweradmin.Zone, error) {
	name, _ := cmd.Flags().GetString("name")
	idStr, _ := cmd.Flags().GetString("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		zone, _, err := client.Zone.GetByID(cmd.Context(), id)
		if err != nil {
			return nil, fmt.Errorf("failed to get zone: %w", err)
		}
		return zone, nil
	}

	zone, _, err := client.Zone.GetByName(cmd.Context(), name)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve zone: %w", err)
	}
	return zone, nil
}

// ResolveUser returns a User identified by --name or --id flag.
func ResolveUser(cmd *cobra.Command, client *poweradmin.Client) (*poweradmin.User, error) {
	name, _ := cmd.Flags().GetString("name")
	idStr, _ := cmd.Flags().GetString("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		user, _, err := client.User.GetByID(cmd.Context(), id)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		return user, nil
	}

	user, _, err := client.User.GetByName(cmd.Context(), name)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve user: %w", err)
	}
	return user, nil
}

// ResolveGroup returns a Group identified by --name or --id flag.
func ResolveGroup(cmd *cobra.Command, client *poweradmin.Client) (*poweradmin.Group, error) {
	name, _ := cmd.Flags().GetString("name")
	idStr, _ := cmd.Flags().GetString("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		group, _, err := client.Group.GetByID(cmd.Context(), id)
		if err != nil {
			return nil, fmt.Errorf("failed to get group: %w", err)
		}
		return group, nil
	}

	group, _, err := client.Group.GetByName(cmd.Context(), name)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve group: %w", err)
	}
	return group, nil
}

// ResolveZoneID returns the numeric zone ID from --zone-name or --zone-id flag.
func ResolveZoneID(cmd *cobra.Command, client *poweradmin.Client) (int, error) {
	zoneName, _ := cmd.Flags().GetString("zone-name")
	zoneIDStr, _ := cmd.Flags().GetString("zone-id")

	if zoneIDStr != "" {
		id, err := strconv.Atoi(zoneIDStr)
		if err != nil {
			return 0, fmt.Errorf("invalid zone-id: %w", err)
		}
		return id, nil
	}

	zone, _, err := client.Zone.GetByName(cmd.Context(), zoneName)
	if err != nil {
		return 0, fmt.Errorf("failed to resolve zone: %w", err)
	}
	return zone.ID, nil
}

// ZoneNameCompletion returns a Cobra completion function that fetches zone names
// from the Poweradmin API. Used for --name and --zone-name flag completion.
// The State is passed directly to avoid depending on the command context,
// which is not available during shell completion.
func ZoneNameCompletion(s *state.State) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := s.Client()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		zones, err := client.Zone.All(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var names []string
		for _, z := range zones {
			if strings.HasPrefix(z.Name, toComplete) {
				names = append(names, z.Name)
			}
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}

// UserNameCompletion returns a Cobra completion function that fetches usernames
// from the Poweradmin API. Used for --name flag completion on user commands.
func UserNameCompletion(s *state.State) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := s.Client()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		users, err := client.User.All(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var names []string
		for _, u := range users {
			if strings.HasPrefix(u.Username, toComplete) {
				names = append(names, u.Username)
			}
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}

// GroupNameCompletion returns a Cobra completion function that fetches group names
// from the Poweradmin API. Used for --name flag completion on group commands.
func GroupNameCompletion(s *state.State) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := s.Client()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		groups, err := client.Group.All(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var names []string
		for _, g := range groups {
			if strings.HasPrefix(g.Name, toComplete) {
				names = append(names, g.Name)
			}
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}

// PermissionTemplateNameCompletion returns a Cobra completion function that
// fetches permission template names from the Poweradmin API.
func PermissionTemplateNameCompletion(s *state.State) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := s.Client()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		templates, _, err := client.PermissionTemplate.List(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var names []string
		for _, t := range templates {
			if strings.HasPrefix(t.Name, toComplete) {
				names = append(names, t.Name)
			}
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}
