// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package server_test

import (
	"context"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/server"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestServerStatusClientError(t *testing.T) {
	fx := testutil.NewFixture(t)
	fx.State.URL = ""

	err := fx.Run(server.NewStatusCmd(nil), nil)
	if err == nil || !strings.Contains(err.Error(), "failed to create client") {
		t.Fatalf("err = %v, want client creation error", err)
	}
}

func TestServerStatusUptimeFormats(t *testing.T) {
	cases := []struct {
		seconds *int
		want    string
	}{
		{new(59), "Uptime:    0m 59s"},
		{new(3725), "Uptime:    1h 2m"},
		{new(93784), "Uptime:    1d 2h 3m"},
		{nil, "Uptime:    unknown"},
	}
	for _, tc := range cases {
		t.Run(tc.want, func(t *testing.T) {
			fx := fixtureWithStatus(t, func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error) {
				return &poweradmin.ServerStatus{Running: false, UptimeSeconds: tc.seconds}, nil, nil
			})
			if err := fx.Run(server.NewStatusCmd(nil), nil); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			out := fx.Stdout.String()
			if !strings.Contains(out, tc.want) || !strings.Contains(out, "Running:   no") {
				t.Errorf("expected %q, got:\n%s", tc.want, out)
			}
		})
	}
}

func TestServerStatusShowMetricsAndNoSlaves(t *testing.T) {
	fx := fixtureWithStatus(t, func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error) {
		s := runningStatus()
		s.Slaves = nil
		return s, nil, nil
	})

	if err := fx.Run(server.NewStatusCmd(nil), []string{"--show-metrics", "--include-slaves"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "udp-queries") || !strings.Contains(out, "No autoprimary servers configured.") {
		t.Errorf("unexpected output:\n%s", out)
	}
}

func TestServerStatusSlaveWithoutDetails(t *testing.T) {
	fx := fixtureWithStatus(t, func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error) {
		s := runningStatus()
		s.Slaves = []poweradmin.SlaveStatus{{IP: "192.0.2.2", Status: "ok"}}
		return s, nil, nil
	})

	if err := fx.Run(server.NewStatusCmd(nil), []string{"--include-slaves"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fx.Stdout.String()
	if !strings.Contains(out, "192.0.2.2") || !strings.Contains(out, "-") {
		t.Errorf("unexpected output:\n%s", out)
	}
}

func TestServerStatusYAML(t *testing.T) {
	fx := fixtureWithStatus(t, func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error) {
		return runningStatus(), nil, nil
	})

	if err := fx.Run(server.NewStatusCmd(nil), []string{"-o", "yaml"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fx.Stdout.String(), "daemon_type: authoritative") {
		t.Errorf("unexpected YAML:\n%s", fx.Stdout.String())
	}
}

func TestNewServerCommand(t *testing.T) {
	cmd := server.NewServerCommand(nil)
	if cmd.Use != "server" {
		t.Fatalf("Use = %q", cmd.Use)
	}
	if len(cmd.Commands()) != 1 || cmd.Commands()[0].Name() != "status" {
		t.Fatalf("subcommands = %v, want [status]", cmd.Commands())
	}
}
