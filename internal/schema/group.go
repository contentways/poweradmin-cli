// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package schema

import "contentways.dev/contentways/poweradmin-go/v2/poweradmin"

// Group is the CLI output schema for a Poweradmin group.
type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	PermTemplID int    `json:"perm_templ_id,omitempty"`
	MemberCount int    `json:"member_count,omitempty"`
	ZoneCount   int    `json:"zone_count,omitempty"`
}

// GroupFromSDK converts a poweradmin SDK Group to the CLI output schema.
func GroupFromSDK(g *poweradmin.Group) Group {
	return Group{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		PermTemplID: g.PermTemplID,
		MemberCount: g.MemberCount,
		ZoneCount:   g.ZoneCount,
	}
}

// GroupList wraps a slice of groups in a root object for JSON output.
type GroupList struct {
	Groups []Group `json:"groups"`
	Count  int     `json:"count"`
}

// GroupListFromSDK converts a slice of SDK Groups to the CLI output schema.
func GroupListFromSDK(groups []*poweradmin.Group) GroupList {
	out := make([]Group, len(groups))
	for i, g := range groups {
		out[i] = GroupFromSDK(g)
	}
	return GroupList{Groups: out, Count: len(out)}
}
