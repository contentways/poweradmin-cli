// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil_test

import (
	"context"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v3/poweradmin"
)

func TestMockRecordClientUnsetFieldsReturnZeroValues(t *testing.T) {
	m := &testutil.MockRecordClient{}
	ctx := context.Background()

	if r, resp, err := m.GetByID(ctx, 1, "rec-1"); r != nil || resp != nil || err != nil {
		t.Errorf("GetByID: expected all nil, got %v, %v, %v", r, resp, err)
	}
	if rs, resp, err := m.List(ctx, 1, poweradmin.RecordListOpts{}); rs != nil || resp != nil || err != nil {
		t.Errorf("List: expected all nil, got %v, %v, %v", rs, resp, err)
	}
	if rs, err := m.All(ctx, 1); rs != nil || err != nil {
		t.Errorf("All: expected all nil, got %v, %v", rs, err)
	}
	if id, resp, err := m.Create(ctx, 1, poweradmin.RecordCreateOpts{}); id != "" || resp != nil || err != nil {
		t.Errorf("Create: expected zero values, got %v, %v, %v", id, resp, err)
	}
	if r, resp, err := m.Update(ctx, 1, "rec-1", poweradmin.RecordUpdateOpts{}); r != nil || resp != nil || err != nil {
		t.Errorf("Update: expected all nil, got %v, %v, %v", r, resp, err)
	}
	if resp, err := m.Delete(ctx, 1, "rec-1"); resp != nil || err != nil {
		t.Errorf("Delete: expected all nil, got %v, %v", resp, err)
	}
	if res, resp, err := m.Bulk(ctx, 1, nil); res != nil || resp != nil || err != nil {
		t.Errorf("Bulk: expected all nil, got %v, %v, %v", res, resp, err)
	}
}

func TestMockRecordClientSetFieldsDelegate(t *testing.T) {
	var calls []string
	m := &testutil.MockRecordClient{
		GetByIDFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Record, *poweradmin.Response, error) {
			calls = append(calls, "GetByID")
			return &poweradmin.Record{ID: recordID}, nil, nil
		},
		ListFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordListOpts) ([]*poweradmin.Record, *poweradmin.Response, error) {
			calls = append(calls, "List")
			return []*poweradmin.Record{{ID: "rec-1"}}, nil, nil
		},
		AllFn: func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
			calls = append(calls, "All")
			return []*poweradmin.Record{{ID: "rec-1"}}, nil
		},
		CreateFn: func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
			calls = append(calls, "Create")
			return "rec-42", nil, nil
		},
		UpdateFn: func(ctx context.Context, zoneID int, recordID string, opts poweradmin.RecordUpdateOpts) (*poweradmin.Record, *poweradmin.Response, error) {
			calls = append(calls, "Update")
			return &poweradmin.Record{ID: recordID}, nil, nil
		},
		DeleteFn: func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
			calls = append(calls, "Delete")
			return nil, nil
		},
		BulkFn: func(ctx context.Context, zoneID int, ops []poweradmin.BulkRecordOperation) (*poweradmin.BulkRecordsResult, *poweradmin.Response, error) {
			calls = append(calls, "Bulk")
			return &poweradmin.BulkRecordsResult{Created: 1}, nil, nil
		},
	}
	ctx := context.Background()

	_, _, _ = m.GetByID(ctx, 1, "rec-1")
	_, _, _ = m.List(ctx, 1, poweradmin.RecordListOpts{})
	_, _ = m.All(ctx, 1)
	_, _, _ = m.Create(ctx, 1, poweradmin.RecordCreateOpts{})
	_, _, _ = m.Update(ctx, 1, "rec-1", poweradmin.RecordUpdateOpts{})
	_, _ = m.Delete(ctx, 1, "rec-1")
	_, _, _ = m.Bulk(ctx, 1, nil)

	want := []string{"GetByID", "List", "All", "Create", "Update", "Delete", "Bulk"}
	if len(calls) != len(want) {
		t.Fatalf("expected %d calls, got %d: %v", len(want), len(calls), calls)
	}
	for i, name := range want {
		if calls[i] != name {
			t.Errorf("call %d: expected %s, got %s", i, name, calls[i])
		}
	}
}
