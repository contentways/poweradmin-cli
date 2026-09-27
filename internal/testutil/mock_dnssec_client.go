// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil

import (
	"context"

	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

// MockDNSSECClient implements poweradmin.IDNSSECClient for testing.
// Each method can be overridden by setting the corresponding function field.
type MockDNSSECClient struct {
	ListKeysFn     func(ctx context.Context, zoneID int) ([]*poweradmin.DNSSECKey, *poweradmin.Response, error)
	GetKeyFn       func(ctx context.Context, zoneID, keyID int) (*poweradmin.DNSSECKey, *poweradmin.Response, error)
	AddKeyFn       func(ctx context.Context, zoneID int, opts poweradmin.DNSSECKeyCreateOpts) (*poweradmin.DNSSECKey, *poweradmin.Response, error)
	SetKeyActiveFn func(ctx context.Context, zoneID, keyID int, active bool) (*poweradmin.DNSSECKey, *poweradmin.Response, error)
	DeleteKeyFn    func(ctx context.Context, zoneID, keyID int) (*poweradmin.Response, error)
	RectifyFn      func(ctx context.Context, zoneID int) (*poweradmin.Response, error)
}

func (m *MockDNSSECClient) ListKeys(ctx context.Context, zoneID int) ([]*poweradmin.DNSSECKey, *poweradmin.Response, error) {
	if m.ListKeysFn != nil {
		return m.ListKeysFn(ctx, zoneID)
	}
	return nil, nil, nil
}

func (m *MockDNSSECClient) GetKey(ctx context.Context, zoneID, keyID int) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
	if m.GetKeyFn != nil {
		return m.GetKeyFn(ctx, zoneID, keyID)
	}
	return nil, nil, nil
}

func (m *MockDNSSECClient) AddKey(ctx context.Context, zoneID int, opts poweradmin.DNSSECKeyCreateOpts) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
	if m.AddKeyFn != nil {
		return m.AddKeyFn(ctx, zoneID, opts)
	}
	return nil, nil, nil
}

func (m *MockDNSSECClient) SetKeyActive(ctx context.Context, zoneID, keyID int, active bool) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
	if m.SetKeyActiveFn != nil {
		return m.SetKeyActiveFn(ctx, zoneID, keyID, active)
	}
	return nil, nil, nil
}

func (m *MockDNSSECClient) DeleteKey(ctx context.Context, zoneID, keyID int) (*poweradmin.Response, error) {
	if m.DeleteKeyFn != nil {
		return m.DeleteKeyFn(ctx, zoneID, keyID)
	}
	return nil, nil
}

func (m *MockDNSSECClient) Rectify(ctx context.Context, zoneID int) (*poweradmin.Response, error) {
	if m.RectifyFn != nil {
		return m.RectifyFn(ctx, zoneID)
	}
	return nil, nil
}
