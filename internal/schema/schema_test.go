// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package schema_test

import (
	"testing"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/schema"
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
