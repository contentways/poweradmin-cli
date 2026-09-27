// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package server_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/server"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func runningStatus() *poweradmin.ServerStatus {
	uptime := 93784
	lastChecked := "2026-09-25 22:33:58"
	errMsg := "Network unreachable"
	return &poweradmin.ServerStatus{
		Running:       true,
		ServerID:      "localhost",
		DaemonType:    "authoritative",
		Version:       "4.9.17",
		UptimeSeconds: &uptime,
		Metrics:       map[string]string{"uptime": "93784", "udp-queries": "12"},
		Slaves:        []poweradmin.SlaveStatus{{IP: "192.0.2.1", Status: "unreachable", LastChecked: &lastChecked, Error: &errMsg}},
	}
}

func fixtureWithStatus(t *testing.T, fn func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error)) *testutil.Fixture {
	return testutil.NewFixtureWithMocks(t, nil, nil).WithServer(&testutil.MockServerClient{StatusFn: fn})
}

func TestServerStatusTable(t *testing.T) {
	fx := fixtureWithStatus(t, func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error) {
		if len(opts.Metrics) != 0 || opts.IncludeSlaves {
			t.Errorf("opts = %+v, want defaults", opts)
		}
		return runningStatus(), nil, nil
	})

	if err := fx.Run(server.NewStatusCmd(nil), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	for _, want := range []string{"Running:   yes", "Version:   4.9.17", "Uptime:    1d 2h 3m"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
	if strings.Contains(out, "METRIC") {
		t.Errorf("metrics must not be shown by default, got:\n%s", out)
	}
}

func TestServerStatusMetricsAndSlaves(t *testing.T) {
	fx := fixtureWithStatus(t, func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error) {
		if strings.Join(opts.Metrics, ",") != "uptime,udp-queries" || !opts.IncludeSlaves {
			t.Errorf("opts = %+v", opts)
		}
		return runningStatus(), nil, nil
	})

	err := fx.Run(server.NewStatusCmd(nil), []string{"--metrics", "uptime,udp-queries", "--include-slaves"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	for _, want := range []string{"METRIC", "udp-queries", "192.0.2.1", "unreachable", "Network unreachable"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestServerStatusJSON(t *testing.T) {
	fx := fixtureWithStatus(t, func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error) {
		return runningStatus(), nil, nil
	})

	if err := fx.Run(server.NewStatusCmd(nil), []string{"-o", "json"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(fx.Stdout.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, fx.Stdout.String())
	}
	if got["running"] != true || got["version"] != "4.9.17" || got["uptime_seconds"] != float64(93784) {
		t.Errorf("unexpected JSON: %v", got)
	}
}

func TestServerStatusUnreachable(t *testing.T) {
	fx := fixtureWithStatus(t, func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error) {
		return nil, nil, &poweradmin.APIError{StatusCode: 503, Message: "PowerDNS server is not reachable"}
	})

	err := fx.Run(server.NewStatusCmd(nil), nil)
	if !errors.Is(err, server.ErrPowerDNSUnreachable) {
		t.Fatalf("err = %v, want ErrPowerDNSUnreachable", err)
	}
}

func TestServerStatusForbidden(t *testing.T) {
	fx := fixtureWithStatus(t, func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error) {
		return nil, nil, &poweradmin.APIError{StatusCode: 403, Message: "Forbidden"}
	})

	err := fx.Run(server.NewStatusCmd(nil), []string{"--include-slaves"})
	if err == nil || errors.Is(err, server.ErrPowerDNSUnreachable) || !strings.Contains(err.Error(), "Forbidden") {
		t.Fatalf("err = %v, want the API's 403 error", err)
	}
}
