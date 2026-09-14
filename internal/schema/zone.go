// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package schema provides CLI-specific output schemas for JSON serialization.
// These types are separate from the SDK types to allow clean JSON output
// with snake_case keys, omitted empty fields and a consistent structure.
package schema

import "contentways.dev/contentways/poweradmin-go/v2/poweradmin"

// Zone is the CLI output schema for a DNS zone.
type Zone struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Masters      string `json:"masters,omitempty"`
	Description  string `json:"description,omitempty"`
	SOASerial    int    `json:"soa_serial,omitempty"`
	DNSSECSigned bool   `json:"dnssec_signed,omitempty"`
}

// ZoneFromSDK converts a poweradmin SDK Zone to the CLI output schema.
func ZoneFromSDK(z *poweradmin.Zone) Zone {
	return Zone{
		ID:           z.ID,
		Name:         z.Name,
		Type:         string(z.Type),
		Masters:      z.Masters,
		Description:  z.Description,
		SOASerial:    z.SOASerial,
		DNSSECSigned: z.DNSSECSigned,
	}
}

// ZoneList wraps a slice of zones in a root object for JSON output.
type ZoneList struct {
	Zones []Zone `json:"zones"`
	Count int    `json:"count"`
}

// ZoneListFromSDK converts a slice of SDK Zones to the CLI output schema.
func ZoneListFromSDK(zones []*poweradmin.Zone) ZoneList {
	out := make([]Zone, len(zones))
	for i, z := range zones {
		out[i] = ZoneFromSDK(z)
	}
	return ZoneList{Zones: out, Count: len(out)}
}

type ZoneWithNameservers struct {
	Zone
	Nameservers []string `json:"nameservers,omitempty"`
}
