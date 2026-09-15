package cli_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/v2/internal/cli"
	"github.com/contentways/poweradmin-cli/v2/internal/state"
)

func TestNewRootCommand(t *testing.T) {
	s := state.New(
		"https://dns.example.com",
		"secret",
	)

	cmd := cli.NewRootCommand(s)

	if cmd.Use != "poweradmin" {
		t.Fatalf("Use = %q", cmd.Use)
	}

	if len(cmd.Commands()) != 6 {
		t.Fatalf("got %d commands, want 6", len(cmd.Commands()))
	}
}
