// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/internal/testutil"
)

func TestZonesImportDryRun(t *testing.T) {
	f, err := os.CreateTemp("", "zone-*.zone")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString(`$ORIGIN example.com.
$TTL 3600
@    IN    SOA    ns1.example.com. admin.example.com. (
                  2026060801 ; serial
                  3600       ; refresh
                  900        ; retry
                  604800     ; expire
                  300 )      ; minimum
@              IN    NS     ns1.example.com.
@              IN    A      1.2.3.4
www            IN    A      1.2.3.4
@              IN    MX     10 mail.example.com.
`)
	f.Close()

	fx := testutil.NewFixtureWithAllMocks(t, &testutil.MockZoneClient{}, &testutil.MockRecordClient{}, nil, nil, nil)
	err = fx.Run(zones.NewImportCmd(nil), []string{"--file", f.Name(), "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "Dry run") {
		t.Errorf("expected dry run output, got:\n%s", out)
	}
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected zone name, got:\n%s", out)
	}
	if !strings.Contains(out, "NS") {
		t.Errorf("expected NS record, got:\n%s", out)
	}
	if !strings.Contains(out, "MX") {
		t.Errorf("expected MX record, got:\n%s", out)
	}
}

func TestZonesImportDryRunDMARC(t *testing.T) {
	f, err := os.CreateTemp("", "zone-*.zone")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString(`$ORIGIN example.com.
$TTL 3600
@    IN    SOA    ns1.example.com. admin.example.com. (
                  2026060801 ; serial
                  3600       ; refresh
                  900        ; retry
                  604800     ; expire
                  300 )      ; minimum
_dmarc         IN    TXT    v=DMARC1; p=reject; rua=mailto:postmaster@example.com
`)
	f.Close()

	fx := testutil.NewFixtureWithAllMocks(t, &testutil.MockZoneClient{}, &testutil.MockRecordClient{}, nil, nil, nil)
	err = fx.Run(zones.NewImportCmd(nil), []string{"--file", f.Name(), "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "v=DMARC1; p=reject") {
		t.Errorf("expected full DMARC content, got:\n%s", out)
	}
}

func TestZonesImportMissingFile(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, &testutil.MockZoneClient{}, &testutil.MockRecordClient{}, nil, nil, nil)
	err := fx.Run(zones.NewImportCmd(nil), []string{})
	if err == nil {
		t.Fatal("expected error when no --file provided")
	}
}

func TestZonesImportFileNotFound(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, &testutil.MockZoneClient{}, &testutil.MockRecordClient{}, nil, nil, nil)
	err := fx.Run(zones.NewImportCmd(nil), []string{"--file", "/tmp/nonexistent.zone"})
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestZonesImportZoneNotFound(t *testing.T) {
	f, err := os.CreateTemp("", "zone-*.zone")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString(`$ORIGIN example.com.
$TTL 3600
@    IN    NS     ns1.example.com.
`)
	f.Close()

	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("not found")
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, &testutil.MockRecordClient{}, nil, nil, nil)
	err = fx.Run(zones.NewImportCmd(nil), []string{"--file", f.Name()})
	if err == nil {
		t.Fatal("expected error when zone not found and --create-zone not set")
	}
}

func TestZonesImportCreateZone(t *testing.T) {
	f, err := os.CreateTemp("", "zone-*.zone")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString(`$ORIGIN example.com.
$TTL 3600
@    IN    NS     ns1.example.com.
@    IN    A      1.2.3.4
`)
	f.Close()

	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("not found")
		},
		CreateFn: func(ctx context.Context, opts poweradmin.ZoneCreateOpts) (int, *poweradmin.Response, error) {
			return 42, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			return "rec-1", nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, mockRecord, nil, nil, nil)
	err = fx.Run(zones.NewImportCmd(nil), []string{"--file", f.Name(), "--create-zone"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "imported") {
		t.Errorf("expected import success message, got:\n%s", out)
	}
}
