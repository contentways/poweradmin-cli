// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package schema

import "github.com/contentways/poweradmin-go/v3/poweradmin"

// Record is the CLI output schema for a DNS record.
type Record struct {
	ID       string `json:"id" yaml:"id"`
	Name     string `json:"name" yaml:"name"`
	Type     string `json:"type" yaml:"type"`
	Content  string `json:"content" yaml:"content"`
	TTL      int    `json:"ttl" yaml:"ttl"`
	Priority int    `json:"priority,omitempty" yaml:"priority,omitempty"`
	Disabled bool   `json:"disabled,omitempty" yaml:"disabled,omitempty"`
}

// RecordFromSDK converts a poweradmin SDK Record to the CLI output schema.
func RecordFromSDK(r *poweradmin.Record) Record {
	return Record{
		ID:       r.ID,
		Name:     r.Name,
		Type:     r.Type,
		Content:  r.Content,
		TTL:      r.TTL,
		Priority: r.Priority,
		Disabled: r.Disabled,
	}
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
