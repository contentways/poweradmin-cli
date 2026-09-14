// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// zoneRecord holds a parsed DNS record from a BIND zone file.
type zoneRecord struct {
	Name     string
	TTL      int
	Type     string
	Content  string
	Priority int
}

// NewImportCmd returns a new "zones import" command instance.
// Imports DNS records from a BIND-compatible zone file into Poweradmin.
// SOA records are skipped — they are managed by Poweradmin internally.
// The zone must already exist unless --create-zone is set.
func NewImportCmd(s *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import a DNS zone from a BIND zone file",
		Long:  `Import DNS records from a BIND-compatible zone file into Poweradmin.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			file, _ := cmd.Flags().GetString("file")
			zoneName, _ := cmd.Flags().GetString("zone-name")
			createZone, _ := cmd.Flags().GetBool("create-zone")
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			if file == "" {
				return fmt.Errorf("--file is required")
			}

			// Open and parse the zone file.
			f, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer f.Close()

			origin, defaultTTL, records, err := parseZoneFile(f, zoneName)
			if err != nil {
				return fmt.Errorf("failed to parse zone file: %w", err)
			}

			// Use $ORIGIN from file if --zone-name not provided.
			if zoneName == "" {
				zoneName = origin
			}
			if zoneName == "" {
				return fmt.Errorf("could not determine zone name — set $ORIGIN in file or use --zone-name")
			}

			_ = defaultTTL

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run — zone: %s, records to import: %d\n", zoneName, len(records))
				for _, r := range records {
					if r.Priority > 0 {
						fmt.Fprintf(cmd.OutOrStdout(), "  %-30s %-6d %-6s %-6d %s\n", r.Name, r.TTL, r.Type, r.Priority, r.Content)
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), "  %-30s %-6d %-6s %s\n", r.Name, r.TTL, r.Type, r.Content)
					}
				}
				return nil
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			// Resolve or create the zone.
			var zoneID int
			zone, _, err := client.Zone.GetByName(cmd.Context(), zoneName)
			if err != nil {
				if !createZone {
					return fmt.Errorf("zone %s not found — use --create-zone to create it", zoneName)
				}
				// Create the zone.
				zoneID, _, err = client.Zone.Create(cmd.Context(), poweradmin.ZoneCreateOpts{
					Name: zoneName,
					Type: "NATIVE",
				})
				if err != nil {
					return fmt.Errorf("failed to create zone: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "created zone %s (id %d)\n", zoneName, zoneID)
			} else {
				zoneID = zone.ID
			}

			// Import records.
			created := 0
			skipped := 0
			for _, r := range records {
				_, _, err := client.Record.Create(cmd.Context(), zoneID, poweradmin.RecordCreateOpts{
					Name:     r.Name,
					Type:     r.Type,
					Content:  r.Content,
					TTL:      r.TTL,
					Priority: r.Priority,
				})
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: failed to create record %s %s: %v\n", r.Name, r.Type, err)
					skipped++
					continue
				}
				created++
			}

			fmt.Fprintf(cmd.OutOrStdout(), "imported %d records into zone %s (skipped: %d)\n", created, zoneName, skipped)
			return nil
		},
	}

	cmd.Flags().String("file", "", "Path to BIND zone file (required)")
	cmd.Flags().String("zone-name", "", "Target zone name (overrides $ORIGIN in file)")
	cmd.Flags().Bool("create-zone", false, "Create zone if it does not exist")
	cmd.Flags().Bool("dry-run", false, "Show what would be imported without making changes")
	cmd.RegisterFlagCompletionFunc("zone-name", base.ZoneNameCompletion(s))
	return cmd
}

// parseZoneFile parses a BIND zone file and returns the origin, default TTL
// and a list of DNS records. SOA records are skipped.
func parseZoneFile(f *os.File, zoneName string) (origin string, defaultTTL int, records []zoneRecord, err error) {
	defaultTTL = 3600
	scanner := bufio.NewScanner(f)

	// Support long lines (e.g. DKIM TXT records).
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and full-line comments.
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		// Strip inline comments only from directive lines ($ORIGIN, $TTL).
		// Record lines may contain semicolons in TXT content (DMARC, SPF, DKIM)
		// so we never strip semicolons from record lines.
		if strings.HasPrefix(line, "$") {
			if idx := strings.Index(line, ";"); idx != -1 {
				line = strings.TrimSpace(line[:idx])
			}
		}

		// Handle $ORIGIN directive.
		if strings.HasPrefix(line, "$ORIGIN") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				origin = strings.TrimSuffix(parts[1], ".")
			}
			continue
		}

		// Handle $TTL directive.
		if strings.HasPrefix(line, "$TTL") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				defaultTTL, _ = strconv.Atoi(parts[1])
			}
			continue
		}

		// Skip SOA records — managed by Poweradmin internally.
		if strings.Contains(line, "SOA") {
			// Skip multiline SOA (lines ending with '(').
			if strings.Contains(line, "(") {
				for scanner.Scan() {
					soaLine := strings.TrimSpace(scanner.Text())
					if strings.Contains(soaLine, ")") {
						break
					}
				}
			}
			continue
		}

		// Parse record line: name [ttl] IN type content
		r, parseErr := parseRecordLine(line, origin, zoneName, defaultTTL)
		if parseErr != nil {
			continue
		}
		records = append(records, r)
	}

	if scanErr := scanner.Err(); scanErr != nil {
		err = fmt.Errorf("error reading file: %w", scanErr)
	}
	return
}

// parseRecordLine parses a single BIND record line into a zoneRecord.
// Handles relative and absolute names, converts @ to zone apex.
func parseRecordLine(line, origin, zoneName string, defaultTTL int) (zoneRecord, error) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return zoneRecord{}, fmt.Errorf("invalid record line: %s", line)
	}

	zone := origin
	if zoneName != "" {
		zone = zoneName
	}

	var r zoneRecord
	idx := 0

	// Parse name.
	name := fields[idx]
	idx++

	// Convert @ to zone apex.
	if name == "@" {
		name = zone
	} else if strings.HasSuffix(name, "."+zone) || name == zone {
		// Already absolute — strip trailing dot.
		name = strings.TrimSuffix(name, ".")
	} else if !strings.Contains(name, "."+zone) {
		// Relative name — make absolute.
		name = name + "." + zone
	} else {
		name = strings.TrimSuffix(name, ".")
	}
	r.Name = name

	// Parse optional TTL.
	r.TTL = defaultTTL
	if ttl, err := strconv.Atoi(fields[idx]); err == nil {
		r.TTL = ttl
		idx++
	}

	// Skip IN class.
	if strings.ToUpper(fields[idx]) == "IN" {
		idx++
	}

	if idx >= len(fields) {
		return zoneRecord{}, fmt.Errorf("invalid record line: %s", line)
	}

	// Parse type.
	r.Type = strings.ToUpper(fields[idx])
	idx++

	if idx >= len(fields) {
		return zoneRecord{}, fmt.Errorf("invalid record line: %s", line)
	}

	// Parse content — rest of the line.
	content := strings.Join(fields[idx:], " ")

	// Handle MX priority — content starts with a number.
	if r.Type == "MX" {
		parts := strings.SplitN(content, " ", 2)
		if len(parts) == 2 {
			r.Priority, _ = strconv.Atoi(parts[0])
			content = parts[1]
		}
	}

	r.Content = content

	// Strip trailing dot from content for record types that contain hostnames.
	switch r.Type {
	case "NS", "CNAME", "MX", "PTR", "SRV":
		r.Content = strings.TrimSuffix(r.Content, ".")
	}

	return r, nil
}
