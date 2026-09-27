// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func dnssecZoneMock(d *poweradmin.ZoneDNSSEC) *testutil.MockZoneClient {
	return &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 7, Name: name, Type: "MASTER"}, nil, nil
		},
		GetDNSSECFn: func(ctx context.Context, id int) (*poweradmin.ZoneDNSSEC, *poweradmin.Response, error) {
			return d, nil, nil
		},
		SetDNSSECFn: func(ctx context.Context, id int, enabled bool) (*poweradmin.ZoneDNSSEC, *poweradmin.Response, error) {
			return &poweradmin.ZoneDNSSEC{Enabled: enabled}, nil, nil
		},
	}
}

func signedDNSSEC() *poweradmin.ZoneDNSSEC {
	return &poweradmin.ZoneDNSSEC{
		Enabled:   true,
		DSRecords: []poweradmin.DSRecord{{KeyTag: 46395, Algorithm: 13, DigestType: 2, Digest: "3dd8ee7d"}},
	}
}

func testKey(id int, active bool) *poweradmin.DNSSECKey {
	dnskey := "257 3 13 mdsswUyr"
	return &poweradmin.DNSSECKey{
		ID: id, Type: poweradmin.DNSSECKeyTypeCSK, KeyTag: 46395, Algorithm: "ecdsa256",
		AlgorithmID: 13, Bits: 256, Active: active, DNSKey: &dnskey, DS: []string{"46395 13 2 3dd8ee7d"},
	}
}

func TestZonesDNSSECStatus(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(signedDNSSEC()), nil)

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "status", "--name", "example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	for _, want := range []string{"DNSSEC: signed", "example.com. IN DS 46395 13 2 3dd8ee7d"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestZonesDNSSECStatusPresignedJSON(t *testing.T) {
	d := signedDNSSEC()
	d.Presigned = true
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(d), nil)

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "status", "--name", "example.com", "-o", "json"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(fx.Stdout.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, fx.Stdout.String())
	}
	if got["zone"] != "example.com" || got["enabled"] != true || got["presigned"] != true {
		t.Errorf("unexpected JSON: %v", got)
	}
}

func TestZonesDNSSECEnable(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil)

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "enable", "--name", "example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "DNSSEC: signed") {
		t.Errorf("expected signed status, got:\n%s", fx.Stdout.String())
	}
}

func TestZonesDNSSECDisableRequiresConfirmation(t *testing.T) {
	called := false
	mock := dnssecZoneMock(nil)
	mock.SetDNSSECFn = func(ctx context.Context, id int, enabled bool) (*poweradmin.ZoneDNSSEC, *poweradmin.Response, error) {
		called = true
		return &poweradmin.ZoneDNSSEC{}, nil, nil
	}
	fx := testutil.NewFixtureWithMocks(t, mock, nil)

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "disable", "--name", "example.com", "--yes"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected SetDNSSEC(false) with --yes")
	}
	if !strings.Contains(fx.Stdout.String(), "DNSSEC: unsigned") {
		t.Errorf("expected unsigned status, got:\n%s", fx.Stdout.String())
	}
}

func TestZonesDNSSECRequiresZone(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil)

	err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "status"})
	if err == nil || !strings.Contains(err.Error(), "--name or --id") {
		t.Fatalf("err = %v, want missing zone error", err)
	}
}

func TestZonesDNSSECRectify(t *testing.T) {
	gotZone := 0
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{
		RectifyFn: func(ctx context.Context, zoneID int) (*poweradmin.Response, error) {
			gotZone = zoneID
			return nil, nil
		},
	})

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "rectify", "--name", "example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotZone != 7 {
		t.Errorf("Rectify(%d), want zone 7", gotZone)
	}
	if !strings.Contains(fx.Stdout.String(), "rectified zone example.com") {
		t.Errorf("unexpected output:\n%s", fx.Stdout.String())
	}
}

func TestZonesDNSSECKeysList(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{
		ListKeysFn: func(ctx context.Context, zoneID int) ([]*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			return []*poweradmin.DNSSECKey{testKey(1, true), testKey(2, false)}, nil, nil
		},
	})

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "keys", "list", "--name", "example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	for _, want := range []string{"KEYTAG", "CSK", "46395", "ecdsa256", "yes", "no"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestZonesDNSSECKeysListJSON(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{
		ListKeysFn: func(ctx context.Context, zoneID int) ([]*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			return []*poweradmin.DNSSECKey{testKey(1, true)}, nil, nil
		},
	})

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "keys", "list", "--name", "example.com", "-o", "json"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got struct {
		Zone  string `json:"zone"`
		Count int    `json:"count"`
		Keys  []struct {
			ID     int      `json:"id"`
			DNSKey string   `json:"dnskey"`
			DS     []string `json:"ds"`
		} `json:"keys"`
	}
	if err := json.Unmarshal(fx.Stdout.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, fx.Stdout.String())
	}
	if got.Zone != "example.com" || got.Count != 1 || got.Keys[0].DNSKey == "" || len(got.Keys[0].DS) != 1 {
		t.Errorf("unexpected JSON: %+v", got)
	}
}

func TestZonesDNSSECKeysAdd(t *testing.T) {
	var gotOpts poweradmin.DNSSECKeyCreateOpts
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{
		AddKeyFn: func(ctx context.Context, zoneID int, opts poweradmin.DNSSECKeyCreateOpts) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			gotOpts = opts
			return testKey(3, opts.Active), nil, nil
		},
	})

	err := fx.Run(zones.NewZonesCommand(nil), []string{
		"dnssec", "keys", "add", "--name", "example.com",
		"--type", "CSK", "--algorithm", "ECDSA256", "--bits", "256", "--active",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := poweradmin.DNSSECKeyCreateOpts{Type: poweradmin.DNSSECKeyTypeCSK, Algorithm: "ecdsa256", Bits: 256, Active: true}
	if gotOpts != want {
		t.Errorf("opts = %+v, want %+v", gotOpts, want)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "key 3") || !strings.Contains(out, "DS 46395 13 2 3dd8ee7d") {
		t.Errorf("unexpected output:\n%s", out)
	}
}

func TestZonesDNSSECKeysAddValidation(t *testing.T) {
	cases := map[string][]string{
		"bad type":          {"--type", "abc", "--algorithm", "ecdsa256", "--bits", "256"},
		"missing algorithm": {"--type", "csk", "--bits", "256"},
		"missing bits":      {"--type", "csk", "--algorithm", "ecdsa256"},
	}
	for name, flags := range cases {
		t.Run(name, func(t *testing.T) {
			fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{})
			args := append([]string{"dnssec", "keys", "add", "--name", "example.com"}, flags...)
			if err := fx.Run(zones.NewZonesCommand(nil), args); err == nil {
				t.Fatal("expected a validation error")
			}
		})
	}
}

func TestZonesDNSSECKeysAddServerError(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{
		AddKeyFn: func(ctx context.Context, zoneID int, opts poweradmin.DNSSECKeyCreateOpts) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			return nil, nil, &poweradmin.APIError{StatusCode: 400, Message: "ECDSA P-256 algorithm must use 256 bits"}
		},
	})

	err := fx.Run(zones.NewZonesCommand(nil), []string{
		"dnssec", "keys", "add", "--name", "example.com", "--type", "csk", "--algorithm", "ecdsa256", "--bits", "384",
	})
	if err == nil || !strings.Contains(err.Error(), "must use 256 bits") {
		t.Fatalf("err = %v, want the server's validation message", err)
	}
}

func TestZonesDNSSECKeysActivateDeactivate(t *testing.T) {
	for _, tc := range []struct {
		cmd    string
		active bool
	}{{"activate", true}, {"deactivate", false}} {
		t.Run(tc.cmd, func(t *testing.T) {
			var gotKey int
			var gotActive *bool
			fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{
				SetKeyActiveFn: func(ctx context.Context, zoneID, keyID int, active bool) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
					gotKey, gotActive = keyID, &active
					return testKey(keyID, active), nil, nil
				},
			})

			if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "keys", tc.cmd, "--name", "example.com", "--key-id", "3"}); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotKey != 3 || gotActive == nil || *gotActive != tc.active {
				t.Errorf("SetKeyActive(key %d, %v), want (3, %v)", gotKey, gotActive, tc.active)
			}
			if !strings.Contains(fx.Stdout.String(), fmt.Sprintf("active: %s", map[bool]string{true: "yes", false: "no"}[tc.active])) {
				t.Errorf("unexpected output:\n%s", fx.Stdout.String())
			}
		})
	}
}

func TestZonesDNSSECKeysDelete(t *testing.T) {
	gotKey := 0
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{
		DeleteKeyFn: func(ctx context.Context, zoneID, keyID int) (*poweradmin.Response, error) {
			gotKey = keyID
			return nil, nil
		},
	})

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "keys", "delete", "--name", "example.com", "--key-id", "3", "--yes"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotKey != 3 {
		t.Errorf("DeleteKey(%d), want 3", gotKey)
	}
}

func TestZonesDNSSECKeysDeleteRequiresKeyID(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{})

	err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "keys", "delete", "--name", "example.com", "--yes"})
	if err == nil || !strings.Contains(err.Error(), "--key-id") {
		t.Fatalf("err = %v, want missing key id error", err)
	}
}
