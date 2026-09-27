// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"
	"strings"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/schema"
	"github.com/contentways/poweradmin-cli/v3/internal/state"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
	"github.com/spf13/cobra"
)

// NewDNSSECKeysCmd returns the "zones dnssec keys" command group.
func NewDNSSECKeysCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage the DNSSEC keys of a zone",
		Long:  `List, add, activate, deactivate and delete the DNSSEC keys of a zone (requires Poweradmin 4.5 or newer).`,
	}

	cmd.AddCommand(NewDNSSECKeysListCmd(s))
	cmd.AddCommand(NewDNSSECKeysAddCmd(s))
	cmd.AddCommand(NewDNSSECKeysActivateCmd(s))
	cmd.AddCommand(NewDNSSECKeysDeactivateCmd(s))
	cmd.AddCommand(NewDNSSECKeysDeleteCmd(s))

	return cmd
}

func activeLabel(active bool) string {
	if active {
		return "yes"
	}
	return "no"
}

// printKey writes a single key in the format chosen with --output.
func printKey(cmd *cobra.Command, message string, k *poweradmin.DNSSECKey) error {
	outputStr, _ := cmd.Flags().GetString("output")
	outputFmt := output.ParseFormat(outputStr)
	if outputFmt.IsStructured() {
		return base.PrintFormatted(cmd, outputFmt, schema.DNSSECKeyFromSDK(k))
	}
	if base.IsQuiet(cmd) {
		return nil
	}

	out := cmd.OutOrStdout()
	algorithm := k.Algorithm
	if algorithm == "" {
		algorithm = fmt.Sprintf("algorithm %d", k.AlgorithmID)
	}
	fmt.Fprintf(out, "%s: key %d (%s, keytag %d, %s, %d bits, active: %s)\n",
		message, k.ID, strings.ToUpper(string(k.Type)), k.KeyTag, algorithm, k.Bits, activeLabel(k.Active))
	for _, ds := range k.DS {
		fmt.Fprintf(out, "  DS %s\n", ds)
	}
	return nil
}

// keyIDFlag reads and validates --key-id.
func keyIDFlag(cmd *cobra.Command) (int, error) {
	keyID, _ := cmd.Flags().GetInt("key-id")
	if keyID <= 0 {
		return 0, fmt.Errorf("--key-id is required")
	}
	return keyID, nil
}

// NewDNSSECKeysListCmd returns a new "zones dnssec keys list" command instance.
func NewDNSSECKeysListCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the DNSSEC keys of a zone",
		Example: `  poweradmin zones dnssec keys list --name example.com
  poweradmin zones dnssec keys list --name example.com -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, zone, err := resolveDNSSECZone(cmd)
			if err != nil {
				return err
			}
			keys, _, err := client.DNSSEC.ListKeys(cmd.Context(), zone.ID)
			if err != nil {
				return fmt.Errorf("failed to list DNSSEC keys: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)
			if outputFmt.IsStructured() {
				return base.PrintFormatted(cmd, outputFmt, schema.DNSSECKeyListFromSDK(zone.Name, keys))
			}

			t := base.NewTable(cmd)
			t.AddHeader("ID", "TYPE", "KEYTAG", "ALGORITHM", "BITS", "ACTIVE")
			for _, k := range keys {
				algorithm := k.Algorithm
				if algorithm == "" {
					algorithm = fmt.Sprintf("%d", k.AlgorithmID)
				}
				t.AddRow(
					fmt.Sprintf("%d", k.ID),
					strings.ToUpper(string(k.Type)),
					fmt.Sprintf("%d", k.KeyTag),
					algorithm,
					fmt.Sprintf("%d", k.Bits),
					activeLabel(k.Active),
				)
			}
			t.Flush()
			return nil
		},
	}

	addZoneFlags(cmd, s)
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")
	return cmd
}

// NewDNSSECKeysAddCmd returns a new "zones dnssec keys add" command instance.
func NewDNSSECKeysAddCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a DNSSEC key to a zone",
		Long: `Add a DNSSEC key to a zone. Keys are created inactive unless --active is set.

Poweradmin validates the combination of algorithm and bits, e.g. ecdsa256
requires 256 bits and the RSA algorithms 1024 or 2048 bits.`,
		Example: `  poweradmin zones dnssec keys add --name example.com --type csk --algorithm ecdsa256 --bits 256 --active
  poweradmin zones dnssec keys add --name example.com --type zsk --algorithm rsasha256 --bits 2048`,
		RunE: func(cmd *cobra.Command, args []string) error {
			keyType, _ := cmd.Flags().GetString("type")
			algorithm, _ := cmd.Flags().GetString("algorithm")
			bits, _ := cmd.Flags().GetInt("bits")
			active, _ := cmd.Flags().GetBool("active")

			keyType = strings.ToLower(keyType)
			switch poweradmin.DNSSECKeyType(keyType) {
			case poweradmin.DNSSECKeyTypeKSK, poweradmin.DNSSECKeyTypeZSK, poweradmin.DNSSECKeyTypeCSK:
			default:
				return fmt.Errorf("--type must be one of ksk, zsk, csk")
			}
			if algorithm == "" {
				return fmt.Errorf("--algorithm is required")
			}
			if bits <= 0 {
				return fmt.Errorf("--bits is required")
			}

			client, zone, err := resolveDNSSECZone(cmd)
			if err != nil {
				return err
			}

			k, _, err := client.DNSSEC.AddKey(cmd.Context(), zone.ID, poweradmin.DNSSECKeyCreateOpts{
				Type:      poweradmin.DNSSECKeyType(keyType),
				Algorithm: strings.ToLower(algorithm),
				Bits:      bits,
				Active:    active,
			})
			if err != nil {
				return fmt.Errorf("failed to add DNSSEC key: %w", err)
			}
			return printKey(cmd, "added "+zone.Name, k)
		},
	}

	addZoneFlags(cmd, s)
	cmd.Flags().String("type", "", "Key type: ksk, zsk or csk (required)")
	cmd.Flags().String("algorithm", "", "Algorithm, e.g. ecdsa256, ed25519, rsasha256 (required)")
	cmd.Flags().Int("bits", 0, "Key size in bits, e.g. 256 for ecdsa256 (required)")
	cmd.Flags().Bool("active", false, "Create the key active")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output")
	cmd.RegisterFlagCompletionFunc("type", cobra.FixedCompletions([]string{"ksk", "zsk", "csk"}, cobra.ShellCompDirectiveNoFileComp))
	cmd.RegisterFlagCompletionFunc("algorithm", cobra.FixedCompletions(
		[]string{"ecdsa256", "ecdsa384", "ed25519", "ed448", "rsasha256", "rsasha512", "rsasha1", "rsasha1-nsec3-sha1"},
		cobra.ShellCompDirectiveNoFileComp,
	))
	return cmd
}

// newDNSSECKeySetActiveCmd returns the "zones dnssec keys activate" or
// "deactivate" command instance.
func newDNSSECKeySetActiveCmd(s *state.State, active bool) *cobra.Command {
	use, verb := "deactivate", "deactivated"
	if active {
		use, verb = "activate", "activated"
	}

	cmd := &cobra.Command{
		Use:     use,
		Short:   strings.ToUpper(use[:1]) + use[1:] + " a DNSSEC key",
		Example: fmt.Sprintf("  poweradmin zones dnssec keys %s --name example.com --key-id 3", use),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyID, err := keyIDFlag(cmd)
			if err != nil {
				return err
			}
			client, zone, err := resolveDNSSECZone(cmd)
			if err != nil {
				return err
			}
			k, _, err := client.DNSSEC.SetKeyActive(cmd.Context(), zone.ID, keyID, active)
			if err != nil {
				return fmt.Errorf("failed to %s DNSSEC key: %w", use, err)
			}
			return printKey(cmd, verb+" in "+zone.Name, k)
		},
	}

	addZoneFlags(cmd, s)
	cmd.Flags().Int("key-id", 0, "Key ID (required, see \"zones dnssec keys list\")")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|yaml")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output")
	return cmd
}

// NewDNSSECKeysActivateCmd returns a new "zones dnssec keys activate" command instance.
func NewDNSSECKeysActivateCmd(s *state.State) *cobra.Command {
	return newDNSSECKeySetActiveCmd(s, true)
}

// NewDNSSECKeysDeactivateCmd returns a new "zones dnssec keys deactivate" command instance.
func NewDNSSECKeysDeactivateCmd(s *state.State) *cobra.Command {
	return newDNSSECKeySetActiveCmd(s, false)
}

// NewDNSSECKeysDeleteCmd returns a new "zones dnssec keys delete" command instance.
func NewDNSSECKeysDeleteCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a DNSSEC key",
		Example: `  poweradmin zones dnssec keys delete --name example.com --key-id 3 --yes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			keyID, err := keyIDFlag(cmd)
			if err != nil {
				return err
			}
			client, zone, err := resolveDNSSECZone(cmd)
			if err != nil {
				return err
			}
			if !base.Confirm(cmd, fmt.Sprintf("Delete DNSSEC key %d of zone %s? [y/N] ", keyID, zone.Name)) {
				return nil
			}
			if _, err := client.DNSSEC.DeleteKey(cmd.Context(), zone.ID, keyID); err != nil {
				return fmt.Errorf("failed to delete DNSSEC key: %w", err)
			}
			if base.IsQuiet(cmd) {
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted DNSSEC key %d of zone %s\n", keyID, zone.Name)
			return nil
		},
	}

	addZoneFlags(cmd, s)
	cmd.Flags().Int("key-id", 0, "Key ID (required, see \"zones dnssec keys list\")")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output after deletion")
	return cmd
}
