// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil

import (
	"context"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
)

// MockRecordClient implements poweradmin.IRecordClient for testing.
// Each method can be overridden by setting the corresponding function field.
// Unset fields return zero values and no error by default.
type MockRecordClient struct {
	GetByIDFn func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Record, *poweradmin.Response, error)
	ListFn    func(ctx context.Context, zoneID int, opts poweradmin.RecordListOpts) ([]*poweradmin.Record, *poweradmin.Response, error)
	AllFn     func(ctx context.Context, zoneID int) ([]*poweradmin.Record, error)
	CreateFn  func(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error)
	UpdateFn  func(ctx context.Context, zoneID int, recordID string, opts poweradmin.RecordUpdateOpts) (*poweradmin.Record, *poweradmin.Response, error)
	DeleteFn  func(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error)
	BulkFn    func(ctx context.Context, zoneID int, ops []poweradmin.BulkRecordOperation) (*poweradmin.BulkRecordsResult, *poweradmin.Response, error)
}

func (m *MockRecordClient) GetByID(ctx context.Context, zoneID int, recordID string) (*poweradmin.Record, *poweradmin.Response, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, zoneID, recordID)
	}
	return nil, nil, nil
}

func (m *MockRecordClient) List(ctx context.Context, zoneID int, opts poweradmin.RecordListOpts) ([]*poweradmin.Record, *poweradmin.Response, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, zoneID, opts)
	}
	return nil, nil, nil
}

func (m *MockRecordClient) All(ctx context.Context, zoneID int) ([]*poweradmin.Record, error) {
	if m.AllFn != nil {
		return m.AllFn(ctx, zoneID)
	}
	return nil, nil
}

func (m *MockRecordClient) Create(ctx context.Context, zoneID int, opts poweradmin.RecordCreateOpts) (string, *poweradmin.Response, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, zoneID, opts)
	}
	return "", nil, nil
}

func (m *MockRecordClient) Update(ctx context.Context, zoneID int, recordID string, opts poweradmin.RecordUpdateOpts) (*poweradmin.Record, *poweradmin.Response, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, zoneID, recordID, opts)
	}
	return nil, nil, nil
}

func (m *MockRecordClient) Delete(ctx context.Context, zoneID int, recordID string) (*poweradmin.Response, error) {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, zoneID, recordID)
	}
	return nil, nil
}

func (m *MockRecordClient) Bulk(ctx context.Context, zoneID int, ops []poweradmin.BulkRecordOperation) (*poweradmin.BulkRecordsResult, *poweradmin.Response, error) {
	if m.BulkFn != nil {
		return m.BulkFn(ctx, zoneID, ops)
	}
	return nil, nil, nil
}
