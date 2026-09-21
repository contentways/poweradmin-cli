// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/records"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestRecordsCreate(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			return "rec-42", nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewCreateCmd(), []string{
		"--zone-name", "example.com",
		"--name", "www.example.com",
		"--type", "A",
		"--content", "1.2.3.4",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "www.example.com") {
		t.Errorf("expected output to contain www.example.com, got:\n%s", out)
	}
	if !strings.Contains(out, "rec-42") {
		t.Errorf("expected output to contain rec-42, got:\n%s", out)
	}
}

func TestRecordsCreateJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			return "rec-42", nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewCreateCmd(), []string{
		"--zone-name", "example.com",
		"--name", "www.example.com",
		"--type", "A",
		"--content", "1.2.3.4",
		"--output", "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"id": "rec-42"`) {
		t.Errorf("expected JSON to contain rec-42, got:\n%s", out)
	}
	if !strings.Contains(out, `"name": "www.example.com"`) {
		t.Errorf("expected JSON to contain www.example.com, got:\n%s", out)
	}
}

func TestRecordsCreateMissingFlags(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, &testutil.MockRecordClient{})
	err := fx.Run(records.NewCreateCmd(), []string{})
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestRecordsCreateError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			return "", nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)
	err := fx.Run(records.NewCreateCmd(), []string{
		"--zone-name", "example.com",
		"--name", "www.example.com",
		"--type", "A",
		"--content", "1.2.3.4",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestRecordsCreateInteractiveSkipsPromptsWhenAllFlagsSet verifies that
// when --interactive is combined with all required values already
// supplied via flags (--zone-name, --name, --type, --content, --ttl all
// set, and --type A so the priority prompt doesn't fire since it's only
// shown for MX/SRV), no huh prompt runs at all — only the final
// confirmation, answered via a piped stdin.
func TestRecordsCreateInteractiveSkipsPromptsWhenAllFlagsSet(t *testing.T) {
	var createdName, createdType, createdContent string
	var createdTTL int

	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			createdName = opts.Name
			createdType = opts.Type
			createdContent = opts.Content
			createdTTL = opts.TTL
			return "rec-42", nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)
	testutil.WithStdin(t, "y\n")

	err := fx.Run(records.NewCreateCmd(), []string{
		"--zone-name", "example.com",
		"--name", "www.example.com",
		"--type", "A",
		"--content", "1.2.3.4",
		"--ttl", "3600",
		"--interactive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if createdName != "www.example.com" {
		t.Errorf("expected name www.example.com, got %q", createdName)
	}
	if createdType != "A" {
		t.Errorf("expected type A, got %q", createdType)
	}
	if createdContent != "1.2.3.4" {
		t.Errorf("expected content 1.2.3.4, got %q", createdContent)
	}
	if createdTTL != 3600 {
		t.Errorf("expected ttl 3600, got %d", createdTTL)
	}
}

// TestRecordsCreateInteractiveDeclineAbortsWithoutCreating verifies that
// declining the final confirmation prevents the API call entirely.
func TestRecordsCreateInteractiveDeclineAbortsWithoutCreating(t *testing.T) {
	created := false
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name, Type: "NATIVE"}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			created = true
			return "rec-42", nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)
	testutil.WithStdin(t, "n\n")

	err := fx.Run(records.NewCreateCmd(), []string{
		"--zone-name", "example.com",
		"--name", "www.example.com",
		"--type", "A",
		"--content", "1.2.3.4",
		"--ttl", "3600",
		"--interactive",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created {
		t.Error("expected record NOT to be created after declining confirmation")
	}
}

func TestRecordsCreateByZoneID(t *testing.T) {
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			if zoneID != 7 {
				t.Errorf("expected zoneID 7, got %d", zoneID)
			}
			return "rec-99", nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, mockRecord)

	err := fx.Run(records.NewCreateCmd(), []string{
		"--zone-id", "7",
		"--name", "www.example.com",
		"--type", "A",
		"--content", "1.2.3.4",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "rec-99") {
		t.Errorf("expected output to contain rec-99, got:\n%s", fx.Stdout.String())
	}
}

func TestRecordsCreateQuiet(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			return "rec-42", nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewCreateCmd(), []string{
		"--zone-name", "example.com",
		"--name", "www.example.com",
		"--type", "A",
		"--content", "1.2.3.4",
		"--quiet",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(fx.Stdout.String())
	if out != "rec-42" {
		t.Errorf("expected quiet output to be just the id, got: %q", out)
	}
}

func TestRecordsCreateMXWithPriority(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	var createdPriority int
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			createdPriority = opts.Priority
			return "rec-mx", nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	err := fx.Run(records.NewCreateCmd(), []string{
		"--zone-name", "example.com",
		"--name", "example.com",
		"--type", "MX",
		"--content", "mail.example.com",
		"--priority", "10",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if createdPriority != 10 {
		t.Errorf("expected priority 10, got %d", createdPriority)
	}
}
