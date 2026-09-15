// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package schema provides CLI-specific output schemas for JSON serialization.
// These types are separate from the SDK types to allow clean JSON output
// with snake_case keys, omitted empty fields and a consistent structure.
package schema

import "github.com/contentways/poweradmin-go/v3/poweradmin"

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

// ZoneMetadata is the CLI output schema for a single zone metadata kind.
type ZoneMetadata struct {
	Kind   string   `json:"kind"`
	Values []string `json:"values"`
}

// ZoneMetadataFromSDK converts a poweradmin SDK ZoneMetadata to the CLI output schema.
func ZoneMetadataFromSDK(m *poweradmin.ZoneMetadata) ZoneMetadata {
	return ZoneMetadata{Kind: m.Kind, Values: m.Values}
}

// ZoneMetadataList wraps a slice of zone metadata entries in a root object for JSON output.
type ZoneMetadataList struct {
	Metadata []ZoneMetadata `json:"metadata"`
	Count    int            `json:"count"`
}

// ZoneMetadataListFromSDK converts a slice of SDK ZoneMetadata to the CLI output schema.
func ZoneMetadataListFromSDK(metadata []*poweradmin.ZoneMetadata) ZoneMetadataList {
	out := make([]ZoneMetadata, len(metadata))
	for i, m := range metadata {
		out[i] = ZoneMetadataFromSDK(m)
	}
	return ZoneMetadataList{Metadata: out, Count: len(out)}
}
