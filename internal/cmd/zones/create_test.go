// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestZonesCreate(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		CreateFn: func(ctx context.Context, opts poweradmin.ZoneCreateOpts) (int, *poweradmin.Response, error) {
			return 42, nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewCreateCmd(), []string{"example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected output to contain example.com, got:\n%s", out)
	}
	if !strings.Contains(out, "42") {
		t.Errorf("expected output to contain id 42, got:\n%s", out)
	}
}

func TestZonesCreateJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		CreateFn: func(ctx context.Context, opts poweradmin.ZoneCreateOpts) (int, *poweradmin.Response, error) {
			return 42, nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.NewCreateCmd(), []string{"example.com", "--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"name": "example.com"`) {
		t.Errorf("expected JSON to contain example.com, got:\n%s", out)
	}
	if !strings.Contains(out, `"id": 42`) {
		t.Errorf("expected JSON to contain id 42, got:\n%s", out)
	}
}

func TestZonesCreateError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		CreateFn: func(ctx context.Context, opts poweradmin.ZoneCreateOpts) (int, *poweradmin.Response, error) {
			return 0, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)
	err := fx.Run(zones.NewCreateCmd(), []string{"example.com"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestZonesCreateInteractiveSkipsPromptsWhenAllFlagsSet verifies that when
// --interactive is combined with all required values already supplied via
// args/flags (including --nameserver and --ttl, so neither the nameserver
// nor the TTL prompt fires), no huh prompt runs at all — only the final
// confirmation, which is answered via a piped stdin.
func TestZonesCreateInteractiveSkipsPromptsWhenAllFlagsSet(t *testing.T) {
	var createdName string
	var createdType poweradmin.ZoneType
	var createdNS string

	mockZone := &testutil.MockZoneClient{
		CreateFn: func(ctx context.Context, opts poweradmin.ZoneCreateOpts) (int, *poweradmin.Response, error) {
			createdName = opts.Name
			createdType = opts.Type
			return 42, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			createdNS = opts.Content
			return "rec-1", nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)
	testutil.WithStdin(t, "y\n")

	err := fx.Run(zones.NewCreateCmd(), []string{
		"example.com",
		"--type", "NATIVE",
		"--nameserver", "ns1.example.com",
		"--ttl", "3600",
		"--interactive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if createdName != "example.com" {
		t.Errorf("expected zone name example.com, got %q", createdName)
	}
	if createdType != poweradmin.ZoneType("NATIVE") {
		t.Errorf("expected zone type NATIVE, got %q", createdType)
	}
	if createdNS != "ns1.example.com" {
		t.Errorf("expected NS record content ns1.example.com, got %q", createdNS)
	}
}

// TestZonesCreateInteractiveDeclineAbortsWithoutCreating verifies that
// declining the final confirmation prevents both the zone and NS record API
// calls entirely. --nameserver and --ttl are set explicitly so the test
// only exercises base.Confirm, never a huh prompt.
func TestZonesCreateInteractiveDeclineAbortsWithoutCreating(t *testing.T) {
	zoneCreated := false
	recordCreated := false

	mockZone := &testutil.MockZoneClient{
		CreateFn: func(ctx context.Context, opts poweradmin.ZoneCreateOpts) (int, *poweradmin.Response, error) {
			zoneCreated = true
			return 42, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			recordCreated = true
			return "rec-1", nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)
	testutil.WithStdin(t, "n\n")

	err := fx.Run(zones.NewCreateCmd(), []string{
		"example.com",
		"--type", "NATIVE",
		"--nameserver", "ns1.example.com",
		"--ttl", "3600",
		"--interactive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if zoneCreated {
		t.Error("expected zone NOT to be created after declining confirmation")
	}
	if recordCreated {
		t.Error("expected NS record NOT to be created after declining confirmation")
	}
}
