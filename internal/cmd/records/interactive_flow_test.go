// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records_test

import (
	"context"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/records"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

// TestRecordsCreateInteractiveFullFlow drives every prompt in
// "records create --interactive" end to end via huh's accessible mode,
// rather than pre-filling required values through flags. This exercises
// the zone-name prompt, the per-field prompts, and the final confirmation
// together, which the flag-driven interactive tests in create_test.go
// deliberately skip.
func TestRecordsCreateInteractiveFullFlow(t *testing.T) {
	testutil.WithAccessiblePrompts(t)

	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	var createdName, createdType, createdContent string
	var createdTTL int
	mockRecord := &testutil.MockRecordClient{
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			createdName = opts.Name
			createdType = opts.Type
			createdContent = opts.Content
			createdTTL = opts.TTL
			return "rec-1", nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	// zone name, record name, type (1 = A, skips the MX/SRV-only priority
	// prompt), content, ttl, then the final "Proceed?" confirmation.
	testutil.WithDelayedStdin(t,
		"example.com\n",
		"www.example.com\n",
		"1\n",
		"1.2.3.4\n",
		"3600\n",
		"y\n",
	)

	err := fx.Run(records.NewCreateCmd(), []string{"--interactive"})
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

// TestRecordsDeleteInteractiveFullFlow drives "records delete --interactive"
// end to end: the zone single-select, the record multi-select, and the
// final confirmation. This is the only test exercising selectZoneName and
// the interactive dispatch branch of records delete; the other interactive
// delete tests (delete_internal_test.go) call deleteSelectedRecords
// directly and never touch the prompts themselves.
func TestRecordsDeleteInteractiveFullFlow(t *testing.T) {
	testutil.WithAccessiblePrompts(t)

	mockZone := &testutil.MockZoneClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Zone, error) {
			return []*poweradmin.Zone{
				{ID: 1, Name: "example.com"},
				{ID: 2, Name: "other.com"},
			}, nil
		},
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	var deletedID string
	mockRecord := &testutil.MockRecordClient{
		ListFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordListOpts) ([]*poweradmin.Record, *poweradmin.Response, error) {
			return []*poweradmin.Record{
				{ID: "rec-1", Name: "www.example.com", Type: "A", Content: "1.2.3.4"},
				{ID: "rec-2", Name: "mail.example.com", Type: "A", Content: "1.2.3.5"},
			}, nil, nil
		},
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			deletedID = recordID
			return nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, mockRecord)

	// zone select (1 = example.com), toggle record 1 (www.example.com),
	// 0 confirms the record selection, then the final "Proceed?" confirmation.
	testutil.WithDelayedStdin(t, "1\n", "1\n", "0\n", "y\n")

	err := fx.Run(records.NewDeleteCmd(nil), []string{"--interactive"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deletedID != "rec-1" {
		t.Errorf("expected rec-1 to be deleted, got %q", deletedID)
	}
	if !strings.Contains(fx.Stdout.String(), "deleted record www.example.com") {
		t.Errorf("expected success output, got:\n%s", fx.Stdout.String())
	}
}
