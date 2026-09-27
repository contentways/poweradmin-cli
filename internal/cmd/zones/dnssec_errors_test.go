// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

var errAPI = errors.New("api error")

// failingDNSSECFixture returns a fixture whose zone lookup works but whose
// DNSSEC calls all fail.
func failingDNSSECFixture(t *testing.T) *testutil.Fixture {
	t.Helper()
	zone := dnssecZoneMock(nil)
	zone.GetDNSSECFn = func(ctx context.Context, id int) (*poweradmin.ZoneDNSSEC, *poweradmin.Response, error) {
		return nil, nil, errAPI
	}
	zone.SetDNSSECFn = func(ctx context.Context, id int, enabled bool) (*poweradmin.ZoneDNSSEC, *poweradmin.Response, error) {
		return nil, nil, errAPI
	}
	return testutil.NewFixtureWithMocks(t, zone, nil).WithDNSSEC(&testutil.MockDNSSECClient{
		ListKeysFn: func(ctx context.Context, zoneID int) ([]*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			return nil, nil, errAPI
		},
		AddKeyFn: func(ctx context.Context, zoneID int, opts poweradmin.DNSSECKeyCreateOpts) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			return nil, nil, errAPI
		},
		SetKeyActiveFn: func(ctx context.Context, zoneID, keyID int, active bool) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			return nil, nil, errAPI
		},
		DeleteKeyFn: func(ctx context.Context, zoneID, keyID int) (*poweradmin.Response, error) {
			return nil, errAPI
		},
		RectifyFn: func(ctx context.Context, zoneID int) (*poweradmin.Response, error) {
			return nil, errAPI
		},
	})
}

func TestZonesDNSSECAPIErrors(t *testing.T) {
	cases := map[string]struct {
		args []string
		want string
	}{
		"status":     {[]string{"dnssec", "status"}, "failed to get DNSSEC status"},
		"enable":     {[]string{"dnssec", "enable"}, "failed to enable DNSSEC"},
		"disable":    {[]string{"dnssec", "disable", "--yes"}, "failed to disable DNSSEC"},
		"rectify":    {[]string{"dnssec", "rectify"}, "failed to rectify zone"},
		"keys list":  {[]string{"dnssec", "keys", "list"}, "failed to list DNSSEC keys"},
		"keys add":   {[]string{"dnssec", "keys", "add", "--type", "csk", "--algorithm", "ecdsa256", "--bits", "256"}, "failed to add DNSSEC key"},
		"activate":   {[]string{"dnssec", "keys", "activate", "--key-id", "3"}, "failed to activate DNSSEC key"},
		"deactivate": {[]string{"dnssec", "keys", "deactivate", "--key-id", "3"}, "failed to deactivate DNSSEC key"},
		"delete":     {[]string{"dnssec", "keys", "delete", "--key-id", "3", "--yes"}, "failed to delete DNSSEC key"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fx := failingDNSSECFixture(t)
			err := fx.Run(zones.NewZonesCommand(nil), append(tc.args, "--name", "example.com"))
			if err == nil || !strings.Contains(err.Error(), tc.want) || !errors.Is(err, errAPI) {
				t.Fatalf("err = %v, want %q wrapping the API error", err, tc.want)
			}
		})
	}
}

func TestZonesDNSSECZoneNotFound(t *testing.T) {
	zone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return nil, nil, errAPI
		},
	}
	fx := testutil.NewFixtureWithMocks(t, zone, nil).WithDNSSEC(&testutil.MockDNSSECClient{})

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "keys", "list", "--name", "missing.example"}); err == nil {
		t.Fatal("expected an error for an unknown zone")
	}
}

func TestZonesDNSSECClientError(t *testing.T) {
	fx := testutil.NewFixture(t)
	fx.State.URL = ""

	err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "status", "--name", "example.com"})
	if err == nil || !strings.Contains(err.Error(), "failed to create client") {
		t.Fatalf("err = %v, want client creation error", err)
	}
}

func TestZonesDNSSECStatusYAML(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(signedDNSSEC()), nil)

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "status", "--name", "example.com", "-o", "yaml"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "zone: example.com") || !strings.Contains(out, "keytag: 46395") {
		t.Errorf("unexpected YAML:\n%s", out)
	}
}

func TestZonesDNSSECStatusUnsignedWithoutDS(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(&poweradmin.ZoneDNSSEC{}), nil)

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "status", "--name", "example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "DNSSEC: unsigned") || strings.Contains(out, "DS records") {
		t.Errorf("unexpected output:\n%s", out)
	}
}

func TestZonesDNSSECStatusPresignedTable(t *testing.T) {
	d := signedDNSSEC()
	d.Presigned = true
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(d), nil)

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "status", "--name", "example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "presigned, managed at the primary") {
		t.Errorf("unexpected output:\n%s", fx.Stdout.String())
	}
}

func TestZonesDNSSECRectifyQuiet(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{})

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "rectify", "--name", "example.com", "-q"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fx.Stdout.Len() != 0 {
		t.Errorf("expected no output with --quiet, got:\n%s", fx.Stdout.String())
	}
}

func TestZonesDNSSECKeysAddQuietAndJSON(t *testing.T) {
	newFx := func(t *testing.T) *testutil.Fixture {
		return testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{
			AddKeyFn: func(ctx context.Context, zoneID int, opts poweradmin.DNSSECKeyCreateOpts) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
				return testKey(5, false), nil, nil
			},
		})
	}
	args := []string{"dnssec", "keys", "add", "--name", "example.com", "--type", "zsk", "--algorithm", "ecdsa256", "--bits", "256"}

	quiet := newFx(t)
	if err := quiet.Run(zones.NewZonesCommand(nil), append(args, "-q")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quiet.Stdout.Len() != 0 {
		t.Errorf("expected no output with --quiet, got:\n%s", quiet.Stdout.String())
	}

	asJSON := newFx(t)
	if err := asJSON.Run(zones.NewZonesCommand(nil), append(args, "-o", "json")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(asJSON.Stdout.String(), `"id": 5`) {
		t.Errorf("unexpected JSON:\n%s", asJSON.Stdout.String())
	}
}

func TestZonesDNSSECKeysUnknownAlgorithmName(t *testing.T) {
	key := testKey(9, true)
	key.Algorithm = ""
	key.AlgorithmID = 253
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{
		ListKeysFn: func(ctx context.Context, zoneID int) ([]*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			return []*poweradmin.DNSSECKey{key}, nil, nil
		},
		SetKeyActiveFn: func(ctx context.Context, zoneID, keyID int, active bool) (*poweradmin.DNSSECKey, *poweradmin.Response, error) {
			return key, nil, nil
		},
	})

	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "keys", "list", "--name", "example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "253") {
		t.Errorf("expected the numeric algorithm in the table, got:\n%s", fx.Stdout.String())
	}

	fx.Stdout.Reset()
	if err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "keys", "activate", "--name", "example.com", "--key-id", "9"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "algorithm 253") {
		t.Errorf("expected the numeric algorithm in the message, got:\n%s", fx.Stdout.String())
	}
}

func TestZonesDNSSECKeysDeleteQuiet(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{})

	err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "keys", "delete", "--name", "example.com", "--key-id", "3", "--yes", "-q"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fx.Stdout.Len() != 0 {
		t.Errorf("expected no output with --quiet, got:\n%s", fx.Stdout.String())
	}
}

func TestZonesDNSSECKeysActivateRequiresKeyID(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{})

	err := fx.Run(zones.NewZonesCommand(nil), []string{"dnssec", "keys", "activate", "--name", "example.com"})
	if err == nil || !strings.Contains(err.Error(), "--key-id") {
		t.Fatalf("err = %v, want missing key id error", err)
	}
}

func TestZonesDNSSECCommandsRequireZone(t *testing.T) {
	for _, args := range [][]string{
		{"dnssec", "enable"},
		{"dnssec", "disable", "--yes"},
		{"dnssec", "rectify"},
		{"dnssec", "keys", "list"},
		{"dnssec", "keys", "add", "--type", "csk", "--algorithm", "ecdsa256", "--bits", "256"},
		{"dnssec", "keys", "activate", "--key-id", "3"},
		{"dnssec", "keys", "delete", "--key-id", "3", "--yes"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			fx := testutil.NewFixtureWithMocks(t, dnssecZoneMock(nil), nil).WithDNSSEC(&testutil.MockDNSSECClient{})
			err := fx.Run(zones.NewZonesCommand(nil), args)
			if err == nil || !strings.Contains(err.Error(), "--name or --id") {
				t.Fatalf("err = %v, want missing zone error", err)
			}
		})
	}
}

func TestZonesDNSSECConfirmationDeclined(t *testing.T) {
	for _, args := range [][]string{
		{"dnssec", "disable", "--name", "example.com"},
		{"dnssec", "keys", "delete", "--name", "example.com", "--key-id", "3"},
	} {
		t.Run(args[1], func(t *testing.T) {
			called := false
			zone := dnssecZoneMock(nil)
			zone.SetDNSSECFn = func(ctx context.Context, id int, enabled bool) (*poweradmin.ZoneDNSSEC, *poweradmin.Response, error) {
				called = true
				return &poweradmin.ZoneDNSSEC{}, nil, nil
			}
			fx := testutil.NewFixtureWithMocks(t, zone, nil).WithDNSSEC(&testutil.MockDNSSECClient{
				DeleteKeyFn: func(ctx context.Context, zoneID, keyID int) (*poweradmin.Response, error) {
					called = true
					return nil, nil
				},
			})
			testutil.WithStdin(t, "n\n")

			if err := fx.Run(zones.NewZonesCommand(nil), args); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if called {
				t.Error("declined confirmation must not change anything")
			}
			if !strings.Contains(fx.Stdout.String(), "Aborted.") {
				t.Errorf("expected abort message, got:\n%s", fx.Stdout.String())
			}
		})
	}
}
