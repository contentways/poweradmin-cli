// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package schema_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/schema"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestZoneFromSDK(t *testing.T) {
	z := &poweradmin.Zone{
		ID:   1,
		Name: "example.com",
		Type: "NATIVE",
	}
	out := schema.ZoneFromSDK(z)
	if out.ID != 1 || out.Name != "example.com" || out.Type != "NATIVE" {
		t.Errorf("ZoneFromSDK = %+v", out)
	}
}

func TestZoneListFromSDK(t *testing.T) {
	zones := []*poweradmin.Zone{
		{ID: 1, Name: "example.com", Type: "NATIVE"},
		{ID: 2, Name: "example.org", Type: "NATIVE"},
	}
	out := schema.ZoneListFromSDK(zones)
	if out.Count != 2 || len(out.Zones) != 2 {
		t.Errorf("ZoneListFromSDK = %+v", out)
	}
}

func TestRecordFromSDK(t *testing.T) {
	r := &poweradmin.Record{
		ID:      "rec-1",
		Name:    "www.example.com",
		Type:    "A",
		Content: "1.2.3.4",
		TTL:     3600,
	}
	out := schema.RecordFromSDK(r)
	if out.ID != "rec-1" || out.Name != "www.example.com" || out.Content != "1.2.3.4" {
		t.Errorf("RecordFromSDK = %+v", out)
	}
}

func TestRecordListFromSDK(t *testing.T) {
	records := []*poweradmin.Record{
		{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4", TTL: 3600},
	}
	out := schema.RecordListFromSDK(records)
	if out.Count != 1 || len(out.Records) != 1 {
		t.Errorf("RecordListFromSDK = %+v", out)
	}
}

func TestUserFromSDK(t *testing.T) {
	u := &poweradmin.User{
		ID:       1,
		Username: "patrick",
		Email:    "patrick@example.com",
		Active:   true,
	}
	out := schema.UserFromSDK(u)
	if out.ID != 1 || out.Username != "patrick" || out.Email != "patrick@example.com" {
		t.Errorf("UserFromSDK = %+v", out)
	}
}

func TestUserListFromSDK(t *testing.T) {
	users := []*poweradmin.User{
		{ID: 1, Username: "patrick", Email: "patrick@example.com", Active: true},
		{ID: 2, Username: "markus", Email: "markus@example.com", Active: true},
	}
	out := schema.UserListFromSDK(users)
	if out.Count != 2 || len(out.Users) != 2 {
		t.Errorf("UserListFromSDK = %+v", out)
	}
}

func TestZoneWithNameservers(t *testing.T) {
	z := &poweradmin.Zone{
		ID:   1,
		Name: "example.com",
		Type: "NATIVE",
	}
	out := schema.ZoneWithNameservers{
		Zone:        schema.ZoneFromSDK(z),
		Nameservers: []string{"ns1.example.com", "ns2.example.com"},
	}
	if out.Name != "example.com" {
		t.Errorf("ZoneWithNameservers.Name = %q", out.Name)
	}
	if len(out.Nameservers) != 2 {
		t.Errorf("ZoneWithNameservers.Nameservers = %v", out.Nameservers)
	}
}

func TestRecordFromSDKPriority(t *testing.T) {
	cases := []struct {
		typ  string
		prio int
		want *int
	}{
		{"MX", 10, new(10)},
		{"SRV", 0, new(0)},
		{"mx", 5, new(5)},
		{"A", 0, nil},
		{"TXT", 7, nil},
	}
	for _, c := range cases {
		out := schema.RecordFromSDK(&poweradmin.Record{Type: c.typ, Priority: c.prio})
		switch {
		case c.want == nil && out.Priority != nil:
			t.Errorf("%s: Priority = %d, want nil", c.typ, *out.Priority)
		case c.want != nil && (out.Priority == nil || *out.Priority != *c.want):
			t.Errorf("%s: Priority = %v, want %d", c.typ, out.Priority, *c.want)
		}
	}
}
