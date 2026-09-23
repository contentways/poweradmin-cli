// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// mock_zone_client_test.go exercises every method of MockZoneClient on both
// its paths: with the corresponding *Fn field set (delegates to it) and
// unset (returns the documented zero values). This is here purely to close
// coverage gaps on mock plumbing that no command test happens to exercise
// (e.g. zone ownership and DNSSEC, which have no CLI command yet) — it
// verifies the mock behaves as documented, not any CLI behavior.
package testutil_test

import (
	"context"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestMockZoneClientUnsetFieldsReturnZeroValues(t *testing.T) {
	m := &testutil.MockZoneClient{}
	ctx := context.Background()

	if z, resp, err := m.GetByID(ctx, 1); z != nil || resp != nil || err != nil {
		t.Errorf("GetByID: expected all nil, got %v, %v, %v", z, resp, err)
	}
	if z, resp, err := m.GetByName(ctx, "x"); z != nil || resp != nil || err != nil {
		t.Errorf("GetByName: expected all nil, got %v, %v, %v", z, resp, err)
	}
	if zs, resp, err := m.List(ctx, poweradmin.ListOpts{}); zs != nil || resp != nil || err != nil {
		t.Errorf("List: expected all nil, got %v, %v, %v", zs, resp, err)
	}
	if zs, err := m.All(ctx); zs != nil || err != nil {
		t.Errorf("All: expected all nil, got %v, %v", zs, err)
	}
	if id, resp, err := m.Create(ctx, poweradmin.ZoneCreateOpts{}); id != 0 || resp != nil || err != nil {
		t.Errorf("Create: expected zero values, got %v, %v, %v", id, resp, err)
	}
	if z, resp, err := m.Update(ctx, 1, poweradmin.ZoneUpdateOpts{}); z != nil || resp != nil || err != nil {
		t.Errorf("Update: expected all nil, got %v, %v, %v", z, resp, err)
	}
	if resp, err := m.Delete(ctx, 1); resp != nil || err != nil {
		t.Errorf("Delete: expected all nil, got %v, %v", resp, err)
	}
	if os, resp, err := m.Owners(ctx, 1); os != nil || resp != nil || err != nil {
		t.Errorf("Owners: expected all nil, got %v, %v, %v", os, resp, err)
	}
	if resp, err := m.AddOwner(ctx, 1, 2); resp != nil || err != nil {
		t.Errorf("AddOwner: expected all nil, got %v, %v", resp, err)
	}
	if resp, err := m.AddOwners(ctx, 1, []int{2, 3}); resp != nil || err != nil {
		t.Errorf("AddOwners: expected all nil, got %v, %v", resp, err)
	}
	if resp, err := m.RemoveOwner(ctx, 1, 2); resp != nil || err != nil {
		t.Errorf("RemoveOwner: expected all nil, got %v, %v", resp, err)
	}
	if d, resp, err := m.GetDNSSEC(ctx, 1); d != nil || resp != nil || err != nil {
		t.Errorf("GetDNSSEC: expected all nil, got %v, %v, %v", d, resp, err)
	}
	if d, resp, err := m.SetDNSSEC(ctx, 1, true); d != nil || resp != nil || err != nil {
		t.Errorf("SetDNSSEC: expected all nil, got %v, %v, %v", d, resp, err)
	}
	if md, resp, err := m.ListMetadata(ctx, 1); md != nil || resp != nil || err != nil {
		t.Errorf("ListMetadata: expected all nil, got %v, %v, %v", md, resp, err)
	}
	if md, resp, err := m.GetMetadata(ctx, 1, "ALLOW-AXFR-FROM"); md != nil || resp != nil || err != nil {
		t.Errorf("GetMetadata: expected all nil, got %v, %v, %v", md, resp, err)
	}
	if resp, err := m.SetMetadata(ctx, 1, "ALLOW-AXFR-FROM", []string{"AUTO-NS"}); resp != nil || err != nil {
		t.Errorf("SetMetadata: expected all nil, got %v, %v", resp, err)
	}
	if resp, err := m.DeleteMetadata(ctx, 1, "ALLOW-AXFR-FROM"); resp != nil || err != nil {
		t.Errorf("DeleteMetadata: expected all nil, got %v, %v", resp, err)
	}
}

func TestMockZoneClientSetFieldsDelegate(t *testing.T) {
	var calls []string
	m := &testutil.MockZoneClient{
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.Zone, *poweradmin.Response, error) {
			calls = append(calls, "GetByID")
			return &poweradmin.Zone{ID: id}, nil, nil
		},
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			calls = append(calls, "GetByName")
			return &poweradmin.Zone{Name: name}, nil, nil
		},
		ListFn: func(ctx context.Context, opts poweradmin.ListOpts) ([]*poweradmin.Zone, *poweradmin.Response, error) {
			calls = append(calls, "List")
			return []*poweradmin.Zone{{ID: 1}}, nil, nil
		},
		AllFn: func(ctx context.Context) ([]*poweradmin.Zone, error) {
			calls = append(calls, "All")
			return []*poweradmin.Zone{{ID: 1}}, nil
		},
		CreateFn: func(ctx context.Context, opts poweradmin.ZoneCreateOpts) (int, *poweradmin.Response, error) {
			calls = append(calls, "Create")
			return 42, nil, nil
		},
		UpdateFn: func(ctx context.Context, id int, opts poweradmin.ZoneUpdateOpts) (*poweradmin.Zone, *poweradmin.Response, error) {
			calls = append(calls, "Update")
			return &poweradmin.Zone{ID: id}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			calls = append(calls, "Delete")
			return nil, nil
		},
		OwnersFn: func(ctx context.Context, zoneID int) ([]*poweradmin.ZoneOwner, *poweradmin.Response, error) {
			calls = append(calls, "Owners")
			return []*poweradmin.ZoneOwner{{UserID: 1}}, nil, nil
		},
		AddOwnerFn: func(ctx context.Context, zoneID, userID int) (*poweradmin.Response, error) {
			calls = append(calls, "AddOwner")
			return nil, nil
		},
		AddOwnersFn: func(ctx context.Context, zoneID int, userIDs []int) (*poweradmin.Response, error) {
			calls = append(calls, "AddOwners")
			return nil, nil
		},
		RemoveOwnerFn: func(ctx context.Context, zoneID, userID int) (*poweradmin.Response, error) {
			calls = append(calls, "RemoveOwner")
			return nil, nil
		},
		GetDNSSECFn: func(ctx context.Context, id int) (*poweradmin.ZoneDNSSEC, *poweradmin.Response, error) {
			calls = append(calls, "GetDNSSEC")
			return &poweradmin.ZoneDNSSEC{}, nil, nil
		},
		SetDNSSECFn: func(ctx context.Context, id int, enabled bool) (*poweradmin.ZoneDNSSEC, *poweradmin.Response, error) {
			calls = append(calls, "SetDNSSEC")
			return &poweradmin.ZoneDNSSEC{}, nil, nil
		},
		ListMetadataFn: func(ctx context.Context, zoneID int) ([]*poweradmin.ZoneMetadata, *poweradmin.Response, error) {
			calls = append(calls, "ListMetadata")
			return []*poweradmin.ZoneMetadata{{Kind: "ALLOW-AXFR-FROM"}}, nil, nil
		},
		GetMetadataFn: func(ctx context.Context, zoneID int, kind string) (*poweradmin.ZoneMetadata, *poweradmin.Response, error) {
			calls = append(calls, "GetMetadata")
			return &poweradmin.ZoneMetadata{Kind: kind}, nil, nil
		},
		SetMetadataFn: func(ctx context.Context, zoneID int, kind string, values []string) (*poweradmin.Response, error) {
			calls = append(calls, "SetMetadata")
			return nil, nil
		},
		DeleteMetadataFn: func(ctx context.Context, zoneID int, kind string) (*poweradmin.Response, error) {
			calls = append(calls, "DeleteMetadata")
			return nil, nil
		},
	}
	ctx := context.Background()

	_, _, _ = m.GetByID(ctx, 1)
	_, _, _ = m.GetByName(ctx, "x")
	_, _, _ = m.List(ctx, poweradmin.ListOpts{})
	_, _ = m.All(ctx)
	_, _, _ = m.Create(ctx, poweradmin.ZoneCreateOpts{})
	_, _, _ = m.Update(ctx, 1, poweradmin.ZoneUpdateOpts{})
	_, _ = m.Delete(ctx, 1)
	_, _, _ = m.Owners(ctx, 1)
	_, _ = m.AddOwner(ctx, 1, 2)
	_, _ = m.AddOwners(ctx, 1, []int{2, 3})
	_, _ = m.RemoveOwner(ctx, 1, 2)
	_, _, _ = m.GetDNSSEC(ctx, 1)
	_, _, _ = m.SetDNSSEC(ctx, 1, true)
	_, _, _ = m.ListMetadata(ctx, 1)
	_, _, _ = m.GetMetadata(ctx, 1, "ALLOW-AXFR-FROM")
	_, _ = m.SetMetadata(ctx, 1, "ALLOW-AXFR-FROM", []string{"AUTO-NS"})
	_, _ = m.DeleteMetadata(ctx, 1, "ALLOW-AXFR-FROM")

	want := []string{
		"GetByID", "GetByName", "List", "All", "Create", "Update", "Delete",
		"Owners", "AddOwner", "AddOwners", "RemoveOwner",
		"GetDNSSEC", "SetDNSSEC",
		"ListMetadata", "GetMetadata", "SetMetadata", "DeleteMetadata",
	}
	if len(calls) != len(want) {
		t.Fatalf("expected %d calls, got %d: %v", len(want), len(calls), calls)
	}
	for i, name := range want {
		if calls[i] != name {
			t.Errorf("call %d: expected %s, got %s", i, name, calls[i])
		}
	}
}
