// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package version

import (
	"fmt"

	internalversion "github.com/contentways/poweradmin-cli/v2/internal/version"
	"github.com/spf13/cobra"
)

// NewVersionCmd returns a new "version" command instance.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of poweradmin",
		Long:  `Print the version and commit hash of the poweradmin CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "poweradmin version %s (commit %s)\n", internalversion.Version, internalversion.Commit)
		},
	}
}
