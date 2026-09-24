// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
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

func TestZonesImportExistingZone(t *testing.T) {
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

	var createdCount int
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 7, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			if zoneID != 7 {
				t.Errorf("expected zoneID 7, got %d", zoneID)
			}
			createdCount++
			return "rec-1", nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, mockRecord, nil, nil, nil)
	err = fx.Run(zones.NewImportCmd(nil), []string{"--file", f.Name()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if createdCount != 2 {
		t.Errorf("expected 2 records created, got %d", createdCount)
	}
	if !strings.Contains(fx.Stdout.String(), "imported 2 records") {
		t.Errorf("expected import success message, got:\n%s", fx.Stdout.String())
	}
}

func TestZonesImportZoneNilWithoutError(t *testing.T) {
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
			return nil, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, mockZone, &testutil.MockRecordClient{}, nil, nil, nil)
	err = fx.Run(zones.NewImportCmd(nil), []string{"--file", f.Name()})
	if err == nil {
		t.Fatal("expected error for nil zone with nil error, got nil")
	}
}

func TestZonesImportRecordCreateFailureIsSkipped(t *testing.T) {
	f, err := os.CreateTemp("", "zone-*.zone")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString(`$ORIGIN example.com.
$TTL 3600
@      IN    NS     ns1.example.com.
www    IN    A      1.2.3.4
mail   IN    A      1.2.3.5
`)
	f.Close()

	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	callCount := 0
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			callCount++
			if opts.Name == "www.example.com" {
				return "", nil, fmt.Errorf("api error")
			}
			return "rec-ok", nil, nil
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, mockRecord, nil, nil, nil)
	err = fx.Run(zones.NewImportCmd(nil), []string{"--file", f.Name()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 3 {
		t.Errorf("expected all 3 records attempted despite one failing, got %d calls", callCount)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "imported 2 records") || !strings.Contains(out, "skipped: 1") {
		t.Errorf("expected 2 imported, 1 skipped, got:\n%s", out)
	}
	if !strings.Contains(fx.Stderr.String(), "failed to create record") {
		t.Errorf("expected warning on stderr, got:\n%s", fx.Stderr.String())
	}
}

func TestZonesImportCreateZoneFailure(t *testing.T) {
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
		CreateFn: func(ctx context.Context, opts poweradmin.ZoneCreateOpts) (int, *poweradmin.Response, error) {
			return 0, nil, fmt.Errorf("api error")
		},
	}

	fx := testutil.NewFixtureWithAllMocks(t, mockZone, &testutil.MockRecordClient{}, nil, nil, nil)
	err = fx.Run(zones.NewImportCmd(nil), []string{"--file", f.Name(), "--create-zone"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestZonesImportExplicitZoneNameOverridesOrigin(t *testing.T) {
	f, err := os.CreateTemp("", "zone-*.zone")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString(`$ORIGIN origin.example.
$TTL 3600
@    IN    NS     ns1.example.com.
`)
	f.Close()

	var requestedName string
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			requestedName = name
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, mockZone, &testutil.MockRecordClient{}, nil, nil, nil)
	err = fx.Run(zones.NewImportCmd(nil), []string{"--file", f.Name(), "--zone-name", "override.example"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if requestedName != "override.example" {
		t.Errorf("expected --zone-name to override $ORIGIN, got %q", requestedName)
	}
}

func TestZonesImportInvalidRecordLineSkipped(t *testing.T) {
	f, err := os.CreateTemp("", "zone-*.zone")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	// "bad" has too few fields to be a valid record line and is silently
	// skipped by parseRecordLine; the valid line after it still imports.
	f.WriteString(`$ORIGIN example.com.
$TTL 3600
bad
www    IN    A      1.2.3.4
`)
	f.Close()

	fx := testutil.NewFixtureWithAllMocks(t, &testutil.MockZoneClient{}, &testutil.MockRecordClient{}, nil, nil, nil)
	err = fx.Run(zones.NewImportCmd(nil), []string{"--file", f.Name(), "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "records to import: 1") {
		t.Errorf("expected only the valid record counted, got:\n%s", out)
	}
}
