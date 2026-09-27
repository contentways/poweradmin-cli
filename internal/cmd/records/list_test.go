// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/records"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestRecordsList(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4", TTL: 3600},
				{ID: "rec-2", Name: "mail.example.com", Type: "MX", Content: "mail.example.com", TTL: 3600},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewListCmd(nil), []string{"--zone-name", "example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "www.example.com") {
		t.Errorf("expected output to contain www.example.com, got:\n%s", out)
	}
	if !strings.Contains(out, "1.2.3.4") {
		t.Errorf("expected output to contain 1.2.3.4, got:\n%s", out)
	}
}

func TestRecordsListJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4", TTL: 3600},
			}, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewListCmd(nil), []string{"--zone-name", "example.com", "--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"name": "www.example.com"`) {
		t.Errorf("expected JSON to contain www.example.com, got:\n%s", out)
	}
}

func TestRecordsListMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, &testutil.MockRecordClient{})
	err := fx.Run(records.NewListCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestRecordsListError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)
	err := fx.Run(records.NewListCmd(nil), []string{"--zone-name", "example.com"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRecordsListByZoneID(t *testing.T) {
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			if zoneID != 9 {
				t.Errorf("expected zoneID 9, got %d", zoneID)
			}
			return []*poweradmin.Record{
				{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4"},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, mockRecord)

	err := fx.Run(records.NewListCmd(nil), []string{"--zone-id", "9"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestRecordsListFilterByType assumes "records list" filters client-side by
// --type, mirroring "zones list --type". If --type isn't implemented this
// way, this test will fail at runtime (unknown flag) rather than at compile
// time — check the failure message if so.
func TestRecordsListFilterByType(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4"},
				{ID: "rec-2", Name: "example.com", Type: "MX", Content: "mail.example.com"},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewListCmd(nil), []string{"--zone-name", "example.com", "--type", "MX"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if strings.Contains(out, "www.example.com") {
		t.Errorf("expected A record to be filtered out, got:\n%s", out)
	}
	if !strings.Contains(out, "mail.example.com") {
		t.Errorf("expected MX record content present, got:\n%s", out)
	}
}

func priorityFixture(t *testing.T) *testutil.Fixture {
	t.Helper()
	longTXT := "v=DKIM1; k=rsa; p=" + strings.Repeat("A", 300)
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			return []*poweradmin.Record{
				{ID: "r1", Name: "example.com", Type: "MX", Content: "mail.example.com", TTL: 60, Priority: 10},
				{ID: "r2", Name: "_imaps._tcp.example.com", Type: "SRV", Content: "1 993 mail.example.com", TTL: 60, Priority: 0},
				{ID: "r3", Name: "www.example.com", Type: "A", Content: "192.0.2.1", TTL: 300},
				{ID: "r4", Name: "dkim._domainkey.example.com", Type: "TXT", Content: longTXT, TTL: 60},
			}, nil
		},
	}
	return testutil.NewFixtureWithMocks(t, mockZone, mockRecord)
}

// CONTENT is the last column, so a long TXT record in full mode must not
// push TTL/PRIO of the other rows to the right.
func TestRecordsListColumnsAndPriority(t *testing.T) {
	fx := priorityFixture(t)
	if err := fx.Run(records.NewListCmd(nil), []string{"--zone-name", "example.com", "-o", "full"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimRight(fx.Stdout.String(), "\n"), "\n")
	if got := strings.Fields(lines[0]); strings.Join(got, " ") != "NAME TYPE TTL PRIO CONTENT" {
		t.Fatalf("header = %q", lines[0])
	}
	rows := map[string][]string{}
	for _, l := range lines[1:] {
		f := strings.Fields(l)
		rows[f[1]] = f
	}
	if f := rows["MX"]; f[2] != "60" || f[3] != "10" || f[4] != "mail.example.com" {
		t.Errorf("MX row = %v, want TTL 60, PRIO 10", f)
	}
	if f := rows["SRV"]; f[2] != "60" || f[3] != "0" || f[4] != "1" {
		t.Errorf("SRV row = %v, want TTL 60, PRIO 0", f)
	}
	if f := rows["A"]; len(f) != 4 || f[3] != "192.0.2.1" {
		t.Errorf("A row = %v, want empty PRIO", f)
	}
	for _, l := range lines {
		if strings.Contains(l, "mail.example.com") && len(l) > 120 {
			t.Errorf("row padded by the long TXT record (%d chars): %q", len(l), l)
		}
	}
}

func TestRecordsListJSONPriority(t *testing.T) {
	fx := priorityFixture(t)
	if err := fx.Run(records.NewListCmd(nil), []string{"--zone-name", "example.com", "-o", "json"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out struct {
		Records []map[string]any `json:"records"`
	}
	if err := json.Unmarshal(fx.Stdout.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	want := map[string]any{"MX": float64(10), "SRV": float64(0)}
	for _, r := range out.Records {
		p, ok := r["priority"]
		typ := r["type"].(string)
		if w, needs := want[typ]; needs {
			if !ok || p != w {
				t.Errorf("%s priority = %v (present %v), want %v", typ, p, ok, w)
			}
		} else if ok {
			t.Errorf("%s must not carry a priority, got %v", typ, p)
		}
	}
}

func TestRecordsListSortByPrio(t *testing.T) {
	fx := priorityFixture(t)
	if err := fx.Run(records.NewListCmd(nil), []string{"--zone-name", "example.com", "--type", "MX", "--sort", "prio"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "mail.example.com") {
		t.Errorf("unexpected output:\n%s", fx.Stdout.String())
	}
}
