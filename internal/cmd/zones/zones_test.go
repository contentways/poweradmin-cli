package zones_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/v2/internal/cmd/zones"
)

func TestNewZonesCommand(t *testing.T) {
	cmd := zones.NewZonesCommand(nil)

	if cmd.Use != "zones" {
		t.Fatalf("Use = %q", cmd.Use)
	}

	if len(cmd.Commands()) != 10 {
		t.Fatalf("got %d subcommands, want 6", len(cmd.Commands()))
	}
}
