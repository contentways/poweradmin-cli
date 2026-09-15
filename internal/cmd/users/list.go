// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/schema"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/spf13/cobra"
)

// NewListCmd returns a new "users list" command instance.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Long:  `List all users in Poweradmin.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			users, err := client.User.All(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to list users: %w", err)
			}

			sortBy, _ := cmd.Flags().GetString("sort")
			if sortBy != "" {
				sort.Slice(users, func(i, j int) bool {
					switch sortBy {
					case "username":
						return users[i].Username < users[j].Username
					case "email":
						return users[i].Email < users[j].Email
					default: // id
						return users[i].ID < users[j].ID
					}
				})
			}

			activeFilter, _ := cmd.Flags().GetBool("active")
			activeSet := cmd.Flags().Changed("active")
			if activeSet {
				filtered := users[:0]
				for _, u := range users {
					if u.Active == activeFilter {
						filtered = append(filtered, u)
					}
				}
				users = filtered
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.UserListFromSDK(users))
			}

			t := base.NewTable(cmd)
			t.AddHeader("ID", "USERNAME", "EMAIL", "ACTIVE")
			for _, u := range users {
				active := "no"
				activeColor := output.RedCode()
				if u.Active {
					active = "yes"
					activeColor = output.GreenCode()
				}
				t.AddColoredRow(
					output.PlainCell(strconv.Itoa(u.ID)),
					output.PlainCell(u.Username),
					output.PlainCell(u.Email),
					output.Cell(active, activeColor),
				)
			}
			t.Flush()
			return nil
		},
	}

	cmd.Flags().Bool("active", false, "Filter by active status (--active or --active=false)")
	cmd.Flags().String("sort", "", "Sort by field. One of: id|username|email")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")
	return cmd
}
