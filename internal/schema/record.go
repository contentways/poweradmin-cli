// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package schema

import (
	"strings"

	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

// Record is the CLI output schema for a DNS record.
type Record struct {
	ID      string `json:"id" yaml:"id"`
	Name    string `json:"name" yaml:"name"`
	Type    string `json:"type" yaml:"type"`
	Content string `json:"content" yaml:"content"`
	TTL     int    `json:"ttl" yaml:"ttl"`
	// Priority is set for record types that carry one (see HasPriority),
	// including a priority of 0, and omitted for all other types.
	Priority *int `json:"priority,omitempty" yaml:"priority,omitempty"`
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`
}

// HasPriority reports whether records of the given type carry a priority
// that PowerDNS stores separately from the content (MX and SRV).
func HasPriority(recordType string) bool {
	switch strings.ToUpper(recordType) {
	case "MX", "SRV":
		return true
	default:
		return false
	}
}

// RecordFromSDK converts a poweradmin SDK Record to the CLI output schema.
func RecordFromSDK(r *poweradmin.Record) Record {
	return Record{
		ID:       r.ID,
		Name:     r.Name,
		Type:     r.Type,
		Content:  r.Content,
		TTL:      r.TTL,
		Priority: priorityOf(r),
		Disabled: r.Disabled,
	}
}

func priorityOf(r *poweradmin.Record) *int {
	if !HasPriority(r.Type) {
		return nil
	}
	p := r.Priority
	return &p
}

// RecordList wraps a slice of records in a root object for JSON output.
type RecordList struct {
	Records []Record `json:"records" yaml:"records"`
	Count   int      `json:"count" yaml:"count"`
}

// RecordListFromSDK converts a slice of SDK Records to the CLI output schema.
func RecordListFromSDK(records []*poweradmin.Record) RecordList {
	out := make([]Record, len(records))
	for i, r := range records {
		out[i] = RecordFromSDK(r)
	}
	return RecordList{Records: out, Count: len(out)}
}
