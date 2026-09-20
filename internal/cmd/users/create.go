// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

import (
	"fmt"
	"os"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewCreateCmd returns a new "users create" command instance.
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Long:  `Create a new user in Poweradmin.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			interactive, _ := cmd.Flags().GetBool("interactive")

			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			email, _ := cmd.Flags().GetString("email")
			fullname, _ := cmd.Flags().GetString("fullname")
			active, _ := cmd.Flags().GetBool("active")

			if interactive {
				var err error

				if username == "" {
					username, err = base.PromptString("Username", "", true)
					if err != nil {
						return err
					}
				}

				if email == "" {
					email, err = base.PromptString("Email", "e.g. user@example.com", true)
					if err != nil {
						return err
					}
				}

				if !cmd.Flags().Changed("fullname") {
					fullname, err = base.PromptString("Full name", "optional", false)
					if err != nil {
						return err
					}
				}

				if !cmd.Flags().Changed("active") {
					active, err = base.PromptBool("Active?", active)
					if err != nil {
						return err
					}
				}
			}

			if username == "" {
				return fmt.Errorf("--username is required")
			}
			if email == "" {
				return fmt.Errorf("--email is required")
			}

			// If password was not provided via flag, prompt interactively.
			// This masked prompt is independent of --interactive and always
			// triggers when --password is omitted.
			if password == "" {
				fmt.Fprint(cmd.OutOrStdout(), "Password: ")
				pw, err := term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					return fmt.Errorf("failed to read password: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout())

				fmt.Fprint(cmd.OutOrStdout(), "Confirm password: ")
				pw2, err := term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					return fmt.Errorf("failed to read password confirmation: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout())

				if string(pw) != string(pw2) {
					return fmt.Errorf("passwords do not match")
				}
				password = string(pw)
				if password == "" {
					return fmt.Errorf("password is required")
				}
			}

			if interactive {
				summary := fmt.Sprintf(
					"Create user:\n  Username: %s\n  Email:    %s\n  Fullname: %s\n  Active:   %v\n\nProceed? [y/N] ",
					username, email, fullname, active,
				)
				if !base.Confirm(cmd, summary) {
					return nil
				}
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			id, _, err := client.User.Create(cmd.Context(), poweradmin.UserCreateOpts{
				Username: username,
				Password: password,
				Email:    email,
				Fullname: fullname,
				Active:   active,
			})
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, map[string]any{
					"id":       id,
					"username": username,
					"email":    email,
				})
			}

			if base.IsQuiet(cmd) {
				fmt.Fprintln(cmd.OutOrStdout(), id)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "created user %s (id %d)\n", username, id)
			return nil
		},
	}

	cmd.Flags().String("username", "", "Username (required)")
	cmd.Flags().String("password", "", "Password (prompted if not provided)")
	cmd.Flags().String("email", "", "Email address (required)")
	cmd.Flags().String("fullname", "", "Full name")
	cmd.Flags().Bool("active", true, "Whether the user is active (default: true)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("quiet", "q", false, "Only print the ID of the created user")
	cmd.Flags().BoolP("interactive", "i", false, "Prompt interactively for missing values")
	return cmd
}
