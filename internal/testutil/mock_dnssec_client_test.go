// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil_test

import (
	"context"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestMockDNSSECClientUnsetFieldsReturnZeroValues(t *testing.T) {
	m := &testutil.MockDNSSECClient{}
	ctx := context.Background()

	if keys, resp, err := m.ListKeys(ctx, 1); keys != nil || resp != nil || err != nil {
		t.Errorf("ListKeys: expected all nil, got %v, %v, %v", keys, resp, err)
	}
	if k, resp, err := m.GetKey(ctx, 1, 2); k != nil || resp != nil || err != nil {
		t.Errorf("GetKey: expected all nil, got %v, %v, %v", k, resp, err)
	}
	if k, resp, err := m.AddKey(ctx, 1, poweradmin.DNSSECKeyCreateOpts{}); k != nil || resp != nil || err != nil {
		t.Errorf("AddKey: expected all nil, got %v, %v, %v", k, resp, err)
	}
	if k, resp, err := m.SetKeyActive(ctx, 1, 2, true); k != nil || resp != nil || err != nil {
		t.Errorf("SetKeyActive: expected all nil, got %v, %v, %v", k, resp, err)
	}
	if resp, err := m.DeleteKey(ctx, 1, 2); resp != nil || err != nil {
		t.Errorf("DeleteKey: expected all nil, got %v, %v", resp, err)
	}
	if resp, err := m.Rectify(ctx, 1); resp != nil || err != nil {
		t.Errorf("Rectify: expected all nil, got %v, %v", resp, err)
	}
}

func TestMockDNSSECClientSetFieldsDelegate(t *testing.T) {
	ctx := context.Background()
	key := &poweradmin.DNSSECKey{ID: 2}
	var calls []string
	m := &testutil.MockDNSSECClient{
		ListKeysFn: func(ctx context.Context, zoneID int) ([]*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			calls = append(calls, "ListKeys")
			return []*poweradmin.DNSSECKey{key}, nil, nil
		},
		GetKeyFn: func(ctx context.Context, zoneID, keyID int) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			calls = append(calls, "GetKey")
			return key, nil, nil
		},
		AddKeyFn: func(ctx context.Context, zoneID int, opts poweradmin.DNSSECKeyCreateOpts) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			calls = append(calls, "AddKey")
			return key, nil, nil
		},
		SetKeyActiveFn: func(ctx context.Context, zoneID, keyID int, active bool) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			calls = append(calls, "SetKeyActive")
			return key, nil, nil
		},
		DeleteKeyFn: func(ctx context.Context, zoneID, keyID int) (*poweradmin.Response, error) {
			calls = append(calls, "DeleteKey")
			return nil, nil
		},
		RectifyFn: func(ctx context.Context, zoneID int) (*poweradmin.Response, error) {
			calls = append(calls, "Rectify")
			return nil, nil
		},
	}

	m.ListKeys(ctx, 1)
	if k, _, _ := m.GetKey(ctx, 1, 2); k != key {
		t.Errorf("GetKey: expected delegated key, got %v", k)
	}
	m.AddKey(ctx, 1, poweradmin.DNSSECKeyCreateOpts{})
	m.SetKeyActive(ctx, 1, 2, true)
	m.DeleteKey(ctx, 1, 2)
	m.Rectify(ctx, 1)

	if len(calls) != 6 {
		t.Errorf("expected all six functions to be called, got %v", calls)
	}
}

func TestMockServerClientUnsetFieldsReturnZeroValues(t *testing.T) {
	m := &testutil.MockServerClient{}

	if s, resp, err := m.Status(context.Background(), poweradmin.ServerStatusOpts{}); s != nil || resp != nil || err != nil {
		t.Errorf("Status: expected all nil, got %v, %v, %v", s, resp, err)
	}
}

func TestMockServerClientSetFieldsDelegate(t *testing.T) {
	status := &poweradmin.ServerStatus{Running: true}
	m := &testutil.MockServerClient{
		StatusFn: func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error) {
			return status, nil, nil
		},
	}

	if s, _, _ := m.Status(context.Background(), poweradmin.ServerStatusOpts{}); s != status {
		t.Errorf("Status: expected delegated status, got %v", s)
	}
}

func TestFixtureWithServerAndDNSSEC(t *testing.T) {
	server := &testutil.MockServerClient{}
	dnssec := &testutil.MockDNSSECClient{}
	fx := testutil.NewFixtureWithMocks(t, nil, nil).WithServer(server).WithDNSSEC(dnssec)

	if fx.State.MockClient.Server != server || fx.State.MockClient.DNSSEC != dnssec {
		t.Error("expected the mocks to be injected into the client")
	}
}
